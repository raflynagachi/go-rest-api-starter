package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestWriteFromError(t *testing.T) {
	type TestStruct struct {
		Name string `validate:"required"`
		Age  int    `validate:"required"`
	}

	mockVal := validator.New()

	tests := []struct {
		name         string
		err          error
		setup        func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "success with validation error",
			err:  errors.New("validation error"),
			setup: func() {
				mockTestStruct := TestStruct{
					Name: "Alice",
				}
				valErrs := mockVal.Struct(mockTestStruct)

				findErrResponse = func(e error) (ErrResponse, error) {
					return ErrResponse{
						Code:    http.StatusBadRequest,
						Message: "",
						Err:     valErrs,
					}, valErrs
				}
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"code":400,"message":"Key: 'TestStruct.Age' Error:Field validation for 'Age' failed on the 'required' tag"}`,
		},
		{
			name: "success with multiple validation error",
			err:  errors.New("validation error"),
			setup: func() {
				mockTestStruct := TestStruct{}
				valErrs := mockVal.Struct(mockTestStruct)

				findErrResponse = func(e error) (ErrResponse, error) {
					return ErrResponse{
						Code:    http.StatusBadRequest,
						Message: "",
						Err:     valErrs,
					}, valErrs
				}
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"code":400,"message":"Key: 'TestStruct.Name' Error:Field validation for 'Name' failed on the 'required' tag,Key: 'TestStruct.Age' Error:Field validation for 'Age' failed on the 'required' tag"}`,
		},
		{
			name: "success with internal server error",
			err:  errors.New("internal server error"),
			setup: func() {
				findErrResponse = func(e error) (ErrResponse, error) {
					return ErrResponse{
						Code:    http.StatusInternalServerError,
						Message: "",
						Err:     errors.New("internal server error"),
					}, errors.New("internal server error")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"code":500,"message":"internal server error"}`,
		},
		{
			name: "encoding error",
			err:  errors.New("some error"),
			setup: func() {
				findErrResponse = func(e error) (ErrResponse, error) {
					return ErrResponse{
						Code:    http.StatusInternalServerError,
						Message: "",
						Err:     errors.New("some error"),
					}, errors.New("some error")
				}
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
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)

			tmpFindErrResponse := findErrResponse
			tmpEncodeJson := encodeJson

			tt.setup()

			defer func() {
				findErrResponse = tmpFindErrResponse
				encodeJson = tmpEncodeJson
			}()

			WriteFromError(recorder, request, tt.err, mockLogger)

			body := recorder.Body.String()
			body = strings.ReplaceAll(body, "\n", "")

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestFindErrResponse(t *testing.T) {
	mockWrappedErr := errors.New("wrapped error")
	mockInnerErr := errors.New("inner error")
	mockNonErrResponse := errors.New("some non-ErrResponse error")

	tests := []struct {
		name            string
		inputError      error
		expectedErrResp ErrResponse
		expectedLastErr error
	}{
		{
			name:            "direct ErrResponse",
			inputError:      ErrResponse{Code: 400, Message: "Bad Request", Err: nil},
			expectedErrResp: ErrResponse{Code: 400, Message: "Bad Request", Err: nil},
			expectedLastErr: nil,
		},
		{
			name: "wrapped ErrResponse",
			inputError: fmt.Errorf("some error: %w",
				ErrResponse{Code: 404, Message: "Not Found", Err: mockWrappedErr}),
			expectedErrResp: ErrResponse{Code: 404, Message: "Not Found", Err: mockWrappedErr},
			expectedLastErr: mockWrappedErr,
		},
		{
			name: "nested ErrResponse",
			inputError: fmt.Errorf("outer error: %w",
				fmt.Errorf("inner error: %w",
					ErrResponse{Code: 403, Message: "Forbidden", Err: mockInnerErr})),
			expectedErrResp: ErrResponse{Code: 403, Message: "Forbidden", Err: mockInnerErr},
			expectedLastErr: mockInnerErr,
		},
		{
			name:            "non-ErrResponse Error",
			inputError:      mockNonErrResponse,
			expectedErrResp: WrapErrInternalServer(nil),
			expectedLastErr: nil,
		},
		{
			name:            "no Error",
			inputError:      nil,
			expectedErrResp: WrapErrInternalServer(nil),
			expectedLastErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errResp, lastErr := FindErrResponse(tt.inputError)

			assert.Equal(t, tt.expectedErrResp, errResp)
			assert.Equal(t, tt.expectedLastErr, lastErr)
		})
	}
}

func TestWrapErrFunctions(t *testing.T) {
	mockBadReqErr := errors.New("bad request error")
	mockNotFoundErr := errors.New("not found error")
	mockInternalErr := errors.New("internal server error")

	tests := []struct {
		name         string
		wrapFunc     func(error) ErrResponse
		inputError   error
		expectedCode int
		expectedErr  error
	}{
		{
			name:         "WrapErrBadRequest",
			wrapFunc:     WrapErrBadRequest,
			inputError:   mockBadReqErr,
			expectedCode: http.StatusBadRequest,
			expectedErr:  mockBadReqErr,
		},
		{
			name:         "WrapErrNotFound",
			wrapFunc:     WrapErrNotFound,
			inputError:   mockNotFoundErr,
			expectedCode: http.StatusNotFound,
			expectedErr:  mockNotFoundErr,
		},
		{
			name:         "WrapErrInternalServer",
			wrapFunc:     WrapErrInternalServer,
			inputError:   mockInternalErr,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  mockInternalErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errResp := tt.wrapFunc(tt.inputError)

			assert.Equal(t, tt.expectedCode, errResp.Code)
			assert.Equal(t, tt.expectedErr, errResp.Err)
		})
	}
}
