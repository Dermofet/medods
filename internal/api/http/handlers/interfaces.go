package handlers

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=./interfaces.go -destination=./handlers_mock.go -package=handlers

type AuthHandlers interface {
	GenerateTokens(c *gin.Context)
	RefreshTokens(c *gin.Context)
}
