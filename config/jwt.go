package config

import (
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

func SetJwtsToCookies(c *fiber.Ctx, email string, name string) {
	tokenJwt, err := GenerateJWT(email, name, time.Now().Add(time.Minute*15))

	if err != nil {
		GlobalErrorHandler(c, err)
	}

	refreshJwt, err := GenerateJWT(email, name, time.Now().Add(time.Hour*24*7))

	if err != nil {
		GlobalErrorHandler(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenJwt,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshJwt,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})
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
