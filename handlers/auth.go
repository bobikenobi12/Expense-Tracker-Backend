package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SignUpHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	signUp := &config.SignUpRequest{}

	if err := c.BodyParser(signUp); err != nil {
		return err
	}

	if err := config.ValidationResponse(signUp); err != nil {
		return err
	}

	user := &models.User{
		Email:       signUp.Email,
		Name:        signUp.Name,
		CountryCode: signUp.CountryCode,
	}

	userSecrets := &models.UserSecrets{
		Password: signUp.Password,
	}

	if err := userSecrets.BeforeInsert(); err != nil {
		log.Println(userSecrets)
		return err
	}

	if _, err := database.DB.Model(userSecrets).Insert(ctx); err != nil {
		return err
	}

	log.Println(userSecrets)

	user.UserSecretsId = userSecrets.ID

	if err := user.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.DB.Model(user).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
