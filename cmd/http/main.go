package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	defer db.Close()

	repo := postgres.New(db)
	usecase := uc.New(cfg, repo)
	handler := hn.New(usecase)

	r := router.New(cfg, handler)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- r.ServeHTTP()
	}()

	// handle shutdown signals
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

	// the context is used as timeout to finish currently handling request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Shutdown(ctx); err != nil {
		log.Fatalf("error %v while shutting down Server\nInitiating force shutdown...", err)
	} else {
		log.Print("server exiting")
	}
}
