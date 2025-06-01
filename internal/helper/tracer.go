package helper

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/trace"
)

var TraceServiceName = "SSO"

func NewGrpcExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	return otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	))
}

func SpanStart(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer(TraceServiceName).Start(ctx, name, opts...)
}

func SpanAttr(attr ...attribute.KeyValue) trace.SpanStartOption {
	return trace.WithAttributes(attr...)
}

func SpanError(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}
