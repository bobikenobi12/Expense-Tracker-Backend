package config

import (
	"ExpenseTracker/database"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	Filter      func(c *fiber.Ctx) bool
	Unathorized fiber.Handler
	Decode      func(c *fiber.Ctx) (*jwt.MapClaims, error)
	Secret      string
	Expiry      int64
}

var JwtConfigDefault = JwtConfig{
	Filter:      nil,
	Decode:      nil,
	Unathorized: nil,
	Secret:      "secret",
	Expiry:      60 * 15,
}

type CustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func New(config JwtConfig) fiber.Handler {
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {

		if cfg.Filter != nil && !cfg.Filter(c) {
			fmt.Println("Middleware skipped")
			return c.Next()
		}
		fmt.Println("Middleware executed")

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("jwtClaims", *claims)
			return c.Next()
		}

		return cfg.Unathorized(c)
	}
}

func configDefault(config ...JwtConfig) JwtConfig {
	if len(config) < 1 {
		return JwtConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = JwtConfigDefault.Filter
	}

	if cfg.Secret == "" {
		cfg.Secret = JwtConfigDefault.Secret
	}

	if cfg.Expiry == 0 {
		cfg.Expiry = JwtConfigDefault.Expiry
	}

	if cfg.Decode == nil {
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {

			cookieToken := c.Cookies("token")

			log.Println("cookieToken", cookieToken)

			if cookieToken == "" {
				return nil, errors.New("missing auth token")
			}

			token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method %v" + token.Header["alg"].(string))
				}
				return []byte(cfg.Secret), nil
			},
			)

			if err != nil {
				return nil, errors.New("error parsing token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				return nil, errors.New("invalid token")
			}

			if expriresAt, ok := claims["exp"]; ok && int64(expriresAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("token is expired")
			}

			return &claims, nil
		}
	}

	if cfg.Unathorized == nil {
		cfg.Unathorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

func SetJwtsToCookies(c *fiber.Ctx, claims *jwt.MapClaims) {
	tokenJwt, err := Encode(claims, 60*15)

	if err != nil {
		GlobalErrorHandler(c, c.SendStatus(fiber.StatusInternalServerError))
	}

	refreshJwt, err := Encode(claims, 60*60*24*7)

	if err != nil {
		GlobalErrorHandler(c, c.SendStatus(fiber.StatusInternalServerError))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenJwt,
		Expires:  time.Now().Add(time.Minute * 15).UTC(),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshJwt,
		Expires:  time.Now().Add(time.Hour * 24 * 7).UTC(),
		HTTPOnly: true,
	})
}

func BlacklistJwt(c *fiber.Ctx) {
	token := c.Cookies("token")

	if token != "" {
		database.RedisClient.Set(c.Context(), token, "blacklisted", time.Minute*15)
	}
}

func Encode(claims *jwt.MapClaims, expiryAfter int64) (string, error) {

	if expiryAfter == 0 {
		expiryAfter = JwtConfigDefault.Expiry
	}

	(*claims)["exp"] = time.Now().UTC().Unix() + expiryAfter

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedTokens, err := token.SignedString([]byte(JwtConfigDefault.Secret))

	if err != nil {
		return "", errors.New("error signing token")
	}

	return signedTokens, nil
}
