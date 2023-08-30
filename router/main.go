package router

import (
	"ExpenseTracker/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/health", handlers.Health)

	expenses := app.Group("/expenses")

	expenses.Post("/types", handlers.InsertExpenseType)
	expenses.Post("/", handlers.InsertExpense)
	expenses.Get("/:id", handlers.SelectExpenseByID)
}
