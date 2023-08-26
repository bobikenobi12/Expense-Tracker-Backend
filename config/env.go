package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadENV() error {
	goEnv := os.Getenv("GO_ENV")

	if goEnv == "" || goEnv == "development" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	return nil
}
