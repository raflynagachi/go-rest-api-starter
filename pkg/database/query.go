package database

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	sqlxIn = sqlx.In
)

// QueryPagination creates a SQL LIMIT and OFFSET clause for pagination.
// It validates that the page is >= 1 and the limit is > 0.
// Returns a string representing the SQL pagination clause.
func QueryPagination(page, limit int) (string, error) {
	if page < 1 {
		return "", fmt.Errorf("page number must be >= 1")
	}
	if limit <= 0 {
		return "", fmt.Errorf("limit must be > 0")
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	return query, nil
}

// BatchSelectContext executes the provided query in batches and collects the results.
func BatchSelectContext(ctx context.Context, db *sqlx.DB, query string, ids []int64, maxBatch int, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("destination must be a non-nil pointer")
	}

	destElem := destValue.Elem()
	if destElem.Kind() != reflect.Slice {
		return errors.New("destination must be a slice")
	}

	sliceType := destElem.Type().Elem()

	for start := 0; start < len(ids); start += maxBatch {
		end := start + maxBatch
		if end > len(ids) {
			end = len(ids)
		}

		idBatch := ids[start:end]
		q, args, err := sqlxIn(query, idBatch)
		if err != nil {
			return errors.Wrap(err, "sqlx.In")
		}

		q = db.Rebind(q)
		slicePtr := reflect.New(reflect.SliceOf(sliceType))
		err = db.SelectContext(ctx, slicePtr.Interface(), q, args...)
		if err != nil {
			return errors.Wrap(err, "SelectContext")
		}

		destElem.Set(reflect.AppendSlice(destElem, slicePtr.Elem()))
	}

	return nil
}
