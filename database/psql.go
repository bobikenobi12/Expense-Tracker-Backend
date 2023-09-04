package database

import (
	"ExpenseTracker/models"
	"context"
	"errors"
	"os"

	"github.com/go-pg/pg/v11"
	"github.com/go-pg/pg/v11/orm"
	"github.com/gofiber/fiber/v2"
)

var PsqlDb *pg.DB

func CreateSchema(ctx context.Context) error {
	models := []interface{}{
		(*models.Expense)(nil),
		(*models.ExpenseType)(nil),
		(*models.User)(nil),
		(*models.UserSecrets)(nil),
		(*models.S3Object)(nil),
		(*models.Currency)(nil),
		(*models.Workspace)(nil),
		(*models.WorkspaceMember)(nil),
		(*models.CurrencyUser)(nil),
		(*models.CurrencyWorkspace)(nil),
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

func InsertDefaultUserEntities(c *fiber.Ctx, user *models.User) error {
	ctx := c.Context()

	defaultWorkspace := &models.Workspace{
		Name:    "Personal",
		OwnerId: user.ID,
	}

	if err := defaultWorkspace.BeforeInsert(); err != nil {
		return err
	}

	if _, err := PsqlDb.Model(defaultWorkspace).Insert(ctx); err != nil {
		return err
	}

	expenseTypes := []models.ExpenseType{}

	types := []string{
		"Food",
		"Transport",
		"Entertainment",
		"Shopping",
		"Health",
		"Other",
	}

	for _, t := range types {
		expenseType := &models.ExpenseType{
			Name: t,
		}

		if err := expenseType.BeforeInsert(); err != nil {
			return err
		}

		expenseTypes = append(expenseTypes, *expenseType)
	}

	if _, err := PsqlDb.Model(&expenseTypes).Insert(ctx); err != nil {
		return err
	}

	return nil
}

func InsertCurrencies() error {
	ctx := context.Background()

	countries := []struct {
		Name    string
		IsoCode string
	}{
		{"United States", "US"},
		{"China", "CN"},
		{"Japan", "JP"},
		{"Germany", "DE"},
		{"India", "IN"},
		{"United Kingdom", "GB"},
		{"France", "FR"},
		{"Brazil", "BR"},
		{"Italy", "IT"},
		{"Canada", "CA"},
		{"Australia", "AU"},
		{"South Korea", "KR"},
		{"Spain", "ES"},
		{"Mexico", "MX"},
		{"Indonesia", "ID"},
		{"Netherlands", "NL"},
		{"Saudi Arabia", "SA"},
		{"Turkey", "TR"},
		{"Switzerland", "CH"},
		{"Sweden", "SE"},
		{"Poland", "PL"},
		{"Belgium", "BE"},
		{"Norway", "NO"},
		{"Austria", "AT"},
		{"UAE", "AE"},
		{"Singapore", "SG"},
		{"Malaysia", "MY"},
		{"Qatar", "QA"},
		{"Thailand", "TH"},
	}

	currencies := []models.Currency{}

	for _, t := range countries {
		existingCurrency := &models.Currency{}

		if err := PsqlDb.Model(existingCurrency).Where("iso_code = ?", t.IsoCode).Select(ctx); err == nil {
			continue
		}
		currencies = append(currencies, models.Currency{
			Name:    t.Name,
			IsoCode: t.IsoCode,
		})
	}
	if len(currencies) == 0 {
		return nil
	}

	_, err := PsqlDb.Model(&currencies).Insert(ctx)

	if err != nil {
		return err
	}

	return nil
}
