package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/payment/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/payment/dto/web/response"
	"github.com/raflynagachi/go-rest-api-starter/internal/payment/model"
	sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	"github.com/raflynagachi/go-rest-api-starter/pkg/auth"
	"github.com/raflynagachi/go-rest-api-starter/pkg/broker"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
	"github.com/raflynagachi/go-rest-api-starter/pkg/validator"
)

func (u *APIUsecaseImpl) GetTransactionByID(ctx context.Context, id int64) (*resp.TransactionResponse, error) {
	t, err := u.repo.GetTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.GetTransactionByID")
		}
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetTransactionByID")
	}
	return toTransactionResponse(t), nil
}

func (u *APIUsecaseImpl) CreateTransaction(ctx context.Context, r *req.CreateTransactionReq) (*resp.TransactionResponse, error) {
	if err := validator.Validate(r); err != nil {
		return nil, errors.Wrap(response.WrapErrBadRequest(err), "APIUsecase.CreateTransaction.Validate")
	}

	t := &model.Transaction{
		OrderID: r.OrderID,
		Amount:  r.Amount,
		Status:  model.StatusPending,
		Method:  r.Method,
		Created: sharedmodel.Created{
			CreatedAt: getTimeNow(),
			CreatedBy: auth.GetEmail(ctx),
		},
	}

	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateTransaction.TxBegin")
	}
	defer func() {
		if txErr := u.repo.TxEnd(tx, err); txErr != nil {
			u.appLogger.ErrorContext(ctx, errors.Wrap(txErr, "APIUsecase.CreateTransaction.TxEnd").Error())
		}
	}()

	id, err := u.repo.InsertTransaction(ctx, tx, t)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateTransaction.InsertTransaction")
	}
	t.ID = id

	// Simulate payment processing (in production: call payment gateway)
	t.Status = model.StatusSucceeded
	if err = u.repo.UpdateTransactionStatus(ctx, tx, id, model.StatusSucceeded); err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateTransaction.UpdateStatus")
	}

	// Publish payment outcome for other services to react
	_ = u.broker.Publish(ctx, broker.SubjectPaymentSucceeded, map[string]interface{}{
		"transaction_id": id,
		"order_id":       r.OrderID,
		"amount":         r.Amount,
	})

	return toTransactionResponse(t), nil
}

func toTransactionResponse(t *model.Transaction) *resp.TransactionResponse {
	return &resp.TransactionResponse{
		ID:      t.ID,
		OrderID: t.OrderID,
		Amount:  t.Amount,
		Status:  t.Status,
		Method:  t.Method,
		CreatedResponse: sharedresp.CreatedResponse{
			CreatedAt: t.CreatedAt,
			CreatedBy: t.CreatedBy,
		},
		UpdatedResponse: sharedresp.UpdatedResponse{
			UpdatedAt: t.UpdatedAt,
			UpdatedBy: t.UpdatedBy,
		},
	}
}
