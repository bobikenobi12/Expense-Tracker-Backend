package models

import (
	"fmt"
	"time"
)

type ExpenseType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Expense struct {
	ID            int64        `json:"id"`
	Amount        float64      `json:"amount"`
	Note          string       `json:"description"`
	Date          time.Time    `json:"date"`
	ExpenseType   *ExpenseType `pg:"rel:has-one" json:"expense_type"`
	ExpenseTypeID int64        `json:"expense_type_id"`
}

func (et *ExpenseType) PrintExpenseType() string {
	return fmt.Sprintf("ExpenseType<%d %s>", et.ID, et.Name)
}

func (e *Expense) PrintExpense() string {
	return fmt.Sprintf("Expense<%d %f %s %s %v>", e.ID, e.Amount, e.Note, e.Date, e.ExpenseType)
}
