package main

import (
	"github.com/kdrkrgz/ecomm-micro/order/config"
	"github.com/kdrkrgz/ecomm-micro/order/internal/adapters/db"
	"github.com/kdrkrgz/ecomm-micro/order/internal/adapters/grpc"
	"github.com/kdrkrgz/ecomm-micro/order/internal/adapters/payment"
	"github.com/kdrkrgz/ecomm-micro/order/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const (
	service     = "order"
	environment = "dev"
	id          = 1
)

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		},
			PrettyPrint: true,
		},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
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
	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (c customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	entry.Data["context"] = span.SpanContext()
	return c.formatter.Format(entry)
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}
