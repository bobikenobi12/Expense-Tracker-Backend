package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"

	"github.com/gofiber/fiber/v2"
)

func SignUpHandler(c *fiber.Ctx) error {
	ctx := c.Context()

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

	hashedPassword, err := config.HashPassword(signUp.Password)
	if err != nil {
		return err
	}

	userSecrets := &models.UserSecrets{
		Password: hashedPassword,
	}

	if err := userSecrets.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(userSecrets).Insert(ctx); err != nil {
		return err
	}

	user.UserSecretsId = userSecrets.ID

	if err := user.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(user).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
