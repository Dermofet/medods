package usecase

import (
	"context"
	"fmt"
	"medods/internal/entity"
	"medods/internal/repository"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type tokenInteractor struct {
	tokenRepo repository.TokenRepository
}

func NewTokenInteractor(tokenRepo repository.TokenRepository) *tokenInteractor {
	return &tokenInteractor{
		tokenRepo: tokenRepo,
	}
}

func (t *tokenInteractor) Generate(ctx context.Context, userIdStr string) (string, string, int, error) {
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return "", "", http.StatusUnprocessableEntity, fmt.Errorf("guid is not valid: %v", err)
	}

	refreshToken, err := entity.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate refresh token: %v", err)
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't hash refresh token: %v", err)
	}

	t.tokenRepo.Create(ctx, &entity.Auth{
		UserID: userId,
		Token:  refreshTokenHash,
	})

	accessToken, err := entity.GenerateAccessToken(userId).String()
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate access token: %v", err)
	}

	return accessToken, refreshToken, http.StatusOK, nil
}

func (t *tokenInteractor) Refresh(ctx context.Context, refreshToken string) (string, string, int, error) {
	auth, err := t.tokenRepo.GetUserId(ctx, []byte(refreshToken))
	if err != nil {
		return "", "", http.StatusConflict, fmt.Errorf("can't find refresh token in usecase: %v", err)
	}

	err = t.tokenRepo.Delete(ctx, auth.Token)
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't delete refresh token: %v", err)
	}

	refreshToken, err = entity.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate refresh token: %v", err)
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate refresh token: %v", err)
	}

	auth.Token = refreshTokenHash

	err = t.tokenRepo.Create(ctx, auth)
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate refresh token: %v", err)
	}

	accessToken, err := entity.GenerateAccessToken(auth.UserID).String()
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("can't generate access token: %v", err)
	}

	return accessToken, refreshToken, http.StatusOK, nil
}
