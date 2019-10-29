package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectRunningStatusMetrics(
	ch chan<- prometheus.Metric, containerName string, status string) {
	runningStatus := 0
	if status == "Running" {
		runningStatus = 1
	}

	ch <- prometheus.MustNewConstMetric(
		runningStatusDesc, prometheus.GaugeValue, float64(runningStatus),
		containerName)
}
