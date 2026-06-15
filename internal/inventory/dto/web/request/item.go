package request

import sharedreq "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"

type CreateItemReq struct {
	SKU      string  `json:"sku"      validate:"required"`
	Name     string  `json:"name"     validate:"required"`
	Quantity int64   `json:"quantity" validate:"required,min=0"`
	Price    float64 `json:"price"    validate:"required,gt=0"`
}

type UpdateStockReq struct {
	Delta int64 `json:"delta" validate:"required"`
}

type ItemFilter struct {
	SKU  string `json:"sku"`
	Name string `json:"name"`
	sharedreq.Pagination
}
