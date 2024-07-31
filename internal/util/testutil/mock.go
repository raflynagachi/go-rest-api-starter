package testutil

import (
	"database/sql"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	sqlmockNew = func() (*sql.DB, sqlmock.Sqlmock, error) {
		return sqlmock.New()
	}
	beginx = func(sqlxDB *sqlx.DB) (*sqlx.Tx, error) {
		return sqlxDB.Beginx()
	}
)

// InitMockDB create mock DB
func InitMockDB() (*sql.DB, *sqlx.DB, sqlmock.Sqlmock, error) {
	db, mockSql, err := sqlmockNew()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "sqlmock.New")
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return db, sqlxDB, mockSql, nil
}

// InitBeginx create mock transaction
func InitBeginx(t *testing.T, sqlxDB *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := beginx(sqlxDB)
	if err != nil {
		return nil, errors.Wrap(err, "Beginx")
	}
	return tx, nil
}
