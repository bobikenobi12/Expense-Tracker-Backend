package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func InitStore() {
	store = session.New()
}

func GetStore(c *fiber.Ctx) *session.Store {
	return store
}
