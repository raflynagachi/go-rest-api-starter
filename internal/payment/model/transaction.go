package model

import sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusSucceeded TransactionStatus = "succeeded"
	StatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID      int64             `db:"id"`
	OrderID int64             `db:"order_id"`
	Amount  float64           `db:"amount"`
	Status  TransactionStatus `db:"status"`
	Method  string            `db:"method"`
	sharedmodel.Created
	sharedmodel.Updated
}
