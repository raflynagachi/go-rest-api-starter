package usecase

import (
	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
)

type APIUsecaseImpl struct {
	cfg  *config.Config
	repo repo.SQLRepo
}

func New(cfg *config.Config, sqlRepo repo.SQLRepo) uc.APIUsecase {
	return &APIUsecaseImpl{
		cfg:  cfg,
		repo: sqlRepo,
	}
}
