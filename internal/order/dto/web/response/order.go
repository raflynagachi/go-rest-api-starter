package response

import (
	"github.com/raflynagachi/go-rest-api-starter/internal/order/model"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type OrderItemResponse struct {
	SKU   string  `json:"sku"`
	Qty   int64   `json:"qty"`
	Price float64 `json:"price"`
}

type OrderResponse struct {
	ID          int64                `json:"id"`
	UserEmail   string               `json:"user_email"`
	Status      model.OrderStatus    `json:"status"`
	TotalAmount float64              `json:"total_amount"`
	Items       []*OrderItemResponse `json:"items"`
	sharedresp.CreatedResponse
	sharedresp.UpdatedResponse
}
