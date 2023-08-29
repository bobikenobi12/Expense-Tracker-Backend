package handlers

import (
	"ExpenseTracker/database"
	"context"

	"github.com/gofiber/fiber/v2"
)

type InsertExpenseType struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
	// ExpenseTypeID int64   `json:"expense_type_id"`
}

func InsertExpense(c *fiber.Ctx) error {
	ctx := context.Background()

	expense := new(InsertExpenseType)
	if err := c.BodyParser(expense); err != nil {
		return err
	}

	_, err := database.DB.Model(expense).Returning("id").Insert(ctx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		// "data":   result.RowsReturned(),
	})
}
