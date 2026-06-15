package usecase

import (
	"time"

	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/order/repository/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/order/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/broker"
	"github.com/raflynagachi/go-rest-api-starter/pkg/client"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type APIUsecaseImpl struct {
	cfg              *config.Config
	appLogger        *logger.Logger
	repo             repo.SQLRepo
	broker           *broker.Broker
	inventoryClient  *client.Client // HTTP → Inventory service
	paymentClient    *client.Client // HTTP → Payment service
	// gRPC clients go here once proto is generated:
	// inventoryGRPC  inventorypb.InventoryServiceClient
	// paymentGRPC    paymentpb.PaymentServiceClient
}

func New(cfg *config.Config, log *logger.Logger, repo repo.SQLRepo, b *broker.Broker) uc.APIUsecase {
	return &APIUsecaseImpl{
		cfg:             cfg,
		appLogger:       log,
		repo:            repo,
		broker:          b,
		inventoryClient: client.New(cfg.Services.InventoryURL),
		paymentClient:   client.New(cfg.Services.PaymentURL),
	}
}

var getTimeNow = time.Now
