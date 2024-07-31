package apperror

import "errors"

var (
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("data must be unique")
)
