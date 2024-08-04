
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.commitHash=`git rev-parse --short HEAD`"

# provide database information in environment
DB_DSN ?= postgres://$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASSWORD)
MIGRATE := docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.17.1 -path=/migrations/ -database "$(DB_DSN)"

build:
	@go build $(LDFLAGS) -o bin/go-rest-api-starter ./cmd/http

run:
	@go run ./cmd/http/

.PHONY: build run

############# TEST #############
mock:
	@mockery --all --dir internal/repository/definition --output=internal/repository/definition/mocks
	@mockery --all --dir internal/usecase/definition --output=internal/usecase/definition/mocks
	@mockery --all --dir internal/handler/definition --output=internal/handler/definition/mocks

test:
	@echo "Running tests..."
	@go test -p=1 -cover -race -covermode=atomic -coverprofile=coverage.out $(PACKAGES)
	@echo "Tests completed."

coverage: test
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage-all.out coverage.html
	@echo "Clean up complete."

PHONY: mock test coverage clean
############# TEST END #############

############# MIGRATIONS #############
migrate-up:
	@echo "Running all database migrations..."
	@$(MIGRATE) up

migrate-down:
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1
.PHONY: migrate-up migrate-down
############# MIGRATIONS END #############