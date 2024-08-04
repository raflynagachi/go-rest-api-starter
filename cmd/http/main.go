package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/router"
	"github.com/raflynagachi/go-rest-api-starter/internal/repository/postgres"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg.Databases[config.ServiceName])
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := postgres.New(db)
	usecase := uc.New(cfg, repo)
	handler := hn.New(usecase)

	r := router.New(cfg, handler)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- r.ServeHTTP()
	}()

	// Handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	case sig := <-stop:
		log.Printf("received signal: %v", sig)
	}

	// Graceful shutdown
	if shutdownErr := r.Shutdown(context.Background()); shutdownErr != nil {
		log.Printf("error during shutdown: %v", shutdownErr)
	}

}
