package models

type Delivery struct {
	ID      int    `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name" validate:"required,alpha"`
	Phone   string `db:"phone" json:"phone" validate:"required,e164"`
	Zip     string `db:"zip" json:"zip" validate:"required,zipcode"`
	City    string `db:"city" json:"city" validate:"required,alpha"`
	Address string `db:"address" json:"address" validate:"required"`
	Region  string `db:"region" json:"region" validate:"required,alphanumunicode"`
	Email   string `db:"email" json:"email" validate:"required,email"`
	OrderID int    `db:"order_id" json:"order_id" validate:"number"`
}
