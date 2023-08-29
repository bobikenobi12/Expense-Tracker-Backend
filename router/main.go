package router

import (
	"ExpenseTracker/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/health", handlers.Health)

	expenses := app.Group("/expenses")

	// expenses.Get("/", handlers.GetAllExpenses)
	// expenses.Get("/:id", handlers.GetExpenseByID)
	expenses.Post("/", handlers.InsertExpense)
}
