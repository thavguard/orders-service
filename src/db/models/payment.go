package models

type Payment struct {
	Id           int    `db:"id" json:"id"`
	Transaction  string `db:"transaction" json:"transaction"`
	RequestId    string `db:"request_id" json:"request_id"`
	Currency     string `db:"currency" json:"currency"`
	Provider     string `db:"provider" json:"provider"`
	Amount       int    `db:"amount" json:"amount"`
	PaymentDt    int    `db:"payment_dt" json:"payment_dt"`
	Bank         string `db:"bank" json:"bank"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost"`
	GoodsTotal   int    `db:"goods_total" json:"goods_total"`
	CustomFee    int    `db:"custom_fee" json:"custom_fee"`
}
