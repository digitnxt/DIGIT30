package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/digitnxt/DIGIT30/pkg/all"
)

func main() {
	// Initialize Prometheus metrics and start the metrics server.
	all.RegisterPrometheusMetrics()
	go all.StartMetricsServer(9464)

	// Initialize Jaeger tracer.
	shutdownTracer := all.InitTracer("identity", "http://jaeger:14268/api/traces")
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()

	// Create a new Gin router.
	r := gin.Default()

	// Setup Documentation endpoints (Swagger/OpenAPI).
	all.SetupDocumentation(r, "./docs/swagger.json")

	// Use InstrumentHandler from the unified library for your handlers.
	r.GET("/ping", all.InstrumentHandler(PingHandler))

	// Register your service with Consul using the discovery functions.
	err := all.RegisterService("identity-service", "identity", "identity", 8080, "http://identity:8080/ping")
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	// Start the HTTP server.
	r.Run() // listens on 0.0.0.0:8080
}

// PingHandler is a sample endpoint that records a business metric.
func PingHandler(c *gin.Context) {
	// Example business metric recording.
	who := "User-123"
	what := "Ping"
	why := "HealthCheck"
	when := time.Now().Format(time.RFC3339)
	how := "HTTP"
	where := "Server-1"
	whom := "Service-Identity"
	account := "Account-123"
	howmuch := "100"

	all.RecordBusinessMetric(who, what, why, when, how, where, whom, account, howmuch)

	c.String(http.StatusOK, "pong")
}