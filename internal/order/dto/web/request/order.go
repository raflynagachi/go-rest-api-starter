package request

import sharedreq "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"

type OrderItemReq struct {
	SKU string `json:"sku" validate:"required"`
	Qty int64  `json:"qty" validate:"required,min=1"`
}

type CreateOrderReq struct {
	Items []OrderItemReq `json:"items" validate:"required,min=1,dive"`
}

type OrderFilter struct {
	UserEmail string `json:"user_email"`
	Status    string `json:"status"`
	sharedreq.Pagination
}
