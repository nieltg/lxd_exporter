package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectCPUMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	cpuState lxdapi.ContainerStateCPU,
) {
	ch <- prometheus.MustNewConstMetric(
		cpuUsageDesc, prometheus.GaugeValue, float64(cpuState.Usage), containerName)
}
