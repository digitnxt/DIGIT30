package observability

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // SDK tracer package, aliased as sdktrace
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace" // Import the trace package without alias
)
// Common Prometheus metrics
var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digit_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digit_http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"path"},
	)
	DurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "digit_http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	BusinessServiceMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digit_business_service_metric_total",
			Help: "Business service metrics with details like who, what, why, when, how, where, whom, account, howmuch",
		},
		[]string{"who", "what", "why", "when", "how", "where", "whom", "account", "howmuch"},
	)
)

// RegisterPrometheusMetrics registers the common metrics with Prometheus.
func RegisterPrometheusMetrics() {
	prometheus.MustRegister(RequestCounter, ErrorCounter, DurationHistogram, BusinessServiceMetric)
}

// StartMetricsServer starts an HTTP server on the specified port to expose Prometheus metrics.
func StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", port)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("failed to start Prometheus metrics server: %v", err)
		}
	}()
}

// InitTracer initializes a Jaeger tracer for the given service and returns a shutdown function.
// jaegerEndpoint should be in the form "http://jaeger:14268/api/traces"
func InitTracer(serviceName, jaegerEndpoint string) func(context.Context) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider( // Updated reference to sdktrace
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes("", attribute.String("service.name", serviceName))),
	)
	otel.SetTracerProvider(tp)
	// Start the Prometheus metrics server on port 9090
	StartMetricsServer(9090)
	return tp.Shutdown
}

// InstrumentHandler wraps a gin.HandlerFunc with Prometheus instrumentation.
func InstrumentHandler(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		timer := prometheus.NewTimer(DurationHistogram.WithLabelValues(path))
		defer timer.ObserveDuration()

		handlerFunc(c)

		statusCode := c.Writer.Status()
		RequestCounter.WithLabelValues(path).Inc()
		if statusCode >= 400 {
			ErrorCounter.WithLabelValues(path).Inc()
		}
	}
}

// RecordBusinessMetric records an event in the business service metric.
func RecordBusinessMetric(who, what, why, when, how, where, whom, account, howmuch string) {
	BusinessServiceMetric.WithLabelValues(who, what, why, when, how, where, whom, account, howmuch).Inc()
}

// StartSpan creates a new tracing span with the given name and returns the updated context and span.
// StartSpan creates a new tracing span with the given name and returns the updated context and span.
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("identity")
	return tracer.Start(ctx, spanName)
}
// TracingMiddleware is a Gin middleware that creates a tracing span for each incoming request.
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := StartSpan(c.Request.Context(), c.FullPath())
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.FullPath()),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}