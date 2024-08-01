package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

func (r *PostgresRepo) TxBegin() (*sqlx.Tx, error) {
	tx, err := database.TxBegin(r.DB)
	if err != nil {
		return nil, errors.Wrap(err, "PostgresRepo.TxBegin.TxBegin")
	}
	return tx, nil
}

func (r *PostgresRepo) TxEnd(tx *sqlx.Tx, err error) error {
	if err != nil {
		if rbErr := database.TxRollback(tx); rbErr != nil {
			return errors.Wrap(err, "PostgresRepo.TxEnd.TxRollback")
		}
		return apperror.ErrTxDone
	}

	err = database.TxCommit(tx)
	if err != nil {
		return errors.Wrap(err, "PostgresRepo.TxEnd.TxCommit")
	}

	return nil
}
