package usecase

import (
	"context"
)

//go:generate mockgen -source=./interfaces.go -destination=./usecases_mock.go -package=usecase

type TokenInteractor interface {
	Generate(ctx context.Context, userId string) (string, string, int, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, int, error)
}
