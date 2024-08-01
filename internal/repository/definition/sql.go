package definition

import (
	"context"

	"github.com/jmoiron/sqlx"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
)

type SQLRepo interface {
	Transaction

	GetUser(ctx context.Context, filter req.UserFilter) ([]*model.User, error)
	CountUser(ctx context.Context, filter req.UserFilter) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	InsertUser(ctx context.Context, tx *sqlx.Tx, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) error
}

type Transaction interface {
	TxBegin() (*sqlx.Tx, error)
	TxEnd(tx *sqlx.Tx, err error) error
}
