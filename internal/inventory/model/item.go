package model

import sharedmodel "github.com/raflynagachi/go-rest-api-starter/internal/model"

type Item struct {
	ID       int64   `db:"id"`
	SKU      string  `db:"sku"`
	Name     string  `db:"name"`
	Quantity int64   `db:"quantity"`
	Price    float64 `db:"price"`
	sharedmodel.Created
	sharedmodel.Updated
	sharedmodel.Deleted
}
