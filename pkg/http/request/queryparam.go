package request

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/guregu/null/v5"
)

// PopulateStructFromQueryParams stores query params to struct
func PopulateStructFromQueryParams(r *http.Request, dst interface{}) error {
	return populateStructFromQueryParams(r, reflect.ValueOf(dst).Elem())
}

// populateStructFromQueryParams is helper function to recursively populate struct fields
func populateStructFromQueryParams(r *http.Request, dstValue reflect.Value) error {
	if dstValue.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct but got %s", dstValue.Kind())
	}

	values := r.URL.Query()
	dstType := dstValue.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldValue := dstValue.Field(i)

		// skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// handle embedded structs except for specific type like time.Time
		if fieldValue.Kind() == reflect.Struct &&
			(field.Type != reflect.TypeOf(time.Time{}) ||
				field.Type != reflect.TypeOf(null.String{})) {
			if err := populateStructFromQueryParams(r, fieldValue); err != nil {
				return err
			}
		}

		// handle fields with tags
		queryTag := field.Tag.Get("json")
		if queryTag == "" {
			queryTag = field.Name
		}

		queryValue := values.Get(queryTag)
		if queryValue == "" {
			continue
		}

		// set the field value based on its type
		switch fieldValue.Interface().(type) {
		case string:
			fieldValue.SetString(queryValue)
		case int, int8, int16, int32, int64:
			intValue, err := strconv.ParseInt(queryValue, 10, fieldValue.Type().Bits())
			if err != nil {
				return fmt.Errorf("invalid integer value for field %s: %v", field.Name, err)
			}
			fieldValue.SetInt(intValue)
		case time.Time:
			timeValue, err := time.Parse(time.RFC3339, queryValue)
			if err != nil {
				return fmt.Errorf("invalid time value for field %s: %v", field.Name, err)
			}
			fieldValue.Set(reflect.ValueOf(timeValue))
		default:
			return fmt.Errorf("unsupported kind %s for field %s", fieldValue.Kind(), field.Name)
		}
	}

	return nil
}
