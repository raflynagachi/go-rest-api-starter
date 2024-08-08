# Go REST API starter kit

[![Go Report Card](https://goreportcard.com/badge/github.com/raflynagachi/go-rest-api-starter)](https://goreportcard.com/report/github.com/raflynagachi/go-rest-api-starter)

## Overview
The Go REST API Starter Kit is a boilerplate project designed to help you quickly set up and develop RESTful APIs using Go.

## Prerequisites
- Go (version 1.22.x)
- Postgres 16
- mockery (for mock testing)

## Setup
1. **Create Database**
    ```sql
    CREATE DATABASE go-rest-api-starter;
    ```
2. **Run DDL Script**
    - Execute the DDL script in `migration` directory to set up the database schema.
    ```sh
    make migrate-up -e DB_HOST=localhost -e DB_PORT=5432 -e DB_NAME=go-rest-api-starter -e DB_USER=postgres -e DB_PASSWORD=postgres
    ```

3. **Configure Environment**
    - Create environment configuration files in the `env` directory with the filename format `{appname}.{env}.json`. Please check the example file.

## How to Run
Make sure to follow **Setup** section first
1. **Build the Project**
    ```sh
    make build # linux/mac
    ```
    Then, run the binary file `bin/go-rest-api-starter`

2. **Run Locally**
    ```sh
    make run
    ```

## How to Test
Run the unit test using:
```sh
make test
```

## How to mock
Run the mock (mockery required) using:
```sh
make mock
```
