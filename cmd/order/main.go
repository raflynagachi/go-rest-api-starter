package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/order/handler"
	"github.com/raflynagachi/go-rest-api-starter/internal/order/handler/router"
	"github.com/raflynagachi/go-rest-api-starter/internal/order/repository/postgres"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/order/usecase"
	"github.com/raflynagachi/go-rest-api-starter/pkg/broker"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

func main() {
	config.ServiceName = "order"

	appLogger := logger.NewLogger(logger.WithEnv(config.Env))

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("failed to load config", logger.ErrAttr(err))
		return
	}

	db, err := database.ConnectDB(cfg.Databases[config.ServiceName])
	if err != nil {
		appLogger.Error("failed to connect database", logger.ErrAttr(err))
		return
	}
	defer db.Close()

	b, err := broker.New(cfg.Services.NatsURL)
	if err != nil {
		appLogger.Error("failed to connect broker", logger.ErrAttr(err))
		return
	}
	defer b.Close()

	repo := postgres.New(db, appLogger)
	usecase := uc.New(cfg, appLogger, repo, b)
	handler := hn.New(usecase, appLogger)
	r := router.New(cfg, appLogger, handler, db)

	serverErr := make(chan error, 1)
	go func() { serverErr <- r.ServeHTTP() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil {
			appLogger.Error("server error", logger.ErrAttr(err))
		}
	case sig := <-stop:
		appLogger.Info("received signal", logger.StringAttr("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Shutdown(ctx); err != nil {
		appLogger.Error("shutdown error", logger.ErrAttr(err))
	} else {
		appLogger.Info("server exiting")
	}
}
