package broker

import "orders/src/db/models"

type OrderMessage struct {
	*models.Order

	Delivery models.Delivery `json:"delivery"`
	Payment  models.Payment  `json:"payment"`
	Items    []models.Item   `json:"items"`
}
