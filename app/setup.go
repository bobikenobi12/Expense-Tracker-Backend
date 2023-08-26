package app

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
)

func SetupAndRunApp() error {
	err := config.LoadENV()
	if err != nil {
		return err
	}

	err = database.StartMongoDb()
	if err != nil {
		return err
	}

	defer database.CloseMongoDb()

	return nil
}
