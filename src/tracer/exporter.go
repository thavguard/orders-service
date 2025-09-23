package tracer

import (
	"context"

	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type SpanExporter interface {
	ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error
	Shutdown(ctx context.Context) error
}

func NewExporter(url string) (SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

}
