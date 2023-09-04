package models

import (
	"time"
)

type ExpenseType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}

type Expense struct {
	ID            int64        `json:"id"`
	Amount        float64      `json:"amount"`
	Note          string       `json:"note"`
	Date          time.Time    `json:"date" pg:"default:now()"`
	ExpenseType   *ExpenseType `pg:"rel:has-one" json:"expense_type"`
	ExpenseTypeID int64        `json:"expense_type_id"`
	WorkspaceID   int64        `json:"workspace_id"`
	CurrencyId    int64        `json:"currency_id"`
}

func (e *Expense) BeforeInsert() error {
	e.Date = time.Now()
	return nil
}

func (et *ExpenseType) BeforeInsert() error {
	et.CreatedAt = time.Now()
	et.UpdatedAt = time.Now()
	return nil
}
