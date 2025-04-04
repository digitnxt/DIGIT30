// pkg/all/all.go
package all

import (
	"github.com/digitnxt/DIGIT30/pkg/observability"
	"github.com/digitnxt/DIGIT30/pkg/discovery"
	"github.com/digitnxt/DIGIT30/pkg/documentation"
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