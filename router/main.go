package router

import (
	"ExpenseTracker/config"
	"ExpenseTracker/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/health", handlers.Health)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/sign-up", handlers.SignUpHandler)
	auth.Post("/login", handlers.LoginHandler)

	app.Use(config.New(config.JwtConfig{}))

	auth.Get("/logout", handlers.LogoutHandler)

	expenses := api.Group("/expenses")
	expenses.Post("/types", handlers.InsertExpenseType)
	expenses.Post("/", handlers.InsertExpense)
	expenses.Get("/:id", handlers.SelectExpenseByID)
}
