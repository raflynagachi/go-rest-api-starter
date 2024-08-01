package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// gtToday check if a date is greater than today
func gtToday(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)

	return date.After(today)
}
