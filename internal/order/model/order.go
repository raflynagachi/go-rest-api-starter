package model

import sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          int64       `db:"id"`
	UserEmail   string      `db:"user_email"`
	Status      OrderStatus `db:"status"`
	TotalAmount float64     `db:"total_amount"`
	sharedmodel.Created
	sharedmodel.Updated
	sharedmodel.Deleted
}

type OrderItem struct {
	ID      int64   `db:"id"`
	OrderID int64   `db:"order_id"`
	SKU     string  `db:"sku"`
	Qty     int64   `db:"qty"`
	Price   float64 `db:"price"`
}
