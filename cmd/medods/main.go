package main

import (
	"context"
	"fmt"
	"log"
	"medods/cmd/medods/config"
	"medods/internal/app"
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type appVersion struct {
	name    string
	version string
}

func (v *appVersion) SetName(name string) {
	v.name = name
}

func (v *appVersion) SetVersion(version string) {
	v.version = version
}

func (v *appVersion) GetRelease() string {
	return fmt.Sprintf("%s@%s", v.name, v.version)
}

func (v *appVersion) LoadFromConfig(cfg *config.Config) {
	v.name = cfg.AppInfo.Name
	v.version = cfg.AppInfo.Version
}

var AppVersion *appVersion

func main() {
	// Parse the application configuration
	cfg, err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("can't parse app config: %v", err)
	}

	AppVersion = &appVersion{}
	AppVersion.LoadFromConfig(cfg)

	fmt.Println("Parsed config:", cfg)

	// Initialize the logger
	logConfig := zap.NewProductionConfig()
	logConfig.Development = cfg.Debug
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("invalid log level: %v", err)
		return
	}
	logConfig.Level = zap.NewAtomicLevelAt(level)
	logConfig.OutputPaths = []string{cfg.PathLog}

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("can't create logger: %v", err)
		return
	}

	defer logger.Sync()

	defer func() {
		if e := recover(); e != nil {
			logger.Fatal("panic error", zap.Error(fmt.Errorf("%s", e)))
		}
	}()

	wg := &sync.WaitGroup{}
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	application := app.NewApp(cfg, logger)
	logger.Info("starting application", zap.String("version", AppVersion.GetRelease()))
	// Запуск приложения
	wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Panic("application start panic", zap.Error(fmt.Errorf("%s", e)))
			}
			wg.Done()
		}()
		application.Start(ctx)
	}()

	// Ожидание завершения контекста для graceful shutdown
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		<-ctx.Done()
		err := application.GracefulShutdown(ctx)
		if err != nil {
			logger.Fatal("graceful shutdown error", zap.Error(err))
		}
	}()
	wg.Wait()

	logger.Warn("application is shutdown")
}
