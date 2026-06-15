PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /mocks)
LDFLAGS  := -ldflags "-X main.commitHash=`git rev-parse --short HEAD`"

DB_DSN ?= postgres://$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&user=$(DB_USER)&password=$(DB_PASSWORD)

############# BUILD #############
build:
	@go build $(LDFLAGS) -o bin/go-rest-api-starter ./cmd/http

build-inventory:
	@go build $(LDFLAGS) -o bin/inventory ./cmd/inventory

build-payment:
	@go build $(LDFLAGS) -o bin/payment ./cmd/payment

build-order:
	@go build $(LDFLAGS) -o bin/order ./cmd/order

build-all: build build-inventory build-payment build-order

.PHONY: build build-inventory build-payment build-order build-all

############# RUN #############
run:
	@go run ./cmd/http/

run-inventory:
	@go run ./cmd/inventory/

run-payment:
	@go run ./cmd/payment/

run-order:
	@go run ./cmd/order/

.PHONY: run run-inventory run-payment run-order

############# TEST #############
mock:
	@mockery --all --dir internal/repository/definition         --output=internal/repository/definition/mocks
	@mockery --all --dir internal/usecase/definition            --output=internal/usecase/definition/mocks
	@mockery --all --dir internal/handler/definition            --output=internal/handler/definition/mocks
	@mockery --all --dir internal/inventory/repository/definition --output=internal/inventory/repository/definition/mocks
	@mockery --all --dir internal/inventory/usecase/definition    --output=internal/inventory/usecase/definition/mocks
	@mockery --all --dir internal/inventory/handler/definition    --output=internal/inventory/handler/definition/mocks
	@mockery --all --dir internal/payment/repository/definition  --output=internal/payment/repository/definition/mocks
	@mockery --all --dir internal/payment/usecase/definition     --output=internal/payment/usecase/definition/mocks
	@mockery --all --dir internal/payment/handler/definition     --output=internal/payment/handler/definition/mocks
	@mockery --all --dir internal/order/repository/definition    --output=internal/order/repository/definition/mocks
	@mockery --all --dir internal/order/usecase/definition       --output=internal/order/usecase/definition/mocks
	@mockery --all --dir internal/order/handler/definition       --output=internal/order/handler/definition/mocks

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

.PHONY: mock test coverage clean

############# PROTO (gRPC) #############
proto:
	@protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/inventory/inventory.proto \
		proto/payment/payment.proto

.PHONY: proto

############# MIGRATIONS #############
MIGRATE      := docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.17.1

migrate-up:
	@echo "Running all database migrations..."
	@$(MIGRATE) -path=/migrations/ -database "$(DB_DSN)" up

migrate-down:
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) -path=/migrations/ -database "$(DB_DSN)" down 1

migrate-up-inventory:
	@$(MIGRATE) -path=/migrations/inventory -database "$(DB_DSN)" up

migrate-up-payment:
	@$(MIGRATE) -path=/migrations/payment -database "$(DB_DSN)" up

migrate-up-order:
	@$(MIGRATE) -path=/migrations/order -database "$(DB_DSN)" up

.PHONY: migrate-up migrate-down migrate-up-inventory migrate-up-payment migrate-up-order
