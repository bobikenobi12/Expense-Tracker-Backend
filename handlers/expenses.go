package handlers

import (
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"context"

	"github.com/gofiber/fiber/v2"
)

func InsertExpenseType(c *fiber.Ctx) error {
	ctx := c.Context()

	expenseType := &models.ExpenseType{}
	if err := c.BodyParser(expenseType); err != nil {
		return err
	}

	if err := expenseType.BeforeInsert(); err != nil {
		return err
	}

	result, err := database.PsqlDb.Model(expenseType).Returning("id").Insert(ctx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   result.RowsReturned(),
	})
}

func InsertExpense(c *fiber.Ctx) error {
	ctx := context.Background()

	expense := &models.Expense{}

	if err := c.BodyParser(expense); err != nil {
		return err
	}

	if err := expense.BeforeInsert(); err != nil {
		return err
	}

	result, err := database.PsqlDb.Model(expense).Returning("id").Insert(ctx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"expense": result.RowsReturned(),
	})
}

func SelectExpenseByID(c *fiber.Ctx) error {
	ctx := c.Context()

	expense := &models.Expense{}

	id := c.Params("id")

	if err := database.PsqlDb.Model(expense).Relation("ExpenseType").Where("expense.id = ?", id).Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid expense id",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"expense": expense,
	})
}
