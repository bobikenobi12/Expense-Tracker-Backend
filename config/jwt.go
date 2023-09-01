package config

import (
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func IssueJwt(email string, name string) (string, error) {
	expiration := time.Now().Add(time.Minute * 15)

	return GenerateJWT(email, name, expiration)
}

func RefreshJwt(email string, name string) (string, error) {
	expiration := time.Now().Add(time.Hour * 24 * 7)

	return GenerateJWT(email, name, expiration)
}

func GenerateJWT(email string, name string, experation time.Time) (string, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "secret"
	}

	claims := Claims{
		Email: email,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experation),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ExpenseTracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func WhitelistJwt(c *fiber.Ctx, jwt string, expiration time.Time) error {
	ctx := c.Context()

	jwtWhitelist := &models.JwtWhitelist{
		Jwt:       jwt,
		ExpiresAt: expiration,
	}

	if err := jwtWhitelist.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(jwtWhitelist).Insert(ctx); err != nil {
		return err
	}

	return nil
}

func BlacklistJwt(c *fiber.Ctx, jwt string) error {
	ctx := c.Context()

	jwtBlacklist := &models.JwtBlacklist{
		Jwt: jwt,
	}

	if _, err := database.PsqlDb.Model(jwtBlacklist).Insert(ctx); err != nil {
		return err
	}

	return nil
}
