package broker

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

func CreateConsumer(ctx context.Context, topic string) *kafka.Reader {

	kafkaUrls := []string{os.Getenv("KAFKA_HOST")}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaUrls,
		Topic:   topic,
		GroupID: "GroupID",
	})

	return reader

}
