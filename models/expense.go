package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExpenseType struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

type Expense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ExpenseType primitive.ObjectID `bson:"expense_type" json:"expense_type"`
	Amount      float64            `bson:"amount" json:"amount"`
	Note        string             `bson:"note" json:"note"`
	Date        string             `bson:"date" json:"date"`
}
