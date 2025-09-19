package models

type Item struct {
	ID          int    `db:"id" json:"id,omitempty"`
	ChrtID      int    `db:"chrt_id" json:"chrt_id" validate:"required,number"`
	TrackNumber string `db:"track_number" json:"track_number" validate:"required,alpha"`
	Price       int    `db:"price" json:"price" validate:"required,number"`
	Rid         string `db:"rid" json:"rid" validate:"required,alphanum"`
	Name        string `db:"name" json:"name" validate:"required,alpha"`
	Sale        int    `db:"sale" json:"sale" validate:"required,number"`
	Size        string `db:"size" json:"size" validate:"required,numeric"`
	NmID        int    `db:"nm_id" json:"nm_id" validate:"required,number"`
	TotalPrice  int    `db:"total_price" json:"total_price" validate:"required,number"`
	Brand       string `db:"brand" json:"brand" validate:"required,alphanum"`
	Status      int    `db:"status" json:"status" validate:"required,number"`
	OrderID     int    `db:"order_id" json:"order_id" validate:"number"`
}
