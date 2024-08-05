package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/config"
)

var (
	sqlxConnect = sqlx.Connect
)

func ConnectDB(cfg *config.Database) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sqlxConnect("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Connect")
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	return db, nil
}
