package main

import (
	"context"
	"github.com/kdrkrgz/ecomm-micro/payment/config"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/adapters/db"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/adapters/grpc"
	"github.com/kdrkrgz/ecomm-micro/payment/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const (
	service     = "payment"
	environment = "dev"
	id          = 1
)

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "timestamp",
				"msg":            "message",
			},
		},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (c customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["traceId"] = span.SpanContext().TraceID().String()
	entry.Data["spanId"] = span.SpanContext().SpanID().String()
	entry.Data["context"] = span.SpanContext()
	return c.formatter.Format(entry)
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://jaeger-otel.jaeger.svc.cluster.local:14278/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	dbAdapter, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		log.Fatalf("Failed to connect database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
