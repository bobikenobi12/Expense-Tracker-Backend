package database

import (
	"ExpenseTracker/models"
	"context"
	"errors"
	"os"

	"github.com/go-pg/pg/v11"
	"github.com/go-pg/pg/v11/orm"
)

var DB *pg.DB

func CreateSchema(ctx context.Context) error {
	models := []interface{}{
		(*models.Expense)(nil),
		(*models.ExpenseType)(nil),
		(*models.User)(nil),
		(*models.UserSecrets)(nil),
	}

	for _, model := range models {
		err := DB.Model(model).CreateTable(ctx, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
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
