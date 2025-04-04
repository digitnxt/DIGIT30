// pkg/all/all.go
package all

import (
	"github.com/digitNxt/digit30/pkg/observability"
	"github.com/digitNxt/digit30/pkg/discovery"
	"github.com/digitNxt/digit30/pkg/documentation"
)

// Observability functions.
var (
	RegisterPrometheusMetrics = observability.RegisterPrometheusMetrics
	StartMetricsServer        = observability.StartMetricsServer
	InitTracer                = observability.InitTracer
	InstrumentHandler         = observability.InstrumentHandler
	RecordBusinessMetric      = observability.RecordBusinessMetric
)

// Discovery functions.
var (
	RegisterService   = discovery.RegisterService
	DeregisterService = discovery.DeregisterService
)

// Documentation functions.
var (
	SetupDocumentation = documentation.SetupDocumentation
)