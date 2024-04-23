package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"medods/internal/entity"
	"medods/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandlers struct {
	interactor usecase.TokenInteractor
}

func NewAuthHandlers(interactor usecase.TokenInteractor) *authHandlers {
	return &authHandlers{
		interactor: interactor,
	}
}

func (a *authHandlers) GenerateTokens(c *gin.Context) {
	guid, ok := c.GetQuery("guid")
	if !ok {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't get guid from query"))
		return
	}

	accessToken, refreshToken, status, err := a.interactor.Generate(context.Background(), guid)

	if err != nil {
		c.AbortWithError(status, err)
		return
	}

	if accessToken != "" && refreshToken != "" {
		c.JSON(http.StatusOK, entity.TokensView{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
		return
	}

	c.AbortWithStatus(status)
}

func (a *authHandlers) RefreshTokens(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't get data: %v", err))
		return
	}

	var rt entity.RefreshToken
	if err := json.Unmarshal(data, &rt); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal: %v", err))
		return
	}

	accessToken, refreshToken, status, err := a.interactor.Refresh(context.Background(), rt.RefreshToken)

	if err != nil {
		c.AbortWithError(status, err)
		return
	}

	if accessToken != "" && refreshToken != "" {
		c.JSON(http.StatusOK, entity.TokensView{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
		return
	}

	c.AbortWithStatus(status)
}
