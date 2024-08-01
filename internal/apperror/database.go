package apperror

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("data must be unique")
	ErrTxDone    = sql.ErrTxDone
)
