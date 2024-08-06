package usecase

import (
	"time"

	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type APIUsecaseImpl struct {
	cfg       *config.Config
	appLogger *logger.Logger
	repo      repo.SQLRepo
}

func New(cfg *config.Config, log *logger.Logger, sqlRepo repo.SQLRepo) uc.APIUsecase {
	return &APIUsecaseImpl{
		cfg:       cfg,
		appLogger: log,
		repo:      sqlRepo,
	}
}

var (
	getTimeNow = time.Now()
)
