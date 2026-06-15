package definition

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/raflynagachi/go-rest-api-starter/internal/payment/model"
)

type SQLRepo interface {
	Transaction

	GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error)
	GetTransactionByOrderID(ctx context.Context, orderID int64) (*model.Transaction, error)
	InsertTransaction(ctx context.Context, tx *sqlx.Tx, t *model.Transaction) (int64, error)
	UpdateTransactionStatus(ctx context.Context, tx *sqlx.Tx, id int64, status model.TransactionStatus) error
}

type Transaction interface {
	TxBegin() (*sqlx.Tx, error)
	TxEnd(tx *sqlx.Tx, err error) error
}
