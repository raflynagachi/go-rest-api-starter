package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	"github.com/raflynagachi/go-rest-api-starter/internal/payment/model"
)

func (r *PostgresRepo) GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := r.DB.Rebind(`
		SELECT id, order_id, amount, status, method, created_at, created_by, updated_at, updated_by
		FROM payment_transactions WHERE id = ?
	`)
	t := &model.Transaction{}
	if err := r.DB.GetContext(ctx, t, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetTransactionByID")
		}
		return nil, errors.Wrap(err, "PostgresRepo.GetTransactionByID")
	}
	return t, nil
}

func (r *PostgresRepo) GetTransactionByOrderID(ctx context.Context, orderID int64) (*model.Transaction, error) {
	query := r.DB.Rebind(`
		SELECT id, order_id, amount, status, method, created_at, created_by, updated_at, updated_by
		FROM payment_transactions WHERE order_id = ?
	`)
	t := &model.Transaction{}
	if err := r.DB.GetContext(ctx, t, query, orderID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetTransactionByOrderID")
		}
		return nil, errors.Wrap(err, "PostgresRepo.GetTransactionByOrderID")
	}
	return t, nil
}

func (r *PostgresRepo) InsertTransaction(ctx context.Context, tx *sqlx.Tx, t *model.Transaction) (int64, error) {
	query := r.DB.Rebind(`
		INSERT INTO payment_transactions (order_id, amount, status, method, created_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`)
	var id int64
	err := tx.GetContext(ctx, &id, query, t.OrderID, t.Amount, t.Status, t.Method, t.CreatedAt, t.CreatedBy)
	return id, errors.Wrap(err, "PostgresRepo.InsertTransaction")
}

func (r *PostgresRepo) UpdateTransactionStatus(ctx context.Context, tx *sqlx.Tx, id int64, status model.TransactionStatus) error {
	query := r.DB.Rebind(`UPDATE payment_transactions SET status = ? WHERE id = ?`)
	_, err := tx.ExecContext(ctx, query, status, id)
	return errors.Wrap(err, "PostgresRepo.UpdateTransactionStatus")
}
