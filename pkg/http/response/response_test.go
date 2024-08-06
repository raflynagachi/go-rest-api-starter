package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name         string
		response     Response
		setup        func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "success write response",
			response: Response{
				Code: http.StatusOK,
				Data: map[string]string{"key": "value"},
			},
			setup:        func() {},
			expectedCode: http.StatusOK,
			expectedBody: `{"code":200,"data":{"key":"value"}}`,
		},
		{
			name: "failed due to encoding error",
			response: Response{
				Code: http.StatusInternalServerError,
				Data: map[string]string{"invalid": "invalid"},
			},
			setup: func() {
				encodeJson = func(w http.ResponseWriter, data interface{}) error {
					return errors.New("some error")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpEncodeJson := encodeJson
			defer func() {
				encodeJson = tmpEncodeJson
			}()

			recorder := httptest.NewRecorder()

			tt.setup()

			WriteResponse(recorder, tt.response, mockLogger)

			body := recorder.Body.String()
			body = strings.ReplaceAll(body, "\n", "")

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestWriteOKResponse(t *testing.T) {
	tests := []struct {
		name         string
		data         interface{}
		expectedBody string
	}{
		{
			name:         "successful response with data",
			data:         map[string]string{"key": "value"},
			expectedBody: `{"code":200,"data":{"key":"value"}}`,
		},
		{
			name:         "successful response with nil data",
			data:         nil,
			expectedBody: `{"code":200,"data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)

			WriteOKResponse(recorder, request, tt.data, mockLogger)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
