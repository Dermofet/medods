package entity

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"medods/cmd/medods/config"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	AccessExpiresTime  = time.Hour * 24
	RefreshExpiresTime = AccessExpiresTime * 7
)

type Token struct {
	Token *jwt.Token
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type Auth struct {
	UserID uuid.UUID `bson:"user_id"`
	Token  []byte    `bson:"token_hash"`
}

type AccessClaims struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type TokensView struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (t *Token) String() (string, error) {
	cfg, err := config.GetAppConfig()
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}

	return t.Token.SignedString([]byte(cfg.ApiKey))
}

func GenerateAccessToken(userId uuid.UUID) *Token {
	return &Token{
		Token: jwt.NewWithClaims(
			jwt.SigningMethodHS512,
			AccessClaims{
				UserID: userId,
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(AccessExpiresTime).Unix(),
					Subject:   "auth",
				},
			},
		),
	}
}

func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("can't generate refresh token: %w", err)
	}
	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}
