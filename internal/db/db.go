package db

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	QueryTimeout = 10 * time.Second
)

type source struct {
	db *mongo.Database
}

func NewSource(db *mongo.Database) *source {
	return &source{
		db: db,
	}
}
