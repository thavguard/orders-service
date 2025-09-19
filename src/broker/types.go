package broker

import (
	"orders/src/db/models"

	"github.com/segmentio/kafka-go"
)

type OrderMessage struct {
	models.Order

	Delivery models.Delivery `json:"delivery"`
	Payment  models.Payment  `json:"payment"`
	Items    []models.Item   `json:"items"`
}

type DQLMessage struct {
	Origin kafka.Message `json:"origin"`
	Reason string        `json:"reason"`
}
