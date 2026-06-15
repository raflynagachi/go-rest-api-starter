package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/order/model"
)

func filterOrder(filter req.OrderFilter) (string, []interface{}) {
	conds := []string{"deleted_at IS NULL"}
	var args []interface{}

	if filter.UserEmail != "" {
		args = append(args, filter.UserEmail)
		conds = append(conds, fmt.Sprintf("user_email = $%d", len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		conds = append(conds, fmt.Sprintf("status = $%d", len(args)))
	}

	return " WHERE " + strings.Join(conds, " AND "), args
}

func (r *PostgresRepo) GetOrders(ctx context.Context, filter req.OrderFilter) ([]*model.Order, error) {
	query := `SELECT id, user_email, status, total_amount, created_at, created_by, updated_at, updated_by FROM orders`
	where, args := filterOrder(filter)
	pagination, err := generatePagination(filter.Page, filter.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "PostgresRepo.GetOrders.generatePagination")
	}

	orders := make([]*model.Order, 0)
	err = r.DB.SelectContext(ctx, &orders, query+where+pagination, args...)
	return orders, errors.Wrap(err, "PostgresRepo.GetOrders")
}

func (r *PostgresRepo) CountOrders(ctx context.Context, filter req.OrderFilter) (int64, error) {
	query := `SELECT COUNT(id) FROM orders`
	where, args := filterOrder(filter)

	var count int64
	err := r.DB.GetContext(ctx, &count, query+where, args...)
	return count, errors.Wrap(err, "PostgresRepo.CountOrders")
}

func (r *PostgresRepo) GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	query := r.DB.Rebind(`
		SELECT id, user_email, status, total_amount,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM orders WHERE id = ? AND deleted_at IS NULL
	`)
	order := &model.Order{}
	if err := r.DB.GetContext(ctx, order, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetOrderByID")
		}
		return nil, errors.Wrap(err, "PostgresRepo.GetOrderByID")
	}
	return order, nil
}

func (r *PostgresRepo) GetOrderItems(ctx context.Context, orderID int64) ([]*model.OrderItem, error) {
	query := r.DB.Rebind(`SELECT id, order_id, sku, qty, price FROM order_items WHERE order_id = ?`)
	items := make([]*model.OrderItem, 0)
	err := r.DB.SelectContext(ctx, &items, query, orderID)
	return items, errors.Wrap(err, "PostgresRepo.GetOrderItems")
}

func (r *PostgresRepo) InsertOrder(ctx context.Context, tx *sqlx.Tx, order *model.Order) (int64, error) {
	query := r.DB.Rebind(`
		INSERT INTO orders (user_email, status, total_amount, created_at, created_by)
		VALUES (?, ?, ?, ?, ?) RETURNING id
	`)
	var id int64
	err := tx.GetContext(ctx, &id, query, order.UserEmail, order.Status, order.TotalAmount, order.CreatedAt, order.CreatedBy)
	return id, errors.Wrap(err, "PostgresRepo.InsertOrder")
}

func (r *PostgresRepo) InsertOrderItem(ctx context.Context, tx *sqlx.Tx, item *model.OrderItem) error {
	query := r.DB.Rebind(`INSERT INTO order_items (order_id, sku, qty, price) VALUES (?, ?, ?, ?)`)
	_, err := tx.ExecContext(ctx, query, item.OrderID, item.SKU, item.Qty, item.Price)
	return errors.Wrap(err, "PostgresRepo.InsertOrderItem")
}
