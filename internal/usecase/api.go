package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	paginationutil "github.com/raflynagachi/go-rest-api-starter/internal/util/pagination"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func (u *APIUsecaseImpl) GetUser(ctx context.Context, filter req.UserFilter) (*resp.ListResponse, error) {
	filter.Pagination.Validate()

	invs, err := u.repo.GetUser(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetUser.GetUser")
	}

	count, err := u.repo.CountUser(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetUser.CountUser")
	}

	userResp := make([]*resp.UserResponse, 0)
	for _, user := range invs {
		userResp = append(userResp, &resp.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			CreatedResponse: resp.CreatedResponse{
				CreatedAt: user.CreatedAt,
				CreatedBy: user.CreatedBy,
			},
			UpdatedResponse: resp.UpdatedResponse{
				UpdatedAt: user.UpdatedAt,
				UpdatedBy: user.UpdatedBy,
			},
		})
	}

	res := &resp.ListResponse{
		Data: userResp,
		PaginationResponse: resp.PaginationResponse{
			Page:      filter.Page,
			Limit:     filter.Limit,
			TotalPage: paginationutil.TotalPage(int64(count), int64(filter.Limit)),
			Total:     count,
		},
	}

	return res, nil
}

func (u *APIUsecaseImpl) GetUserByID(ctx context.Context, id int64) (*resp.UserResponse, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.GetUserByID.GetUserByID")
		}
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetUserByID.GetUserByID")
	}

	res := &resp.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		CreatedResponse: resp.CreatedResponse{
			CreatedAt: user.CreatedAt,
			CreatedBy: user.CreatedBy,
		},
		UpdatedResponse: resp.UpdatedResponse{
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
		},
	}

	return res, nil
}

func (u *APIUsecaseImpl) CreateUser(ctx context.Context, request *req.CreateUpdateUserReq) error {
	panic("need to be implemented")
}

func (u *APIUsecaseImpl) UpdateUser(ctx context.Context, id int64, request *req.CreateUpdateUserReq) error {
	panic("need to be implemented")
}
