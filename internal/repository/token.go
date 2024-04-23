package repository

import (
	"context"
	"fmt"
	"medods/internal/db"
	"medods/internal/entity"
)

type tokenRepository struct {
	source db.TokenSource
}

func NewTokenRepository(source db.TokenSource) *tokenRepository {
	return &tokenRepository{
		source: source,
	}
}

func (r *tokenRepository) Create(ctx context.Context, auth *entity.Auth) error {
	err := r.source.CreateToken(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) GetUserId(ctx context.Context, token []byte) (*entity.Auth, error) {
	auth, err := r.source.GetUserIdFromToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("can't find token in repo: %w", err)
	}

	return auth, nil
}

func (r *tokenRepository) Delete(ctx context.Context, token []byte) error {
	err := r.source.DeleteToken(ctx, token)
	if err != nil {
		return fmt.Errorf("can't delete token: %w", err)
	}
	return nil
}
