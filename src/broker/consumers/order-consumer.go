package consumers

import (
	"context"
	"encoding/json"
	"log"
	"orders/src/broker"
	"orders/src/db/models"
	"orders/src/metrics"
	"orders/src/service"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/semaphore"
)

type OrderConsumer struct {
	broker          *broker.Broker
	orderService    service.OrderService
	deliveryService service.DeliveryService
	itemService     service.ItemService
	paymentService  service.PaymentService
	metrics         *metrics.Metrics
	tp              *trace.TracerProvider
}

func NewOrderConsumer(metrics *metrics.Metrics, tp *trace.TracerProvider, orderService service.OrderService,
	deliveryService service.DeliveryService,
	itemService service.ItemService,
	paymentService service.PaymentService) *OrderConsumer {

	broker := broker.NewBroker(tp, "orders")

	return &OrderConsumer{tp: tp, broker: broker, orderService: orderService, deliveryService: deliveryService, itemService: itemService, paymentService: paymentService, metrics: metrics}
}

func (c *OrderConsumer) handleMessage(ctx context.Context, msg *kafka.Message) {

	switch string(msg.Key) {

	case "order":

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

			return

		}

		c.metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "success").Inc()

	case "payment":

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

			return
		}

		c.metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "success").Inc()

	case "item":

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

			return

		}

		c.metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "success").Inc()

	case "delivery":

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()

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

			c.metrics.KafkaMessagesDLQ.WithLabelValues(msg.Topic).Inc()
			return
		}

		c.metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "success").Inc()

	}

}

func (c *OrderConsumer) Run(ctx context.Context) {

	var wg sync.WaitGroup

	const maxWorkers = 20 // TODO: настроить кол-во пулов к БД
	sem := semaphore.NewWeighted(maxWorkers)

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
					c.metrics.KafkaMessagesConsumed.WithLabelValues(message.Topic, "error").Inc()

					continue
				}

				wg.Add(1)

				go func(msg *kafka.Message) {
					defer wg.Done()

					if err := sem.Acquire(ctx, 1); err != nil {
						log.Printf("Error sem.Acquire: %v\n", err)
						return
					}

					defer func() {
						sem.Release(1)
					}()

					ctx = c.broker.Trace(ctx, msg)
					tr := c.tp.Tracer("orders-consumer")
					_, span := tr.Start(ctx, "handle-order")

					c.handleMessage(ctx, msg)

					span.End()

				}(message)

			}
		}
	}()
}

func (c *OrderConsumer) Close() (error, error) {
	return c.broker.Close()
}
