package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer trace.Tracer
}

func NewTracer(name string) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(name),
	}
}

type SpanContext struct {
	ctx  context.Context
	span trace.Span
}

func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, *SpanContext) {
	ctx, span := t.tracer.Start(ctx, spanName, opts...)
	return ctx, &SpanContext{
		ctx:  ctx,
		span: span,
	}
}

func (sc *SpanContext) End() {
	sc.span.End()
}

func (sc *SpanContext) EndWithError(err error) {
	if err != nil {
		sc.span.RecordError(err)
		sc.span.SetStatus(codes.Error, err.Error())
	}
	sc.span.End()
}

func (sc *SpanContext) SetAttributes(attrs ...attribute.KeyValue) {
	sc.span.SetAttributes(attrs...)
}

func (sc *SpanContext) AddEvent(name string, attrs ...attribute.KeyValue) {
	sc.span.AddEvent(name, trace.WithAttributes(attrs...))
}

func (sc *SpanContext) RecordError(err error, attrs ...attribute.KeyValue) {
	sc.span.RecordError(err, trace.WithAttributes(attrs...))
}

func (sc *SpanContext) SetStatus(code codes.Code, description string) {
	sc.span.SetStatus(code, description)
}

func (sc *SpanContext) Context() context.Context {
	return sc.ctx
}

func (sc *SpanContext) Span() trace.Span {
	return sc.span
}

func TraceFunction[T any](ctx context.Context, tracer *Tracer, spanName string, fn func(context.Context) (T, error)) (T, error) {
	ctx, spanCtx := tracer.Start(ctx, spanName)
	defer spanCtx.End()

	start := time.Now()
	result, err := fn(ctx)
	duration := time.Since(start)

	spanCtx.SetAttributes(
		attribute.String("function.name", spanName),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		spanCtx.RecordError(err)
		spanCtx.SetStatus(codes.Error, err.Error())
	} else {
		spanCtx.SetStatus(codes.Ok, "success")
	}

	return result, err
}

func TraceVoidFunction(ctx context.Context, tracer *Tracer, spanName string, fn func(context.Context) error) error {
	ctx, spanCtx := tracer.Start(ctx, spanName)
	defer spanCtx.End()

	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	spanCtx.SetAttributes(
		attribute.String("function.name", spanName),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		spanCtx.RecordError(err)
		spanCtx.SetStatus(codes.Error, err.Error())
	} else {
		spanCtx.SetStatus(codes.Ok, "success")
	}

	return err
}
