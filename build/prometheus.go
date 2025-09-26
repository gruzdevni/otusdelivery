package build

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const metricsEndpoint = "/metrics"

func (b *Builder) prometheus() *prometheus.Registry {
	if b.prometheusRegistry != nil {
		return b.prometheusRegistry
	}

	b.prometheusRegistry = prometheus.NewRegistry()
	b.prometheusRegistry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), //nolint:exhaustruct
		collectors.NewGoCollector(),
	)

	return b.prometheusRegistry
}
