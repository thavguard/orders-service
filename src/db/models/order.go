package models

import "time"

type Order struct {
	ID                int       `db:"id" json:"id,omitempty"`
	OrderUID          string    `db:"order_uid" json:"order_uid" validate:"required"`
	TrackNumber       string    `db:"track_number" json:"track_number" validate:"required"`
	Entry             string    `db:"entry" json:"entry" validate:"required"`
	Locale            string    `db:"locale" json:"locale" validate:"required,iso3166_1_alpha2"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerID        string    `db:"customer_id" json:"customer_id" validate:"required"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service" validate:"required"`
	Shardkey          string    `db:"shardkey" json:"shardkey" validate:"required,numeric"`
	SmID              int       `db:"sm_id" json:"sm_id" validate:"required,number"`
	DateCreated       time.Time `db:"date_created" json:"date_created" validate:"required"`
	OofShard          string    `db:"oof_shard" json:"oof_shard" validate:"required,numeric"`
	DeliveryID        int       `db:"delivery_id" json:"delivery_id" validate:"required,number"`
	PaymentID         int       `db:"payment_id" json:"payment_id" validate:"required,number"`
}
