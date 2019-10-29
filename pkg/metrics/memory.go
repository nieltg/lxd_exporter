package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectMemoryMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	memState lxdapi.ContainerStateMemory,
) {
	collector.collectRAMMetrics(ch, containerName, memState)

	ch <- prometheus.MustNewConstMetric(swapUsageDesc,
		prometheus.GaugeValue, float64(memState.SwapUsage), containerName)
	ch <- prometheus.MustNewConstMetric(swapUsagePeakDesc,
		prometheus.GaugeValue, float64(memState.SwapUsagePeak), containerName)
}

func (collector *collector) collectRAMMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	memState lxdapi.ContainerStateMemory,
) {
	ch <- prometheus.MustNewConstMetric(
		memUsageDesc, prometheus.GaugeValue, float64(memState.Usage), containerName)
	ch <- prometheus.MustNewConstMetric(memUsagePeakDesc,
		prometheus.GaugeValue, float64(memState.UsagePeak), containerName)
}
