package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orders/src/db/queries"
	"orders/src/internal/broker"
)

func saveOrderToDb(ctx context.Context, dbService *queries.DBService, order *broker.OrderMessage) (int, error) {
	fmt.Println("START DB SERVICE")
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

	orderDb, err := dbService.CreateOrder(ctx, &orderDto)

	if err != nil {
		return 0, err
	}

	for i := range order.Items {
		currentItem := &order.Items[i]

		fmt.Printf("ORDER ID IN ITEM: %v\n", orderDb.Id)
		currentItem.OrderId = orderDb.Id

		_, err := dbService.CreateItem(ctx, currentItem)

		if err != nil {
			return 0, err
		}

	}

	return orderDb.Id, nil

}

func ListenOrders(ctx context.Context, dbService *queries.DBService) {
	reader := broker.CreateConsumer(ctx, "orders")

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
			fmt.Println(obj)
			orderID, err := saveOrderToDb(ctx, dbService, obj)

			if err != nil {
				fmt.Printf("SOME ERROR: %v\n", err)
			}

			fmt.Printf("ORDER ID: %v\n", orderID)

		}

	}

	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
