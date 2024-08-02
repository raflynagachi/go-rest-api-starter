package validator

import (
	"testing"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTranslatorFunctions(t *testing.T) {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	tag := "testtag"
	msg := "This field is testtag"
	checkStr := func(fl validator.FieldLevel) bool {
		return fl.Field().String() != ""
	}

	if err := validate.RegisterValidation(tag, checkStr); err != nil {
		t.Fatalf("Failed to register validation: %v", err)
	}

	if err := validate.RegisterTranslation(
		tag, trans,
		registerTranslator(tag, msg),
		translate,
	); err != nil {
		t.Fatalf("Failed to register translation: %v", err)
	}

	type TestStruct struct {
		Name string `validate:"testtag"`
	}

	t.Run("registerTranslator", func(t *testing.T) {
		testStruct := TestStruct{}
		err := validate.Struct(testStruct).(validator.ValidationErrors)
		if assert.NotNil(t, err) {
			assert.Contains(t, err[0].Translate(trans), msg, "Expected error message to contain translation")
		}
	})

	t.Run("translate", func(t *testing.T) {
		testStruct := TestStruct{}
		valErrs := validate.Struct(testStruct)

		var fe validator.FieldError
		if valErrs, ok := valErrs.(validator.ValidationErrors); ok {
			for _, fieldErr := range valErrs {
				fe = fieldErr
				break
			}
		}

		result := translate(trans, fe)
		assert.Equal(t, msg, result, "Expected translated message did not match")
	})
}
