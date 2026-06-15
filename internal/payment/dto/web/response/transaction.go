package response

import (
	"github.com/raflynagachi/go-rest-api-starter/internal/payment/model"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
)

type TransactionResponse struct {
	ID      int64                    `json:"id"`
	OrderID int64                    `json:"order_id"`
	Amount  float64                  `json:"amount"`
	Status  model.TransactionStatus  `json:"status"`
	Method  string                   `json:"method"`
	sharedresp.CreatedResponse
	sharedresp.UpdatedResponse
}
