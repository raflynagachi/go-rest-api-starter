package definition

import (
	"context"

	"github.com/jmoiron/sqlx"
	req "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/inventory/model"
)

type SQLRepo interface {
	Transaction

	GetItems(ctx context.Context, filter req.ItemFilter) ([]*model.Item, error)
	CountItems(ctx context.Context, filter req.ItemFilter) (int64, error)
	GetItemByID(ctx context.Context, id int64) (*model.Item, error)
	GetItemBySKU(ctx context.Context, sku string) (*model.Item, error)
	InsertItem(ctx context.Context, tx *sqlx.Tx, item *model.Item) (int64, error)
	UpdateStock(ctx context.Context, tx *sqlx.Tx, id, delta int64) error
}

type Transaction interface {
	TxBegin() (*sqlx.Tx, error)
	TxEnd(tx *sqlx.Tx, err error) error
}
