package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
)

type Broker struct {
	topic  string
	reader *otelkafkakonsumer.Reader
	writer *kafka.Writer
}

func NewBroker(tp *trace.TracerProvider, topic string) *Broker {

	kafkaUrls := []string{os.Getenv("KAFKA_HOST")}

	r, err := otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: kafkaUrls,
			Topic:   topic,
			GroupID: "GroupID",
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String("opentel-autocommit-cg"),
			},
		),
	)

	if err != nil {
		log.Printf("ERROR IN INIT READER: %v\n", err)
	}

	DLQTopic := topic + "." + "errors"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkaUrls,
		Topic:   DLQTopic,
	})

	w.AllowAutoTopicCreation = true

	return &Broker{topic: topic, reader: r, writer: w}
}

func (b *Broker) Read(ctx context.Context) (*kafka.Message, error) {
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

func (b *Broker) Trace(ctx context.Context, message *kafka.Message) context.Context {
	ctx = b.reader.TraceConfig.Propagator.Extract(ctx, otelkafkakonsumer.NewMessageCarrier(message))

	return ctx
}

func (b *Broker) Close() (error, error) {
	readerError := b.reader.Close()
	writerError := b.writer.Close()

	return readerError, writerError

}
