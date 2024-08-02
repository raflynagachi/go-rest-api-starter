package validator

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestGtTodayValidator(t *testing.T) {
	type TestStruct struct {
		InvalidField string    `validate:"omitempty,gt_today"`
		Date         time.Time `validate:"gt_today"`
	}

	validate := validator.New()
	validate.RegisterValidation("gt_today", gtToday)

	tests := []struct {
		name     string
		input    TestStruct
		expected bool
	}{
		{
			name: "success with future date",
			input: TestStruct{
				Date: time.Now().Add(48 * time.Hour), // Future time
			},
			expected: true,
		},
		{
			name: "failed due to invalid today date",
			input: TestStruct{
				Date: time.Now().Truncate(24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "failed due to invalid past date",
			input: TestStruct{
				Date: time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour),
			},
			expected: false,
		},
		{
			name: "failed due to empty date",
			input: TestStruct{
				Date: time.Time{},
			},
			expected: false,
		},
		{
			name: "failed due to invalid type",
			input: TestStruct{
				InvalidField: "invalid",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
