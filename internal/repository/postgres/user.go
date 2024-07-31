package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
)

func (r *PostgresRepo) GetUser(ctx context.Context, filter req.UserFilter) ([]*model.User, error) {
	panic("need to be implemented")
}
func (r *PostgresRepo) CountUser(ctx context.Context, filter req.UserFilter) (int64, error) {
	panic("need to be implemented")
}
func (r *PostgresRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	panic("need to be implemented")
}
func (r *PostgresRepo) InsertUser(ctx context.Context, tx *sqlx.Tx, model *model.User) (int64, error) {
	panic("need to be implemented")
}
func (r *PostgresRepo) UpdateUser(ctx context.Context, tx *sqlx.Tx, model *model.User) error {
	panic("need to be implemented")
}
