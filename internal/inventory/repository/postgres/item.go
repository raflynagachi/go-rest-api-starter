package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/inventory/model"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

func filterItem(filter req.ItemFilter) (string, []interface{}) {
	conds := []string{"deleted_at IS NULL"}
	var args []interface{}

	if filter.SKU != "" {
		args = append(args, filter.SKU)
		conds = append(conds, fmt.Sprintf("sku = $%d", len(args)))
	}
	if filter.Name != "" {
		args = append(args, "%"+filter.Name+"%")
		conds = append(conds, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	return " WHERE " + strings.Join(conds, " AND "), args
}

func (r *PostgresRepo) GetItems(ctx context.Context, filter req.ItemFilter) ([]*model.Item, error) {
	query := `SELECT id, sku, name, quantity, price, created_at, created_by, updated_at, updated_by FROM items`
	where, args := filterItem(filter)
	pagination, err := generatePagination(filter.Page, filter.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "PostgresRepo.GetItems.generatePagination")
	}

	items := make([]*model.Item, 0)
	err = r.DB.SelectContext(ctx, &items, query+where+pagination, args...)
	return items, errors.Wrap(err, "PostgresRepo.GetItems.SelectContext")
}

func (r *PostgresRepo) CountItems(ctx context.Context, filter req.ItemFilter) (int64, error) {
	query := `SELECT COUNT(id) FROM items`
	where, args := filterItem(filter)

	var count int64
	err := r.DB.GetContext(ctx, &count, query+where, args...)
	return count, errors.Wrap(err, "PostgresRepo.CountItems.GetContext")
}

func (r *PostgresRepo) GetItemByID(ctx context.Context, id int64) (*model.Item, error) {
	query := r.DB.Rebind(`
		SELECT id, sku, name, quantity, price,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM items WHERE id = ? AND deleted_at IS NULL
	`)
	item := &model.Item{}
	err := r.DB.GetContext(ctx, item, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetItemByID")
	}
	return item, errors.Wrap(err, "PostgresRepo.GetItemByID")
}

func (r *PostgresRepo) GetItemBySKU(ctx context.Context, sku string) (*model.Item, error) {
	query := r.DB.Rebind(`
		SELECT id, sku, name, quantity, price,
		       created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM items WHERE sku = ? AND deleted_at IS NULL
	`)
	item := &model.Item{}
	err := r.DB.GetContext(ctx, item, query, sku)
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(apperror.ErrNotFound, "PostgresRepo.GetItemBySKU")
	}
	return item, errors.Wrap(err, "PostgresRepo.GetItemBySKU")
}

func (r *PostgresRepo) InsertItem(ctx context.Context, tx *sqlx.Tx, item *model.Item) (int64, error) {
	query := r.DB.Rebind(`
		INSERT INTO items (sku, name, quantity, price, created_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`)
	var id int64
	err := tx.GetContext(ctx, &id, query, item.SKU, item.Name, item.Quantity, item.Price, item.CreatedAt, item.CreatedBy)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == database.ERR_PQ_CODE_DUPLICATE {
			return 0, errors.Wrap(apperror.ErrDuplicate, "PostgresRepo.InsertItem")
		}
		return 0, errors.Wrap(err, "PostgresRepo.InsertItem")
	}
	return id, nil
}

func (r *PostgresRepo) UpdateStock(ctx context.Context, tx *sqlx.Tx, id, delta int64) error {
	query := r.DB.Rebind(`UPDATE items SET quantity = quantity + ? WHERE id = ? AND deleted_at IS NULL`)
	result, err := tx.ExecContext(ctx, query, delta, id)
	if err != nil {
		return errors.Wrap(err, "PostgresRepo.UpdateStock")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.Wrap(apperror.ErrNotFound, "PostgresRepo.UpdateStock")
	}
	return nil
}
