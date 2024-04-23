package http

import (
	"fmt"
	"medods/internal/api/http/handlers"
	"medods/internal/db"
	"medods/internal/repository"
	"medods/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
)

type routerHandlers struct {
	authHandlers handlers.AuthHandlers
}

type router struct {
	router   *gin.Engine
	db       *mongo.Database
	handlers routerHandlers
	logger   *zap.Logger
}

func NewRouter(db *mongo.Database, logger *zap.Logger) *router {
	return &router{
		router: gin.New(),
		db:     db,
		logger: logger,
	}
}

func (r *router) Init() error {
	r.router.Use(
		gin.Logger(),
		gin.CustomRecovery(r.recovery),
	)
	err := r.registerRoutes()
	if err != nil {
		return fmt.Errorf("can't init router: %w", err)
	}

	return nil
}

func (r *router) recovery(c *gin.Context, recovered any) {
	defer func() {
		if e := recover(); e != nil {
			r.logger.Fatal("http server panic", zap.Error(fmt.Errorf("%s", recovered)))
		}
	}()
	c.AbortWithStatus(http.StatusInternalServerError)
}

func (r *router) registerRoutes() error {
	r.router.NoMethod(handlers.NotImplementedHandler)
	r.router.NoRoute(handlers.NotImplementedHandler)

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r.router.Use(corsMiddleware)

	pgSource := db.NewSource(r.db)
	tokenRepository := repository.NewTokenRepository(pgSource)
	tokenInteractor := usecase.NewTokenInteractor(tokenRepository)
	r.handlers.authHandlers = handlers.NewAuthHandlers(tokenInteractor)

	authGroup := r.router.Group("/auth")
	authGroup.POST("/new-tokens", r.handlers.authHandlers.GenerateTokens)
	authGroup.POST("/refresh-tokens", r.handlers.authHandlers.RefreshTokens)

	return nil
}
