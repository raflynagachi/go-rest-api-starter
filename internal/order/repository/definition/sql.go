package definition

import (
	"context"

	"github.com/jmoiron/sqlx"
	req "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/order/model"
)

type SQLRepo interface {
	Transaction

	GetOrders(ctx context.Context, filter req.OrderFilter) ([]*model.Order, error)
	CountOrders(ctx context.Context, filter req.OrderFilter) (int64, error)
	GetOrderByID(ctx context.Context, id int64) (*model.Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]*model.OrderItem, error)
	InsertOrder(ctx context.Context, tx *sqlx.Tx, order *model.Order) (int64, error)
	InsertOrderItem(ctx context.Context, tx *sqlx.Tx, item *model.OrderItem) error
}

type Transaction interface {
	TxBegin() (*sqlx.Tx, error)
	TxEnd(tx *sqlx.Tx, err error) error
}
