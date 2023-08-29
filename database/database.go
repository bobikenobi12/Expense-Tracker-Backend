package database

import (
	"context"
	"errors"
	"os"

	"github.com/go-pg/pg/v11"
	"github.com/go-pg/pg/v11/orm"
)

var DB *pg.DB

func GetModel(model interface{}) *orm.Query {
	return DB.Model(model)
}

func NewDbConn() error {
	ctx := context.Background()

	user := os.Getenv("USER")
	if user == "" {
		return errors.New("you must set your 'USER' environmental variable")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return errors.New("you must set your 'PASSWORD' environmental variable")
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return errors.New("you must set your 'DATABASE' environmental variable")
	}

	options := &pg.Options{
		User:     user,
		Password: password,
		Database: database,
	}

	DB = pg.Connect(options)

	if err := DB.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func CloseConn() {
	ctx := context.Background()
	DB.Close(ctx)
}
