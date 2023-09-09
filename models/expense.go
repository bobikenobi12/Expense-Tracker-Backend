package models

import (
	"time"
)

type ExpenseType struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}

type Expense struct {
	ID            uint64       `json:"id"`
	Amount        float64      `json:"amount"`
	Note          string       `json:"note"`
	Date          time.Time    `json:"date" pg:"default:now()"`
	ExpenseType   *ExpenseType `pg:"rel:has-one" json:"expense_type"`
	ExpenseTypeID uint64       `json:"expense_type_id"`
	WorkspaceID   uint64       `json:"workspace_id"`
	CurrencyId    uint64       `json:"currency_id"`
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
