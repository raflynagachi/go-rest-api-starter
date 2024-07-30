mock:
	mockery --all --dir internal/repository/definition --output=internal/repository/definition/mocks
	mockery --all --dir internal/usecase/definition --output=internal/usecase/definition/mocks

test:
	go test ./... --cover --race

build:
	go build -o bin/go-rest-api-starter ./cmd/http

run:
	go run ./cmd/http/

.PHONY: mock test build run