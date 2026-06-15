package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/response"
	"github.com/raflynagachi/go-rest-api-starter/internal/inventory/model"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	paginationutil "github.com/raflynagachi/go-rest-api-starter/internal/util/pagination"
	"github.com/raflynagachi/go-rest-api-starter/pkg/auth"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
	"github.com/raflynagachi/go-rest-api-starter/pkg/validator"
	sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"
)

func (u *APIUsecaseImpl) GetItems(ctx context.Context, filter req.ItemFilter) (*sharedresp.ListResponse, error) {
	filter.Pagination.Validate()

	items, err := u.repo.GetItems(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetItems.GetItems")
	}

	count, err := u.repo.CountItems(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetItems.CountItems")
	}

	data := make([]*resp.ItemResponse, 0, len(items))
	for _, item := range items {
		data = append(data, toItemResponse(item))
	}

	return &sharedresp.ListResponse{
		Data: data,
		PaginationResponse: sharedresp.PaginationResponse{
			Page:      filter.Page,
			Limit:     filter.Limit,
			TotalPage: paginationutil.TotalPage(count, int64(filter.Limit)),
			Total:     count,
		},
	}, nil
}

func (u *APIUsecaseImpl) GetItemByID(ctx context.Context, id int64) (*resp.ItemResponse, error) {
	item, err := u.repo.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.GetItemByID")
		}
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetItemByID")
	}
	return toItemResponse(item), nil
}

func (u *APIUsecaseImpl) GetItemBySKU(ctx context.Context, sku string) (*resp.ItemResponse, error) {
	item, err := u.repo.GetItemBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.GetItemBySKU")
		}
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetItemBySKU")
	}
	return toItemResponse(item), nil
}

func (u *APIUsecaseImpl) CreateItem(ctx context.Context, r *req.CreateItemReq) error {
	if err := validator.Validate(r); err != nil {
		return errors.Wrap(response.WrapErrBadRequest(err), "APIUsecase.CreateItem.Validate")
	}

	item := &model.Item{
		SKU:      r.SKU,
		Name:     r.Name,
		Quantity: r.Quantity,
		Price:    r.Price,
		Created: sharedmodel.Created{
			CreatedAt: getTimeNow(),
			CreatedBy: auth.GetEmail(ctx),
		},
	}

	tx, err := u.repo.TxBegin()
	if err != nil {
		return errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateItem.TxBegin")
	}
	defer func() {
		if txErr := u.repo.TxEnd(tx, err); txErr != nil {
			u.appLogger.ErrorContext(ctx, errors.Wrap(txErr, "APIUsecase.CreateItem.TxEnd").Error())
		}
	}()

	_, err = u.repo.InsertItem(ctx, tx, item)
	if err != nil {
		if errors.Is(err, apperror.ErrDuplicate) {
			return errors.Wrap(response.WrapErrConflict(err), "APIUsecase.CreateItem.InsertItem")
		}
		return errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateItem.InsertItem")
	}
	return nil
}

func (u *APIUsecaseImpl) UpdateStock(ctx context.Context, id int64, r *req.UpdateStockReq) error {
	if err := validator.Validate(r); err != nil {
		return errors.Wrap(response.WrapErrBadRequest(err), "APIUsecase.UpdateStock.Validate")
	}

	_, err := u.repo.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.UpdateStock.GetItemByID")
		}
		return errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.UpdateStock.GetItemByID")
	}

	tx, err := u.repo.TxBegin()
	if err != nil {
		return errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.UpdateStock.TxBegin")
	}
	defer func() {
		if txErr := u.repo.TxEnd(tx, err); txErr != nil {
			u.appLogger.ErrorContext(ctx, errors.Wrap(txErr, "APIUsecase.UpdateStock.TxEnd").Error())
		}
	}()

	err = u.repo.UpdateStock(ctx, tx, id, r.Delta)
	if err != nil {
		return errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.UpdateStock.UpdateStock")
	}
	return nil
}

func toItemResponse(item *model.Item) *resp.ItemResponse {
	return &resp.ItemResponse{
		ID:       item.ID,
		SKU:      item.SKU,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
		CreatedResponse: sharedresp.CreatedResponse{
			CreatedAt: item.CreatedAt,
			CreatedBy: item.CreatedBy,
		},
		UpdatedResponse: sharedresp.UpdatedResponse{
			UpdatedAt: item.UpdatedAt,
			UpdatedBy: item.UpdatedBy,
		},
	}
}
