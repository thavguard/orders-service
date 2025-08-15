package models

type Item struct {
	Id          int    `db:"id" json:"id"`
	ChrtId      int    `db:"chrt_id" json:"chrt_id"`
	TrackNumber string `db:"track_number" json:"track_number"`
	Price       int    `db:"price" json:"price"`
	Rid         string `db:"rid" json:"rid"`
	Name        string `db:"name" json:"name"`
	Sale        int    `db:"sale" json:"sale"`
	Size        string `db:"size" json:"size"`
	Nm_id       int    `db:"nm_id" json:"nm_id"`
	TotalPrice  int    `db:"total_price" json:"total_price"`
	Brand       string `db:"brand" json:"brand"`
	Status      int    `db:"status" json:"status"`
	OrderId     int    `db:"order_id" json:"-"`
}
