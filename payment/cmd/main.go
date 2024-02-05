package main

import (
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"os"
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

func main() {

}
