package definition

import (
	"context"

	req "github.com/raflynagachi/go-rest-api-starter/internal/payment/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/payment/dto/web/response"
)

type APIUsecase interface {
	GetTransactionByID(ctx context.Context, id int64) (*resp.TransactionResponse, error)
	CreateTransaction(ctx context.Context, r *req.CreateTransactionReq) (*resp.TransactionResponse, error)
}
