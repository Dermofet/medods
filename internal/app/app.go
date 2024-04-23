package app

import (
	"context"
	"fmt"
	"medods/cmd/medods/config"
	"medods/internal/api/http"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type app struct {
	config     *config.Config
	dbClient   *mongo.Client
	logger     *zap.Logger
	httpServer http.Server
}

func NewApp(cfg *config.Config, logger *zap.Logger) *app {
	return &app{
		config: cfg,
		logger: logger,
	}
}

func (a *app) Start(ctx context.Context) {
	appCtx, cancelApp := context.WithCancel(ctx)
	defer func() {
		if e := recover(); e != nil {
			a.logger.Panic("application shutdown", zap.Error(fmt.Errorf("%s", e)))
			cancelApp()
		}
	}()
	// Инициализируем БД
	dbClient, err := a.initDb(appCtx,
		a.config.DB.Host,
		a.config.DB.Port,
		a.config.DB.Name,
		a.config.DB.Username,
		a.config.DB.Password,
	)
	if err != nil {
		a.logger.Fatal("init db error", zap.Error(err))
	}
	a.dbClient = dbClient
	db := a.dbClient.Database(a.config.DB.Name)

	// Запуск миграций
	err = a.startMigrate(appCtx, migrationsPath, a.config.DB.Name, a.dbClient)
	if err != nil {
		a.logger.Error("db migration error", zap.Error(err))
	}

	wg := &sync.WaitGroup{}
	// Старт HTTP-сервера
	wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				a.logger.Panic("http start panic", zap.Error(fmt.Errorf("%s", e)))
			}
			wg.Done()
		}()
		addr := fmt.Sprintf("%s:%d", a.config.HttpServer.Host, a.config.HttpServer.Port)
		a.httpServer = http.NewServer(addr, db, a.logger)
		if a.httpServer == nil {
			cancelApp()
			a.logger.Fatal("can't create http server")
			return
		}
		err := a.httpServer.Run(appCtx)
		// Отменяем контекст, если HTTP-сервер завершил работу
		cancelApp()
		if err != nil {
			a.logger.Error("can't start http server", zap.Error(err))
			return
		}
	}()

	wg.Wait()
}

// GracefulShutdown graceful shutdown приложения
func (a *app) GracefulShutdown(ctx context.Context) error {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("can't shutdown http-server: %w", err)
	}

	err = a.dbClient.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("can't shutdown db: %w", err)
	}
	return nil
}

// initDb инициализация базы данных
func (a *app) initDb(
	ctx context.Context,
	host string,
	port int,
	dbName string,
	user string,
	password string,
) (*mongo.Client, error) {
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", user, password, host, port, dbName)))
	if err != nil {
		return nil, err
	}

	return db, nil
}
