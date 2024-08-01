package testutil

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInitMockDB(t *testing.T) {
	t.Run("success mock DB", func(t *testing.T) {
		db, mockSql, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New() error = %v", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		gotDB, gotSQLXDB, gotMockSql, gotErr := InitMockDB()
		assert.NoError(t, gotErr)
		assert.IsType(t, db, gotDB)
		assert.IsType(t, sqlxDB, gotSQLXDB)
		assert.IsType(t, mockSql, gotMockSql)
	})

	t.Run("failure in sqlmockNew", func(t *testing.T) {
		tmpSqlmockNew := sqlmockNew
		defer func() {
			sqlmockNew = tmpSqlmockNew
		}()

		sqlmockNew = func() (*sql.DB, sqlmock.Sqlmock, error) {
			return nil, nil, MockErr
		}

		gotDB, gotSQLXDB, gotMockSql, gotErr := InitMockDB()
		assert.Error(t, gotErr)
		assert.Nil(t, gotDB)
		assert.Nil(t, gotSQLXDB)
		assert.Nil(t, gotMockSql)
	})
}

func TestInitBeginx(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	tests := []struct {
		name    string
		setup   func()
		wantTx  bool
		wantErr bool
	}{
		{
			name: "success begin transaction",
			setup: func() {
				mockSql.ExpectBegin()
			},
			wantTx:  true,
			wantErr: false,
		},
		{
			name: "failure in beginx",
			setup: func() {
				mockSql.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			wantTx:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := InitBeginx(sqlxDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitBeginx() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (got != nil) != tt.wantTx {
				t.Errorf("InitBeginx() tx = %v, wantTx %v", got, tt.wantTx)
			}
		})
	}
}
