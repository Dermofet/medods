package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		logger.Error("get errors while handle request", zap.Reflect("errors", c.Errors.Errors()))
	}
}
