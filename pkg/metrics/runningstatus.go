package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var runningStatusDesc = prometheus.NewDesc("lxd_container_running_status",
	"Container Running Status",
	[]string{"container_name"}, nil,
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
