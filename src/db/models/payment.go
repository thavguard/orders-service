package models

type Payment struct {
	ID           int    `db:"id" json:"id,omitempty"`
	Transaction  string `db:"transaction" json:"transaction" validate:"required,alphanumunicode"`
	RequestID    string `db:"request_id" json:"request_id" validate:"numeric"`
	Currency     string `db:"currency" json:"currency" validate:"required,iso4217"`
	Provider     string `db:"provider" json:"provider" validate:"required,alphanumunicode"`
	Amount       int    `db:"amount" json:"amount" validate:"required,number"`
	PaymentDt    int    `db:"payment_dt" json:"payment_dt" validate:"required,number"`
	Bank         string `db:"bank" json:"bank" validate:"required"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost" validate:"required,number"`
	GoodsTotal   int    `db:"goods_total" json:"goods_total" validate:"required,number"`
	CustomFee    int    `db:"custom_fee" json:"custom_fee" validate:"required,number"`
	OrderID      int    `db:"order_id" json:"order_id" validate:"number"`
}
