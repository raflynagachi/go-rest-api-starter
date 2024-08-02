package validator

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate = validator.New()

	enLocale = en.New()
	uni      = ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")
)

func init() {
	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatalf("Error registering translations: %v", err)
	}

	// register custom validation tag
	if err := validate.RegisterValidation("gt_today", gtToday); err != nil {
		log.Fatalf("Error registering validation: %v", err)
	}

	if err := validate.RegisterTranslation(
		"gt_today", trans,
		registerTranslator("gt_today", "{0} must be a date after today"),
		translate,
	); err != nil {
		log.Fatalf("Error registering translation: %v", err)
	}

	// register json tag for error validations
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate(s any) error {
	return validate.Struct(s)
}

func GetTranslator() ut.Translator {
	return trans
}
