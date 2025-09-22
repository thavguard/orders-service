package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(url string, serviceName string) (*tracesdk.TracerProvider, error) {
	exporter, err := NewExporter(url)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(tp) // !!!!!!!!!!!

	return tp, nil
}
