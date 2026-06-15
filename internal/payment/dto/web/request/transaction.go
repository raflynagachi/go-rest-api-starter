package request

type CreateTransactionReq struct {
	OrderID int64   `json:"order_id" validate:"required,min=1"`
	Amount  float64 `json:"amount"   validate:"required,gt=0"`
	Method  string  `json:"method"   validate:"required,oneof=credit_card bank_transfer wallet"`
}
