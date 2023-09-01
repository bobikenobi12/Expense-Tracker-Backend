package config

import (
	"ExpenseTracker/database"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func ValidateJwt(c *fiber.Ctx) error {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "secret"
	}

	token, err := jwt.ParseWithClaims(c.Cookies("token"), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}, jwt.WithLeeway(5*time.Minute))

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		c.Locals("claims", claims)
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
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

func BlacklistJwt(c *fiber.Ctx) {
	token := c.Cookies("token")

	if token != "" {
		database.RedisClient.Set(c.Context(), token, "blacklisted", time.Minute*15)
	}
}

func GenerateJWT(email string, name string, experation time.Time) (string, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "secret"
	}

	claims := CustomClaims{
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
