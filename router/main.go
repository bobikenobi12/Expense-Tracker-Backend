package router

import (
	"ExpenseTracker/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/health", handlers.Health)

	expenses := app.Group("/expenses")

	expenses.Get("/", handlers.GetExpenses)
	expenses.Post("/", handlers.AddExpense)

}
