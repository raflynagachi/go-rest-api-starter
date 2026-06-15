package usecase

import (
	"time"

	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/payment/repository/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/payment/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/broker"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type APIUsecaseImpl struct {
	cfg       *config.Config
	appLogger *logger.Logger
	repo      repo.SQLRepo
	broker    *broker.Broker
}

func New(cfg *config.Config, log *logger.Logger, repo repo.SQLRepo, b *broker.Broker) uc.APIUsecase {
	return &APIUsecaseImpl{cfg: cfg, appLogger: log, repo: repo, broker: b}
}

var getTimeNow = time.Now
