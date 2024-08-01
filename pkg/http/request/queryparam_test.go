package request

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/guregu/null/v5"
)

type TestStruct struct {
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	JoinedAt time.Time `json:"joined_at"`
	TestProfile
}

type TestProfile struct {
	URL             string      `json:"url"`
	email           string      // unexported test
	Fullname        string      // missing json tag
	UnsupportedNull null.String `json:"unsupported_null"`
}

func TestPopulateStructFromQueryParams(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expected    TestStruct
		expectError bool
	}{
		{
			name: "success decode queryparam",
			queryParams: map[string]string{
				"name":      "Alice",
				"age":       "30",
				"joined_at": "2023-08-01T12:00:00Z",
				"url":       "www.example.com",
			},
			expected: TestStruct{
				Name:     "Alice",
				Age:      30,
				JoinedAt: time.Date(2023, 8, 1, 12, 0, 0, 0, time.UTC),
				TestProfile: TestProfile{
					URL: "www.example.com",
				},
			},
			expectError: false,
		},
		{
			name: "success with missing query parameter",
			queryParams: map[string]string{
				"name": "Bob",
			},
			expected: TestStruct{
				Name:     "Bob",
				Age:      0,
				JoinedAt: time.Time{},
			},
			expectError: false,
		},
		{
			name: "failed due to invalid destination value",
			queryParams: map[string]string{
				"name": "Boby",
			},
			expected:    TestStruct{},
			expectError: true,
		},
		{
			name: "failed due to invalid integer value",
			queryParams: map[string]string{
				"name": "Charlie",
				"age":  "invalid",
			},
			expected: TestStruct{
				Name: "Charlie",
			},
			expectError: true,
		},
		{
			name: "failed due to invalid time value",
			queryParams: map[string]string{
				"name":      "Diana",
				"age":       "25",
				"joined_at": "invalid-time",
			},
			expected: TestStruct{
				Name: "Diana",
				Age:  25,
			},
			expectError: true,
		},
		{
			name: "failed due to unsupported type",
			queryParams: map[string]string{
				"name":             "Eve",
				"unsupported_null": "invalid-type",
			},
			expected: TestStruct{
				Name: "Eve",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test?"+encodeQueryParams(tt.queryParams), nil)

			var (
				result TestStruct
				err    error
			)

			if tt.name == "failed due to invalid destination value" {
				err = PopulateStructFromQueryParams(req, &tt.name)
			} else {
				err = PopulateStructFromQueryParams(req, &result)
			}

			if (err != nil) != tt.expectError {
				t.Errorf("PopulateStructFromQueryParams() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PopulateStructFromQueryParams() result = %v, want %v", result, tt.expected)
			}
		})
	}
}

func encodeQueryParams(params map[string]string) string {
	values := ""
	for key, value := range params {
		if values != "" {
			values += "&"
		}
		values += key + "=" + value
	}
	return values
}
