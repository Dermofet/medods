package app

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.mongodb.org/mongo-driver/mongo"
)

const migrationsPath = "migrations"

//go:embed  migrations/*.json
var fs embed.FS

func (a *app) startMigrate(ctx context.Context, migratePath string, dbName string, dbClient *mongo.Client) error {
	err := dbClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("db connection not alive: %w", err)
	}
	driver, err := mongodb.WithInstance(dbClient, &mongodb.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("db migration database driver error: %w", err)
	}
	source, err := iofs.New(fs, migratePath)
	if err != nil {
		return fmt.Errorf("db migration source driver error: %w", err)
	}
	instance, err := migrate.NewWithInstance("fs", source, dbName, driver)
	if err != nil {
		return fmt.Errorf("db migration instance error: %w", err)
	}
	if err := instance.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("db migration up error: %w", err)
	}

	return nil
}
