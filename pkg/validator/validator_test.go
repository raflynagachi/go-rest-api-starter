package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"-"`
	}

	t.Run("Valid struct", func(t *testing.T) {
		validStruct := TestStruct{Name: "Valid Name"}
		err := Validate(validStruct)
		assert.NoError(t, err, "Expected no validation errors")
	})

	t.Run("Invalid struct", func(t *testing.T) {
		invalidStruct := TestStruct{}
		err := Validate(invalidStruct)
		assert.Error(t, err, "Expected validation error")
		assert.Contains(t, err.Error(), "name", "Expected error to contain field name")
	})
}

func TestGetTranslator(t *testing.T) {
	t.Run("GetTranslator", func(t *testing.T) {
		translator := GetTranslator()
		assert.NotNil(t, translator, "Expected translator to be non-nil")

		tag := "required"
		field := "Name"
		expectedMsg := "Name is a required field"
		msg, err := translator.T(tag, field)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedMsg, msg, "Expected translated message did not match")
		}
	})
}
