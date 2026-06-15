# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
make build        # compile binary to bin/go-rest-api-starter
make run          # run via go run
make test         # run all tests with -race and coverage
make coverage     # run tests + open HTML coverage report
make mock         # regenerate mocks (requires mockery)
make migrate-up   # run migrations (pass DB_HOST/PORT/NAME/USER/PASSWORD as env vars)
make migrate-down # revert last migration step
```

Run a single test:
```sh
go test -run TestFunctionName ./internal/handler/
```

## Environment Configuration

Config is loaded from `env/{appname}.{env}.json` where `appname = "go-rest-api-starter"` and `env` comes from the `env` environment variable (defaults to `development`). Copy `env/example.development.json` to `env/go-rest-api-starter.development.json` and fill in values before running.

## Architecture

This project follows clean architecture with three internal layers, each with a `definition/` package containing interfaces and a `mocks/` subdirectory for mockery-generated test doubles:

```
Handler → Usecase → Repository (Postgres)
```

**Wiring** (`cmd/http/main.go`): logger → config → db → `postgres.New(repo)` → `usecase.New(uc)` → `handler.New(hn)` → `router.New(r)`

**Layers:**
- `internal/handler/` — HTTP handlers using `httprouter`. Each method decodes the request, calls one usecase method, and writes a response via `pkg/http/response`.
- `internal/usecase/` — Business logic: validates requests with `pkg/validator`, orchestrates repository calls, wraps DB errors with HTTP status codes before returning.
- `internal/repository/postgres/` — Raw SQL queries using `sqlx`. Returns `apperror.ErrNotFound` / `apperror.ErrDuplicate` for domain errors; wraps mutation queries in transactions via `TxBegin`/`TxEnd`.

**Supporting packages:**
- `internal/dto/web/request|response` — HTTP-layer DTOs, distinct from `internal/model/` domain structs.
- `pkg/http/response` — `WrapErrBadRequest/NotFound/InternalServer` create typed `ErrResponse` values. `WriteFromError` unwraps the error chain to find the `ErrResponse`, logs it, and writes JSON. Internal server error messages are sanitized before sending to the client.
- `pkg/logger/` — Wraps `log/slog`. Uses a pretty handler in development and a context-aware JSON handler in production. Set to discard in tests (via `WithEnv("test")`).

**Error flow:** Repository returns `apperror.ErrNotFound` → usecase wraps it in `response.WrapErrNotFound(err)` → handler calls `response.WriteFromError(w, r, err, log)`.

## Testing Patterns

- Table-driven tests throughout.
- Package-level `var` function pointers (e.g., `var populateStructFromQueryParams = request.PopulateStructFromQueryParams`) allow test overrides without interface mocks.
- `internal/util/testutil/` provides `InitMockDB()` (sqlmock + sqlx) and a shared `MockErr` sentinel.
- `internal/util/random/` provides data generators like `randomutil.RandomUser()` for test fixtures.
- Handler tests use `httptest.NewRecorder()` and route through the real `httprouter` instance to exercise routing alongside handler logic.
- After adding new interface methods, regenerate mocks with `make mock`.
