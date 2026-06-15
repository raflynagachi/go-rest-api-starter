package definition

import (
	"context"

	req "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/response"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type APIUsecase interface {
	GetOrders(ctx context.Context, filter req.OrderFilter) (*sharedresp.ListResponse, error)
	GetOrderByID(ctx context.Context, id int64) (*resp.OrderResponse, error)
	CreateOrder(ctx context.Context, r *req.CreateOrderReq) (*resp.OrderResponse, error)
}
