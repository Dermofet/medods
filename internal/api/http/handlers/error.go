package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotImplementedHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}
