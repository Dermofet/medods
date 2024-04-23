package repository

import (
	"context"
	"medods/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=repository

type TokenRepository interface {
	Create(ctx context.Context, auth *entity.Auth) error
	GetUserId(ctx context.Context, token []byte) (*entity.Auth, error)
	Delete(ctx context.Context, token []byte) error
}
