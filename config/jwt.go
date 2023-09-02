package config

import (
	"ExpenseTracker/database"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type JwtConfig struct {
	Filter      func(c *fiber.Ctx) bool
	Unathorized func(c *fiber.Ctx, err error) error
	Decode      func(c *fiber.Ctx, cookieType string) (*jwt.MapClaims, error)
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
	Id    string `json:"id"`
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

		tokenClaims, err := cfg.Decode(c, "token")

		if err == nil {
			c.Locals("jwtClaims", *tokenClaims)
			return c.Next()
		}

		refreshTokenClaims, err := cfg.Decode(c, "refresh_token")

		if err == nil {
			newToken, err := Encode(refreshTokenClaims, 60*15)

			if err != nil {
				return cfg.Unathorized(c, err)
			}

			blacklisted, err := database.RedisClient.Get(c.Context(), newToken).Result()

			if err != nil {
				if err == redis.Nil {
					log.Println("key does not exist")
				} else {
					return cfg.Unathorized(c, errors.New("error checking if token is blacklisted"))
				}
			}

			if blacklisted == "blacklisted" {
				return cfg.Unathorized(c, errors.New("token is blacklisted"))
			}

			c.Cookie(&fiber.Cookie{
				Name:     "token",
				Value:    newToken,
				Expires:  time.Now().Add(time.Minute * 15).UTC(),
				HTTPOnly: true,
			})

			c.Locals("jwtClaims", *refreshTokenClaims)

			return c.Next()
		}

		return cfg.Unathorized(c, err)
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
		cfg.Decode = func(c *fiber.Ctx, cookieType string) (*jwt.MapClaims, error) {

			cookieToken := c.Cookies(cookieType)

			if cookieToken == "" {
				return nil, errors.New("missing auth token")
			}

			blacklisted, err := database.RedisClient.Get(c.Context(), cookieToken).Result()

			if err != nil {
				if err == redis.Nil {
					log.Println("key does not exist")

				} else {
					return nil, errors.New("error checking if token is blacklisted")
				}
			}

			if blacklisted == "blacklisted" {
				return nil, errors.New("token is blacklisted")
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
		cfg.Unathorized = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
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
		database.RedisClient.Set(c.Context(), token, "blacklisted", time.Minute)
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
