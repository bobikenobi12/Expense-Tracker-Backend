package database

import (
	"ExpenseTracker/models"
	"context"
	"errors"
	"os"

	"github.com/go-pg/pg/v11"
	"github.com/go-pg/pg/v11/orm"
)

var PsqlDb *pg.DB

func CreateSchema(ctx context.Context) error {
	models := []interface{}{
		(*models.Expense)(nil),
		(*models.ExpenseType)(nil),
		(*models.User)(nil),
		(*models.UserSecrets)(nil),
		(*models.S3Object)(nil),
	}

	for _, model := range models {
		err := PsqlDb.Model(model).CreateTable(ctx, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func NewPsqlDbConn() error {
	ctx := context.Background()

	user := os.Getenv("PSQL_USER")
	if user == "" {
		return errors.New("you must set your 'USER' environmental variable")
	}

	password := os.Getenv("PSQL_PASSWORD")
	if password == "" {
		return errors.New("you must set your 'PASSWORD' environmental variable")
	}

	database := os.Getenv("PSQL_DATABASE")
	if database == "" {
		return errors.New("you must set your 'DATABASE' environmental variable")
	}

	options := &pg.Options{
		User:     user,
		Password: password,
		Database: database,
	}

	PsqlDb = pg.Connect(options)

	if err := PsqlDb.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func ClosePsqlConn() {
	ctx := context.Background()
	PsqlDb.Close(ctx)
}

func CheckIfEmailExists(email string) error {
	ctx := context.Background()

	user := &models.User{}

	err := PsqlDb.Model(user).Where("email = ?", email).Select(ctx)
	switch err {
	case pg.ErrNoRows:
		return nil
	case pg.ErrMultiRows:
		return errors.New("a user with this email already exists")
	case nil:
		return errors.New("a user with this email already exists")
	}

	return err
}
