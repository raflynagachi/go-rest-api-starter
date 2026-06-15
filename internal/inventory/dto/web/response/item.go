package response

import (
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type ItemResponse struct {
	ID       int64   `json:"id"`
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
	sharedresp.CreatedResponse
	sharedresp.UpdatedResponse
}
