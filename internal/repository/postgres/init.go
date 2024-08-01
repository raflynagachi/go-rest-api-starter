package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

type PostgresRepo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) repo.SQLRepo {
	return &PostgresRepo{
		DB: db,
	}
}

var (
	generatePagination = database.QueryPagination
)
