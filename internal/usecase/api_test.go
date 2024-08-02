package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/config"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	paginationutil "github.com/raflynagachi/go-rest-api-starter/internal/util/pagination"
	randomutil "github.com/raflynagachi/go-rest-api-starter/internal/util/random"
	"github.com/raflynagachi/go-rest-api-starter/internal/util/testutil"
)

func TestAPIUsecaseImpl_GetUser(t *testing.T) {
	mockUser := randomutil.RandomUser()
	mockUserResp := []*resp.UserResponse{{
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
	mockCount := int64(len(mockUserResp))
	mockPagination := req.Pagination{
		Page:  1,
		Limit: 1,
	}
	mockFilter := req.UserFilter{
		Pagination: mockPagination,
	}
	mockResp := &resp.ListResponse{
		Data: mockUserResp,
		PaginationResponse: resp.PaginationResponse{
			Page:      mockPagination.Page,
			TotalPage: paginationutil.TotalPage(mockCount, int64(mockPagination.Limit)),
			Limit:     mockPagination.Limit,
			Total:     mockCount,
		},
	}

	type fields struct {
		cfg  *config.Config
		repo repo.SQLRepo
	}
	type args struct {
		ctx    context.Context
		filter req.UserFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    *resp.ListResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx:    context.Background(),
				filter: mockFilter,
			},
			setup: func() {
				mockRepo.On("GetUser", context.Background(), mockFilter).
					Once().Return([]*model.User{mockUser}, nil)
				mockRepo.On("CountUser", context.Background(), mockFilter).
					Once().Return(mockCount, nil)
			},
			want:    mockResp,
			wantErr: false,
		},
		{
			name: "failed due to get user error",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx:    context.Background(),
				filter: mockFilter,
			},
			setup: func() {
				mockRepo.On("GetUser", context.Background(), mockFilter).
					Once().Return(nil, testutil.MockErr)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed due to count user error",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx:    context.Background(),
				filter: mockFilter,
			},
			setup: func() {
				mockRepo.On("GetUser", context.Background(), mockFilter).
					Once().Return([]*model.User{mockUser}, nil)
				mockRepo.On("CountUser", context.Background(), mockFilter).
					Once().Return(int64(0), testutil.MockErr)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &APIUsecaseImpl{
				cfg:  tt.fields.cfg,
				repo: tt.fields.repo,
			}

			tt.setup()

			got, err := u.GetUser(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIUsecaseImpl.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIUsecaseImpl.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIUsecaseImpl_GetUserByID(t *testing.T) {
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

	type fields struct {
		cfg  *config.Config
		repo repo.SQLRepo
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    *resp.UserResponse
		wantErr bool
	}{
		{
			name: "success get user by userID",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx: context.Background(),
				id:  mockUser.ID,
			},
			setup: func() {
				mockRepo.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(mockUser, nil)
			},
			want:    mockResp,
			wantErr: false,
		},
		{
			name: "failed due to user not found",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx: context.Background(),
				id:  mockUser.ID,
			},
			setup: func() {
				mockRepo.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(nil, apperror.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed due to connection error",
			fields: fields{
				cfg:  mockCfg,
				repo: mockRepo,
			},
			args: args{
				ctx: context.Background(),
				id:  mockUser.ID,
			},
			setup: func() {
				mockRepo.On("GetUserByID", context.Background(), mockUser.ID).
					Once().Return(nil, testutil.MockErr)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &APIUsecaseImpl{
				cfg:  tt.fields.cfg,
				repo: tt.fields.repo,
			}

			tt.setup()

			got, err := u.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIUsecaseImpl.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIUsecaseImpl.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
