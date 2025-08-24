package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orders/src/db/queries"
	"orders/src/internal/broker"

	"github.com/segmentio/kafka-go"
)

func saveOrderToDb(ctx context.Context, dbService *queries.DBService, order *broker.OrderMessage) (int, error) {
	delivery, err := dbService.CreateDelivery(ctx, &order.Delivery)

	if err != nil {
		return 0, err
	}

	payment, err := dbService.CreatePayment(ctx, &order.Payment)

	if err != nil {
		return 0, err
	}

	orderDto := order.Order
	orderDto.DeliveryId = delivery.Id
	orderDto.PaymentId = payment.Id

	orderDb, err := dbService.CreateOrder(ctx, orderDto)

	if err != nil {
		return 0, err
	}

	for i := range order.Items {
		currentItem := &order.Items[i]

		currentItem.OrderId = orderDb.Id

		_, err := dbService.CreateItem(ctx, currentItem)

		if err != nil {
			return 0, err
		}

	}

	return orderDb.Id, nil

}

func ListenOrders(ctx context.Context, dbService *queries.DBService) *kafka.Reader {
	reader := broker.CreateConsumer(ctx, "orders")

	go func() {
		for {
			m, err := reader.ReadMessage(ctx)

			if err != nil {
				log.Fatal("Failed while listen topic:", err)
				break
			}

			obj := &broker.OrderMessage{}

			err = json.Unmarshal(m.Value, obj)

			if err != nil {
				fmt.Printf("Error in parse JSON %v\n", err)
			} else {
				_, err := saveOrderToDb(ctx, dbService, obj)

				if err != nil {
					fmt.Printf("SOME ERROR: %v\n", err)
				}

			}

		}

		if err := reader.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}

	}()

	return reader
}
