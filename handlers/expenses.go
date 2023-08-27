package handlers

import (
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetExpenses godoc
// @Summary Get all expenses
// @Description Get all expenses
// @Tags Expenses
// @Accept json
// @Produce json
// @Success 200 {array} Expense
// @Router /expenses [get]
func GetExpenses(c *fiber.Ctx) error {
	expenses := database.GetCollection("expenses")

	var results []models.Expense

	cursor, err := expenses.Find(c.Context(), nil)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("No expenses found")
		}
		panic(err)
	}

	if err = cursor.All(c.Context(), &results); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving expenses")
	}

	return c.JSON(results)
}

func AddExpense(c *fiber.Ctx) error {
	expenses := database.GetCollection("expenses")

	expense := new(models.Expense)

	if err := c.BodyParser(expense); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing expense")

	}

	result, err := expenses.InsertOne(c.Context(), expense)

	if err != nil {
		log.Println(err)
		c.Status(fiber.StatusInternalServerError).SendString("Error adding expense")
		return err
	}

	c.JSON(result)

	return nil
}
