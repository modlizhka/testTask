package model

import "time"
import "github.com/shopspring/decimal"

type Operation struct {
	ID        int             `json:"id"`
	Type      string          `json:"operation_type"`
	Sender    int             `json:"sender"`
	Recipient int             `json:"recipient"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}

type User struct {
	ID         int             `json:"id"`
	Balance    decimal.Decimal `json:"balance"`
	Operations []Operation
}

type Payment struct {
	SenderID    int             `json:"sender"`
	RecipientID int             `json:"recipient"`
	Amount      decimal.Decimal `json:"amount"`
}
type Replenishment struct {
	RecipientID int             `json:"recipient"`
	Amount      decimal.Decimal `json:"amount"`
}
