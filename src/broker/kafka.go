package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

type Broker struct {
	topic  string
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewBroker(topic string) *Broker {

	kafkaUrls := []string{os.Getenv("KAFKA_HOST")}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaUrls,
		Topic:   topic,
		GroupID: "GroupID",
	})

	DLQTopic := topic + "." + "errors"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkaUrls,
		Topic:   DLQTopic,
	})

	w.AllowAutoTopicCreation = true

	return &Broker{topic: topic, reader: r, writer: w}
}

func (b *Broker) Read(ctx context.Context) (kafka.Message, error) {
	return b.reader.ReadMessage(ctx)
}

func (b *Broker) PushDQL(ctx context.Context, key string, dlqMessage DQLMessage, repetable string, maxRetries string) error {

	dlqMessageJSON, err := json.Marshal(dlqMessage)

	if err != nil {
		log.Printf("ERROR IN dlqMessageJSON: %v\n", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: dlqMessageJSON,
		Headers: []kafka.Header{
			protocol.Header{Key: "repetable", Value: []byte(repetable)},
			protocol.Header{Key: "max_retries", Value: []byte(maxRetries)},
		},
	}

	return b.writer.WriteMessages(ctx, message)
}

func (b *Broker) Close() (error, error) {
	readerError := b.reader.Close()
	writerError := b.writer.Close()

	return readerError, writerError

}
