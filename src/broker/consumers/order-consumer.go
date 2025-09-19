package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"orders/src/broker"
	"orders/src/db/models"
	"orders/src/service"
	"sync"

	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	broker          *broker.Broker
	orderService    *service.OrderService
	deliveryService *service.DeliveryService
	itemService     *service.ItemService
	paymentService  *service.PaymentService
}

func NewOrderConsumer(orderService *service.OrderService,
	deliveryService *service.DeliveryService,
	itemService *service.ItemService,
	paymentService *service.PaymentService) *OrderConsumer {

	broker := broker.NewBroker("orders")

	return &OrderConsumer{broker: broker, orderService: orderService, deliveryService: deliveryService, itemService: itemService, paymentService: paymentService}
}

func (c *OrderConsumer) handleMessage(ctx context.Context, msg kafka.Message) {

	switch string(msg.Key) {

	case "order":
		fmt.Printf("ORDER: %v\n", msg.Value)

		var order models.Order

		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("ERROR IN CASE ORDER KAFKA: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err := c.broker.PushDQL(ctx, "order", dlq, "false", "0"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

			return
		}

		if _, err := c.orderService.CreateOrder(ctx, order); err != nil {
			log.Printf("ERROR IN CASE ORDER SERVICE: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err := c.broker.PushDQL(ctx, "order", dlq, "true", "5"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}
		}

	case "payment":
		fmt.Printf("PAYMENT: %v\n", msg.Value)

		var payment models.Payment

		if err := json.Unmarshal(msg.Value, &payment); err != nil {
			log.Printf("ERROR IN CASE PAYMENT KAFKA: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err = c.broker.PushDQL(ctx, "order", dlq, "false", "0"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

			return
		}

		if _, err := c.paymentService.CreatePayment(ctx, &payment); err != nil {
			log.Printf("ERROR IN CASE PAYMENT SERVICE: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err = c.broker.PushDQL(ctx, "order", dlq, "true", "5"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

		}

	case "item":
		fmt.Printf("ITEM: %v\n", msg.Value)

		var item models.Item

		if err := json.Unmarshal(msg.Value, &item); err != nil {
			log.Printf("ERROR IN CASE ITEM KAFKA: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err = c.broker.PushDQL(ctx, "order", dlq, "false", "0"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

			return
		}

		if _, err := c.itemService.CreateItem(ctx, &item); err != nil {
			log.Printf("ERROR IN CASE ITEM SERVICE: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err = c.broker.PushDQL(ctx, "order", dlq, "true", "5"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

		}

	case "delivery":
		fmt.Printf("DELIVERY: %v\n", msg.Value)

		var delivery models.Delivery

		if err := json.Unmarshal(msg.Value, &delivery); err != nil {
			log.Printf("ERROR IN CASE DELIVERY KAFKA: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err := c.broker.PushDQL(ctx, "order", dlq, "false", "0"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}

			return
		}

		if _, err := c.deliveryService.CreateDelivery(ctx, &delivery); err != nil {
			log.Printf("ERROR IN CASE DELIVERY SERVICE: %v\n", err)

			dlq := broker.DQLMessage{
				Origin: msg,
				Reason: err.Error(),
			}

			if err := c.broker.PushDQL(ctx, "order", dlq, "true", "5"); err != nil {
				log.Printf("EROR IN PushDQL: %v\n", err)
			}
		}

	}
}

func (c *OrderConsumer) Run(ctx context.Context) {

	var wg sync.WaitGroup

	const maxWorkers = 20 // TODO: настроить кол-во пулов к БД
	sem := make(chan struct{}, maxWorkers)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("CTX IS DONE, WAITING FOR GORUTINES...\n")
				wg.Wait()
				log.Printf("ALL GORUTINES IS DONE! BYE BYE\n")
				return

			default:
				message, err := c.broker.Read(ctx)
				if err != nil {
					log.Printf("error reading message: %v", err)
					continue
				}

				sem <- struct{}{}
				wg.Add(1)

				go func(msg kafka.Message) {
					defer wg.Done()
					defer func() {
						<-sem
					}()

					c.handleMessage(ctx, msg)

				}(message)

			}
		}
	}()
}

func (c *OrderConsumer) Close() (error, error) {
	return c.broker.Close()
}
