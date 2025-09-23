package filldata

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type brokerSerive struct {
	Writer *kafka.Writer
}

func createProvider() *brokerSerive {
	fmt.Println(os.Getenv("KAFKA_HOST"))

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		Topic:   "orders",
	})

	return &brokerSerive{Writer: w}
}

func (broker *brokerSerive) close() {
	broker.Writer.Close()
}

func (broker *brokerSerive) writeMessages(ctx context.Context, messages []kafka.Message) {
	err := broker.Writer.WriteMessages(ctx, messages...)

	if err != nil {
		fmt.Printf("failed to write messages: %v\n", err)
	}

	if err := broker.Writer.Close(); err != nil {
		fmt.Printf("failed to close writer: %v\n", err)
	}
}

func getDeliveryJSON() string {
	json := `{
		"name": "Sergey",
		"phone": "+79235858077",
		"zip": "12345",
		"city": "Moscow",
		"address": "Автозаводская 23с16",
		"region": "Moscow",
		"email": "test@test.com"
		}`

	return json
}

func getPaymentJSON() string {
	json := `{
		"transaction": "transaction",
		"request_id": "1234",
		"currency": "RUB",
		"provider": "provider",
		"amount": 100,
		"payment_dt": 100,
		"bank": "Т-Банк",
		"delivery_cost": 100,
		"goods_total": 100,
		"custom_fee": 100
		}`

	return json
}

func getItemsJSON() string {
	json := `{
		"chrt_id": 1,
		"track_number": "HEYEDSADSa",
		"price": 200,
		"rid": "RIDRIDRID",
		"name": "Sergey",
		"sale": 50,
		"size": "100",
		"nm_id": 123,
		"total_price": 500,
		"brand": "SUPERBRAND",
		"status": 1,
		"order_id": 1
		}`

	return json
}

func getOrderJSON() string {
	json := `{
		"order_uid": "SKDMALSDMLSA",
		"track_number": "KEIEIIEIE",
		"entry": "HEY",
		"locale": "US",
		"internal_signature": "HEYSADA",
		"customer_id": "123",
		"delivery_service": "WBTWIN",
		"shardkey": "44412",
		"sm_id": 123,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "123312",
		"delivery_id": 1,
		"payment_id": 1
		}`

	return json
}

func writeMessages(ctx context.Context, brokerService *brokerSerive) {

	payment := getPaymentJSON()
	delivery := getDeliveryJSON()
	items := getItemsJSON()
	order := getOrderJSON()

	var messages []kafka.Message

	for range 100 {
		messagePayment := kafka.Message{
			Key:   []byte("payment"),
			Value: []byte(payment),
		}

		messageDelivery := kafka.Message{
			Key:   []byte("delivery"),
			Value: []byte(delivery),
		}

		messageItem := kafka.Message{
			Key:   []byte("item"),
			Value: []byte(items),
		}

		messageOrder := kafka.Message{
			Key:   []byte("order"),
			Value: []byte(order),
		}

		messages = append(messages, messagePayment, messageDelivery, messageItem, messageOrder)
	}

	brokerService.writeMessages(ctx, messages)
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Print("No .env file found")
	}
}

func FillData(ctx context.Context) {

	// Init provider
	brokerService := *createProvider()
	defer brokerService.close()

	// Write messages
	writeMessages(ctx, &brokerService)

}
