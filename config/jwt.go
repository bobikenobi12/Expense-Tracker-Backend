package config

import (
	"os"
	"time"

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
