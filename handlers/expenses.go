package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func InsertExpenseType(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("jwtClaims").(jwt.MapClaims)

	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}
	expenseTypeReq := &config.InsertExpenseTypeRequest{}

	if err := c.BodyParser(expenseTypeReq); err != nil {
		return err
	}

	if err := config.ValidationResponse(expenseTypeReq); err != nil {
		return err
	}

	expenseType := &models.ExpenseType{
		Name: expenseTypeReq.Name,
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

	claims := c.Locals("jwtClaims").(jwt.MapClaims)

	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	expenseReq := &config.InsertExpenseRequest{}

	if err := c.BodyParser(expenseReq); err != nil {
		return err
	}

	if err := config.ValidationResponse(expenseReq); err != nil {
		return err
	}

	if err := database.PsqlDb.Model(&models.ExpenseType{ID: expenseReq.TypeId}).WherePK().Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid expense type id",
		})
	}

	if err := database.PsqlDb.Model(&models.Workspace{ID: expenseReq.WorkspaceId}).WherePK().Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid workspace id",
		})
	}

	if err := database.PsqlDb.Model(&models.Currency{ID: expenseReq.CurrencyId}).WherePK().Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid currency id",
		})
	}

	expense := &models.Expense{
		Amount:        expenseReq.Amount,
		Note:          expenseReq.Note,
		ExpenseTypeID: expenseReq.TypeId,
		WorkspaceID:   expenseReq.WorkspaceId,
		CurrencyId:    expenseReq.CurrencyId,
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

	claims := c.Locals("jwtClaims").(jwt.MapClaims)

	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claims["id"].(float64)

	expenseById := config.GetExpenseByIdRequest{}

	if err := c.ParamsParser(expenseById); err != nil {
		return err
	}

	if err := config.ValidationResponse(expenseById); err != nil {
		return err
	}

	expense := &models.Expense{}

	if err := database.PsqlDb.Model(expense).Relation("ExpenseType").Where("expense.id = ?", expenseById.Id).Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid expense id",
		})
	}

	if err := database.UserCanFetchExpense(c, uint64(userId), uint64(expense.WorkspaceID)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"expense": expense,
	})
}

func SelectExpenses(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("jwtClaims").(jwt.MapClaims)

	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claims["id"].(float64)

	paginationReq := config.PaginationRequest{}

	if err := c.QueryParser(&paginationReq); err != nil {
		return err
	}

	if err := config.ValidationResponse(paginationReq); err != nil {
		return err
	}

	if err := database.UserCanFetchExpense(c, uint64(userId), paginationReq.WorkspaceId); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	expenses := make([]*models.Expense, 0)

	if err := database.PsqlDb.Model(&expenses).Relation("ExpenseType").Offset(int((paginationReq.Page - 1) * paginationReq.Size)).Limit(int(paginationReq.Size)).Select(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"expenses": expenses,
	})
}
