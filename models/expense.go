package models

import (
	"fmt"
	"time"
)

type ExpenseType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
}

type Expense struct {
	ID            int64        `json:"id"`
	Amount        float64      `json:"amount"`
	Note          string       `json:"note"`
	Date          time.Time    `json:"date" pg:"default:now()"`
	ExpenseType   *ExpenseType `pg:"rel:has-one" json:"expense_type"`
	ExpenseTypeID int64        `json:"expense_type_id"`
}

func (et *ExpenseType) PrintExpenseType() string {
	return fmt.Sprintf("ExpenseType<%d %s>", et.ID, et.Name)
}

func (e *Expense) PrintExpense() string {
	return fmt.Sprintf("Expense<%d %f %s %s %v %d>", e.ID, e.Amount, e.Note, e.Date, e.ExpenseType, e.ExpenseTypeID)
}

func (e *Expense) BeforeInsert() error {
	e.Date = time.Now()
	return nil
}

func (et *ExpenseType) BeforeInsert() error {
	et.CreatedAt = time.Now()
	return nil
}
