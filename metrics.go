package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "line_"

// Metrics implements the prometheus.Metrics interface and
// exposes gitea metrics for prometheus
type Metrics struct {
	ReceiveCount *prometheus.Desc
	SendCount    *prometheus.Desc
}

// NewMetrics returns a new Metrics with all prometheus.Desc initialized
func NewMetrics() Metrics {

	return Metrics{
		ReceiveCount: prometheus.NewDesc(
			namespace+"receive_count",
			"Number of receive count",
			nil, nil,
		),
		SendCount: prometheus.NewDesc(
			namespace+"send_count",
			"Number of send count",
			nil, nil,
		),
	}
}

// Describe returns all possible prometheus.Desc
func (c Metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ReceiveCount
	ch <- c.SendCount
}

// Collect returns the metrics with values
func (c Metrics) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.ReceiveCount,
		prometheus.GaugeValue,
		float64(ReceiveCount),
	)
	ch <- prometheus.MustNewConstMetric(
		c.SendCount,
		prometheus.GaugeValue,
		float64(SendCount),
	)
}
