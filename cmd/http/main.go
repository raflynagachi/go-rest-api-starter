package main

import (
	"fmt"
	"log"

	"github.com/raflynagachi/go-rest-api-starter/config"
	"github.com/raflynagachi/go-rest-api-starter/internal/repository/postgres"
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

	// TODO: start usecase, handler and router
	fmt.Println(repo)
}
