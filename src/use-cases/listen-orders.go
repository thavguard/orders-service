package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orders/src/db/repositories"
	"orders/src/internal/broker"
	"orders/src/mycache"
	"strconv"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/segmentio/kafka-go"
)

type ListenOrdersSerivce struct {
	DbService *repositories.DBRepository
	Cache     *mycache.RedisService
}

func (service *ListenOrdersSerivce) saveOrderToDb(ctx context.Context, order *broker.OrderMessage) (int, error) {
	delivery, err := service.DbService.CreateDelivery(ctx, &order.Delivery)

	if err != nil {
		return 0, err
	}

	payment, err := service.DbService.CreatePayment(ctx, &order.Payment)

	if err != nil {
		return 0, err
	}

	orderDto := order.Order
	orderDto.DeliveryId = delivery.Id
	orderDto.PaymentId = payment.Id

	orderDb, err := service.DbService.CreateOrder(ctx, orderDto)

	if err != nil {
		return 0, err
	}

	for i := range order.Items {
		currentItem := &order.Items[i]

		currentItem.OrderId = orderDb.Id

		_, err := service.DbService.CreateItem(ctx, currentItem)

		if err != nil {
			return 0, err
		}

	}

	orderDbId := orderDb.Id

	// Save to cache

	redisKey := "order_" + strconv.Itoa(orderDbId)

	if err = service.Cache.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   redisKey,
		Value: &orderDb,
		TTL:   time.Hour,
	}); err != nil {
		fmt.Printf("ERROR IN REDIS LISTEN ORDER: %v\n", err)
	}

	return orderDbId, nil

}

func (service *ListenOrdersSerivce) ListenOrders(ctx context.Context) *kafka.Reader {
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
				orderId, err := service.saveOrderToDb(ctx, obj)

				fmt.Printf("NEW ORDER SAVED: %v\n", orderId)

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
