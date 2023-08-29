package database

import (
	"context"
	"errors"
	"os"

	"github.com/go-pg/pg/v11"
)

func NewDbConn() (*pg.DB, error) {
	ctx := context.Background()
	user := os.Getenv("USER")
	if user == "" {
		return nil, errors.New("you must set your 'USER' environmental variable")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, errors.New("you must set your 'PASSWORD' environmental variable")
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return nil, errors.New("you must set your 'DATABASE' environmental variable")
	}

	options := &pg.Options{
		User:     user,
		Password: password,
		Database: database,
	}

	db := pg.Connect(options)
	if db == nil {
		return nil, errors.New("failed to connect to database")
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
