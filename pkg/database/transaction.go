package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TxBegin creates a new transaction
func TxBegin(db *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// TxRollback performs a rollback transaction
func TxRollback(tx *sqlx.Tx) error {
	if rbErr := tx.Rollback(); rbErr != nil {
		return fmt.Errorf("transaction rollback failed: %w", rbErr)
	}

	return nil
}

// TxCommit commits the transaction
func TxCommit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}
