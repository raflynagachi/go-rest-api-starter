package definition

import (
	"context"

	req "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/response"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type APIUsecase interface {
	GetItems(ctx context.Context, filter req.ItemFilter) (*sharedresp.ListResponse, error)
	GetItemByID(ctx context.Context, id int64) (*resp.ItemResponse, error)
	GetItemBySKU(ctx context.Context, sku string) (*resp.ItemResponse, error)
	CreateItem(ctx context.Context, r *req.CreateItemReq) error
	UpdateStock(ctx context.Context, id int64, r *req.UpdateStockReq) error
}
