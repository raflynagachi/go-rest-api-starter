package usecase

import (
	"context"

	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

func (u *APIUsecaseImpl) GetUser(ctx context.Context, filter req.UserFilter) (*resp.ListResponse, error) {
	panic("need to be implemented")
}

func (u *APIUsecaseImpl) GetUserByID(ctx context.Context, id int64) (*resp.UserResponse, error) {
	panic("need to be implemented")
}

func (u *APIUsecaseImpl) CreateUser(ctx context.Context, request *req.CreateUpdateUserReq) error {
	panic("need to be implemented")
}

func (u *APIUsecaseImpl) UpdateUser(ctx context.Context, id int64, request *req.CreateUpdateUserReq) error {
	panic("need to be implemented")
}
