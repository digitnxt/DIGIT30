// pkg/all/all.go
package all

import (
	"github.com/digitnxt/digit/pkg/observability"
	"github.com/digitnxt/digit/pkg/discovery"
	"github.com/digitnxt/digit/pkg/documentation"
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