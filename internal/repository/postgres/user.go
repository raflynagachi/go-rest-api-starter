package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

func (r *PostgresRepo) GetUser(ctx context.Context, filter req.UserFilter) ([]*model.User, error) {
	panic("need to be implemented")
}

func (r *PostgresRepo) CountUser(ctx context.Context, filter req.UserFilter) (int64, error) {
	panic("need to be implemented")
}

func (r *PostgresRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT
			id, email, created_at, created_by,
			updated_at, updated_by, deleted_at, deleted_by
		FROM users
		WHERE id = ?
	`

	query = r.DB.Rebind(query)

	user := &model.User{}
	err := r.DB.GetContext(ctx, user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetUserByID.GetContext")
		}
		return nil, errors.Wrap(err, "PostgresRepo.GetUserByID.GetContext")
	}

	return user, nil
}

func (r *PostgresRepo) InsertUser(ctx context.Context, tx *sqlx.Tx, user *model.User) (int64, error) {
	query := `
		INSERT INTO users (email, created_at, created_by)
		VALUES (:email, :created_at, :created_by)
	`

	result, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == database.ERR_PQ_CODE_DUPLICATE {
				return 0, errors.Wrap(apperror.ErrDuplicate, "PostgresRepo.InsertUser.NamedExecContext")
			}
		}
		return 0, errors.Wrap(err, "PostgresRepo.InsertUser.NamedExecContext")
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "PostgresRepo.InsertUser.LastInserId")
	}

	return lastID, nil
}

func (r *PostgresRepo) UpdateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := `
		UPDATE users SET
			email= COALESCE(:email, email),
			updated_at = COALESCE(:updated_at, updated_at),
			updated_by = COALESCE(:updated_by, updated_by)
		WHERE id = :id
	`

	_, err := r.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == database.ERR_PQ_CODE_DUPLICATE {
				return errors.Wrap(apperror.ErrDuplicate, "PostgresRepo.UpdateUser.NamedExecContext")
			}
		}
		return errors.Wrap(err, "PostgresRepo.UpdateUser.NamedExecContext")
	}

	return nil
}
