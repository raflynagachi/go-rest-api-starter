package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/payment/repository/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type PostgresRepo struct {
	DB        *sqlx.DB
	appLogger *logger.Logger
}

func New(db *sqlx.DB, log *logger.Logger) repo.SQLRepo {
	return &PostgresRepo{DB: db, appLogger: log}
}
