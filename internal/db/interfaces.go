package db

import (
	"context"
	"medods/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./source_mock.go -package=db

type TokenSource interface {
	CreateToken(ctx context.Context, auth *entity.Auth) error
	GetUserIdFromToken(ctx context.Context, token []byte) (*entity.Auth, error)
	DeleteToken(ctx context.Context, token []byte) error
}
