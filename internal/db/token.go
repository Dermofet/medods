package db

import (
	"context"
	"fmt"
	"medods/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (s *source) CreateToken(ctx context.Context, auth *entity.Auth) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	result, err := s.db.Collection("tokens").InsertOne(dbCtx, auth)
	if err != nil {
		return fmt.Errorf("can't insert token: %w", err)
	}
	if result.InsertedID == nil {
		return fmt.Errorf("can't insert token: %w", err)
	}

	return nil
}

func (s *source) GetUserIdFromToken(ctx context.Context, token []byte) (*entity.Auth, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	var auth entity.Auth

	cursor, err := s.db.Collection("tokens").Find(dbCtx, bson.M{"token_hash": bson.M{"$exists": true}})
	if err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.Decode(&auth)
		if err != nil {
			return nil, fmt.Errorf("can't decode token: %w", err)
		}

		err := bcrypt.CompareHashAndPassword(auth.Token, token)
		if err == nil {
			return &auth, nil
		}
	}
	return nil, fmt.Errorf("can't find token in db")
}

func (s *source) DeleteToken(ctx context.Context, token []byte) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	result, err := s.db.Collection("tokens").DeleteOne(dbCtx, bson.M{"token_hash": token})
	if err != nil {
		return fmt.Errorf("can't delete token: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("can't delete token: %w", err)
	}

	return nil
}
