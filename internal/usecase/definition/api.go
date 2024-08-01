package definition

import (
	"context"

	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type APIUsecase interface {
	GetUser(ctx context.Context, filter req.UserFilter) (*resp.ListResponse, error)
	GetUserByID(ctx context.Context, id int64) (*resp.UserResponse, error)
	CreateUser(ctx context.Context, request *req.CreateUpdateUserReq) error
	UpdateUser(ctx context.Context, id int64, request *req.CreateUpdateUserReq) error
}
