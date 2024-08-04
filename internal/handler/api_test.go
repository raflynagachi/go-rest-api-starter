package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	randomutil "github.com/raflynagachi/go-rest-api-starter/internal/util/random"
	"github.com/raflynagachi/go-rest-api-starter/internal/util/testutil"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func TestAPIHandlerImpl_GetUser(t *testing.T) {
	mockUser := randomutil.RandomUser()
	mockInvResp := []*resp.UserResponse{{
		ID:    mockUser.ID,
		Email: mockUser.Email,
		CreatedResponse: resp.CreatedResponse{
			CreatedAt: mockUser.CreatedAt,
			CreatedBy: mockUser.CreatedBy,
		},
		UpdatedResponse: resp.UpdatedResponse{
			UpdatedAt: mockUser.UpdatedAt,
			UpdatedBy: mockUser.UpdatedBy,
		},
	}}
	mockResp := &resp.ListResponse{
		Data:               mockInvResp,
		PaginationResponse: resp.PaginationResponse{},
	}

	tmpPopulateStructFromQueryParams := populateStructFromQueryParams
	defer func() {
		populateStructFromQueryParams = tmpPopulateStructFromQueryParams
	}()

	type args struct {
		request *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "success get user",
			args: func(t *testing.T) args {
				populateStructFromQueryParams = func(r *http.Request, dst interface{}) error {
					return nil
				}

				request, err := http.NewRequest(http.MethodGet, "/users", http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("GetUser", context.Background(), req.UserFilter{}).
					Once().Return(mockResp, nil)

				return args{request: request}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed due to invalid query param",
			args: func(t *testing.T) args {
				populateStructFromQueryParams = func(r *http.Request, dst interface{}) error {
					return testutil.MockErr
				}

				req, err := http.NewRequest(http.MethodGet, "/users", http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed due to internal server error",
			args: func(t *testing.T) args {
				populateStructFromQueryParams = func(r *http.Request, dst interface{}) error {
					return nil
				}

				request, err := http.NewRequest(http.MethodGet, "/users", http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("GetUser", context.Background(), req.UserFilter{}).
					Once().Return(nil, response.WrapErrInternalServer(testutil.MockErr))

				return args{request: request}
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			mockHandler.Router.ServeHTTP(resp, tArgs.request)
			res := resp.Result()

			if res.StatusCode != tt.wantCode {
				t.Errorf("APIHandler.GetUser() code = %v, wantCode %v", res.StatusCode, tt.wantCode)
				return
			}
		})
	}
}

func TestAPIHandlerImpl_GetUserByID(t *testing.T) {
	mockUser := randomutil.RandomUser()
	mockResp := &resp.UserResponse{
		ID:    mockUser.ID,
		Email: mockUser.Email,
		CreatedResponse: resp.CreatedResponse{
			CreatedAt: mockUser.CreatedAt,
			CreatedBy: mockUser.CreatedBy,
		},
		UpdatedResponse: resp.UpdatedResponse{
			UpdatedAt: mockUser.UpdatedAt,
			UpdatedBy: mockUser.UpdatedBy,
		},
	}

	type args struct {
		request *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "success get user by ID",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(mockResp, nil)

				return args{request: req}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed due to invalid query param",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%s", "invalidID")
				req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}
				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed due to user not found",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(nil, response.WrapErrNotFound(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "failed due to internal server error",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(nil, response.WrapErrInternalServer(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			mockHandler.Router.ServeHTTP(resp, tArgs.request)
			res := resp.Result()

			if res.StatusCode != tt.wantCode {
				t.Errorf("APIHandler.GetUserByID() code = %v, wantCode %v", res.StatusCode, tt.wantCode)
				return
			}
		})
	}
}

func TestAPIHandlerImpl_CreateUser(t *testing.T) {
	mockUser := randomutil.RandomUser()

	mockReq := &req.CreateUpdateUserReq{
		Email: mockUser.Email,
	}

	type args struct {
		request *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "success create user",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("CreateUser", req.Context(), reqModel).
					Once().Return(nil)

				return args{request: req}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed due to json encode error",
			args: func(t *testing.T) args {
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid")))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed due to bad request",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("CreateUser", req.Context(), reqModel).
					Once().Return(response.WrapErrBadRequest(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed_errorInternalServer",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("CreateUser", req.Context(), reqModel).
					Once().Return(response.WrapErrInternalServer(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			mockHandler.Router.ServeHTTP(resp, tArgs.request)
			res := resp.Result()

			if res.StatusCode != tt.wantCode {
				t.Errorf("APIHandler.CreateUser() code = %v, wantCode %v", res.StatusCode, tt.wantCode)
				return
			}
		})
	}
}

func TestAPIHandlerImpl_UpdateUser(t *testing.T) {
	mockUser := randomutil.RandomUser()

	mockReq := &req.CreateUpdateUserReq{
		Email: mockUser.Email,
	}

	type args struct {
		request *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "success_updateUser",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("UpdateUser", req.Context(), mockUser.ID, reqModel).
					Once().Return(nil)

				return args{request: req}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "failed_errorInvalidIDParam",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%s", "invalid")
				req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer([]byte("invalid")))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed_errorJsonDecode",
			args: func(t *testing.T) args {
				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer([]byte("invalid")))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed_errorBadRequest",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("UpdateUser", req.Context(), mockUser.ID, reqModel).
					Once().Return(response.WrapErrBadRequest(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "failed_errorInternalServer",
			args: func(t *testing.T) args {
				reqModel := mockReq

				reqBody, err := json.Marshal(reqModel)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				path := fmt.Sprintf("/users/%d", mockUser.ID)
				req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("fail to create request: %v", err)
				}

				mockUc.On("UpdateUser", req.Context(), mockUser.ID, reqModel).
					Once().Return(response.WrapErrInternalServer(testutil.MockErr))

				return args{request: req}
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			mockHandler.Router.ServeHTTP(resp, tArgs.request)
			res := resp.Result()

			if res.StatusCode != tt.wantCode {
				t.Errorf("APIHandler.UpdateUser() code = %v, wantCode %v", res.StatusCode, tt.wantCode)
				return
			}
		})
	}
}
