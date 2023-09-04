package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	if err := database.CheckIfEmailExists(signUp.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
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

	defaultWorkspace := &models.Workspace{
		Name:    "Personal",
		OwnerId: user.ID,
	}

	if err := defaultWorkspace.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(defaultWorkspace).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
		"workspaces": []models.Workspace{
			*defaultWorkspace,
		},
	})
}

func LoginHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	login := &config.LoginRequest{}

	if err := c.BodyParser(login); err != nil {
		return err
	}

	if err := config.ValidationResponse(login); err != nil {
		return err
	}

	user := &models.User{}

	if err := database.PsqlDb.Model(user).Where("email = ?", login.Email).Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "A user with this email does not exist",
		})
	}

	userSecrets := &models.UserSecrets{
		ID: user.UserSecretsId,
	}

	if err := database.PsqlDb.Model(userSecrets).WherePK().Select(ctx); err != nil {
		return err
	}

	match := config.CheckPasswordHash(login.Password, userSecrets.Password)

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Incorrect password",
		})
	}

	config.SetJwtsToCookies(c, &jwt.MapClaims{
		"email":       user.Email,
		"name":        user.Name,
		"id":          user.ID,
		"prof_pic_id": user.ProfilePicId,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"msg":    "Logged in",
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	config.BlacklistJwt(c)

	c.ClearCookie()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"msg":    "Logged out",
	})
}
