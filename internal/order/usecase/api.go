package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/internal/apperror"
	req "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/request"
	resp "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/response"
	"github.com/raflynagachi/go-rest-api-starter/internal/order/model"
	sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"
	sharedresp "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/response"
	paginationutil "github.com/raflynagachi/go-rest-api-starter/internal/util/pagination"
	"github.com/raflynagachi/go-rest-api-starter/pkg/auth"
	"github.com/raflynagachi/go-rest-api-starter/pkg/broker"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
	"github.com/raflynagachi/go-rest-api-starter/pkg/validator"
)

// inventoryItemResp is the shape returned by GET /items?sku=X from the inventory service.
type inventoryItemResp struct {
	Data struct {
		ID    int64   `json:"id"`
		Price float64 `json:"price"`
	} `json:"data"`
}

// paymentChargeResp is the shape returned by POST /transactions from the payment service.
type paymentChargeResp struct {
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

func (u *APIUsecaseImpl) GetOrders(ctx context.Context, filter req.OrderFilter) (*sharedresp.ListResponse, error) {
	filter.Pagination.Validate()

	orders, err := u.repo.GetOrders(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetOrders")
	}

	count, err := u.repo.CountOrders(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetOrders.CountOrders")
	}

	data := make([]*resp.OrderResponse, 0, len(orders))
	for _, o := range orders {
		data = append(data, toOrderResponse(o, nil))
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

func (u *APIUsecaseImpl) GetOrderByID(ctx context.Context, id int64) (*resp.OrderResponse, error) {
	order, err := u.repo.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, errors.Wrap(response.WrapErrNotFound(err), "APIUsecase.GetOrderByID")
		}
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetOrderByID")
	}

	items, err := u.repo.GetOrderItems(ctx, id)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.GetOrderByID.GetOrderItems")
	}

	return toOrderResponse(order, items), nil
}

func (u *APIUsecaseImpl) CreateOrder(ctx context.Context, r *req.CreateOrderReq) (*resp.OrderResponse, error) {
	if err := validator.Validate(r); err != nil {
		return nil, errors.Wrap(response.WrapErrBadRequest(err), "APIUsecase.CreateOrder.Validate")
	}

	// Step 1: Verify each item's price and availability via Inventory service (HTTP).
	// Replace with gRPC once proto is generated: inventoryGRPC.CheckStock(ctx, &inventorypb.CheckStockRequest{...})
	orderItems := make([]*model.OrderItem, 0, len(r.Items))
	var total float64
	for _, item := range r.Items {
		var invResp inventoryItemResp
		path := fmt.Sprintf("/items?sku=%s", item.SKU)
		if err := u.inventoryClient.Get(ctx, path, &invResp); err != nil {
			return nil, errors.Wrap(response.WrapErrBadRequest(
				fmt.Errorf("item %s not available: %w", item.SKU, err),
			), "APIUsecase.CreateOrder.CheckInventory")
		}

		lineTotal := invResp.Data.Price * float64(item.Qty)
		total += lineTotal
		orderItems = append(orderItems, &model.OrderItem{
			SKU:   item.SKU,
			Qty:   item.Qty,
			Price: invResp.Data.Price,
		})
	}

	// Step 2: Persist order + items in a transaction.
	userEmail := auth.GetEmail(ctx)
	order := &model.Order{
		UserEmail:   userEmail,
		Status:      model.StatusPending,
		TotalAmount: total,
		Created: sharedmodel.Created{
			CreatedAt: getTimeNow(),
			CreatedBy: userEmail,
		},
	}

	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateOrder.TxBegin")
	}
	defer func() {
		if txErr := u.repo.TxEnd(tx, err); txErr != nil {
			u.appLogger.ErrorContext(ctx, errors.Wrap(txErr, "APIUsecase.CreateOrder.TxEnd").Error())
		}
	}()

	orderID, err := u.repo.InsertOrder(ctx, tx, order)
	if err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateOrder.InsertOrder")
	}
	order.ID = orderID

	for _, item := range orderItems {
		item.OrderID = orderID
		if err = u.repo.InsertOrderItem(ctx, tx, item); err != nil {
			return nil, errors.Wrap(response.WrapErrInternalServer(err), "APIUsecase.CreateOrder.InsertOrderItem")
		}
	}

	// Step 3: Charge via Payment service (HTTP).
	// Replace with gRPC: paymentGRPC.Charge(ctx, &paymentpb.ChargeRequest{...})
	var payResp paymentChargeResp
	chargeBody := map[string]interface{}{
		"order_id": orderID,
		"amount":   total,
		"method":   "credit_card",
	}
	if err = u.paymentClient.Post(ctx, "/transactions", chargeBody, &payResp); err != nil {
		return nil, errors.Wrap(response.WrapErrInternalServer(
			fmt.Errorf("payment failed: %w", err),
		), "APIUsecase.CreateOrder.Charge")
	}

	// Step 4: Publish event for async downstream work (e.g., deduct inventory stock).
	_ = u.broker.Publish(ctx, broker.SubjectOrderCreated, map[string]interface{}{
		"order_id":       orderID,
		"user_email":     userEmail,
		"total_amount":   total,
		"payment_tx_id":  payResp.Data.ID,
	})

	order.Status = model.StatusConfirmed
	return toOrderResponse(order, orderItems), nil
}

func toOrderResponse(o *model.Order, items []*model.OrderItem) *resp.OrderResponse {
	itemResps := make([]*resp.OrderItemResponse, 0, len(items))
	for _, i := range items {
		itemResps = append(itemResps, &resp.OrderItemResponse{
			SKU:   i.SKU,
			Qty:   i.Qty,
			Price: i.Price,
		})
	}
	return &resp.OrderResponse{
		ID:          o.ID,
		UserEmail:   o.UserEmail,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		Items:       itemResps,
		CreatedResponse: sharedresp.CreatedResponse{
			CreatedAt: o.CreatedAt,
			CreatedBy: o.CreatedBy,
		},
		UpdatedResponse: sharedresp.UpdatedResponse{
			UpdatedAt: o.UpdatedAt,
			UpdatedBy: o.UpdatedBy,
		},
	}
}
