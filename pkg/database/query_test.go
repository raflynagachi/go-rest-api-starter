package database

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchSelectContext(t *testing.T) {
	type User struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}

	tests := []struct {
		name        string
		query       string
		ids         []int64
		maxBatch    int
		setup       func(mock sqlmock.Sqlmock)
		dest        interface{}
		expected    []User
		expectedErr bool
	}{
		{
			name:     "success with multiple batches",
			query:    "SELECT * FROM table WHERE id IN (?)",
			ids:      []int64{1, 2, 3, 4, 5},
			maxBatch: 2,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM table WHERE id IN \(\?, \?\)`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Alice").AddRow(2, "Bob"))
				mock.ExpectQuery(`SELECT \* FROM table WHERE id IN \(\?, \?\)`).
					WithArgs(3, 4).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(3, "Charlie").AddRow(4, "David"))
				mock.ExpectQuery(`SELECT \* FROM table WHERE id IN \(\?\)`).
					WithArgs(5).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(5, "Eve"))
			},
			dest:        &[]User{},
			expected:    []User{{1, "Alice"}, {2, "Bob"}, {3, "Charlie"}, {4, "David"}, {5, "Eve"}},
			expectedErr: false,
		},
		{
			name:     "error due to connection error",
			query:    "SELECT * FROM table WHERE id IN (?)",
			ids:      []int64{1},
			maxBatch: 1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM table WHERE id IN (?)").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			dest:        &[]User{},
			expected:    nil,
			expectedErr: true,
		},
		{
			name:     "error when sqlx.In error",
			query:    "SELECT * FROM table WHERE id IN (?)",
			ids:      []int64{1, 2, 3},
			maxBatch: 2,
			setup: func(mock sqlmock.Sqlmock) {
				sqlxIn = func(query string, args ...interface{}) (string, []interface{}, error) {
					return "", []interface{}{}, errors.New("some error")
				}
			},
			dest:        &[]User{},
			expectedErr: true,
		},
		{
			name:        "error due to destination is not a slice",
			query:       "SELECT * FROM table WHERE id IN (?)",
			ids:         []int64{1},
			maxBatch:    1,
			setup:       func(mock sqlmock.Sqlmock) {},
			dest:        new(int),
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "error due to destination is a nil pointer",
			query:       "SELECT * FROM table WHERE id IN (?)",
			ids:         []int64{1},
			maxBatch:    1,
			setup:       func(mock sqlmock.Sqlmock) {},
			dest:        nil,
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			originalSqlxIn := sqlxIn
			defer func() {
				sqlxIn = originalSqlxIn
			}()

			tt.setup(mock)

			err = BatchSelectContext(context.Background(), sqlxDB, tt.query, tt.ids, tt.maxBatch, tt.dest)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.dest != nil {
					result := reflect.ValueOf(tt.dest).Elem().Interface()
					assert.ElementsMatch(t, tt.expected, result)
				}
			}
		})
	}
}
