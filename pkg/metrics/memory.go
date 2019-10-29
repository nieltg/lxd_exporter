package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	memUsageDesc = prometheus.NewDesc("lxd_container_mem_usage",
		"Container Memory Usage",
		[]string{"container_name"}, nil,
	)
	memUsagePeakDesc = prometheus.NewDesc("lxd_container_mem_usage_peak",
		"Container Memory Usage Peak",
		[]string{"container_name"}, nil,
	)

	swapUsageDesc = prometheus.NewDesc("lxd_container_swap_usage",
		"Container Swap Usage",
		[]string{"container_name"}, nil,
	)
	swapUsagePeakDesc = prometheus.NewDesc("lxd_container_swap_usage_peak",
		"Container Swap Usage Peak",
		[]string{"container_name"}, nil,
	)
)

func (collector *collector) collectMemoryMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	memState lxdapi.ContainerStateMemory,
) {
	collector.collectRAMMetrics(ch, containerName, memState)
	collector.collectSwapMetrics(ch, containerName, memState)
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

func (collector *collector) collectSwapMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	memState lxdapi.ContainerStateMemory,
) {
	ch <- prometheus.MustNewConstMetric(swapUsageDesc,
		prometheus.GaugeValue, float64(memState.SwapUsage), containerName)
	ch <- prometheus.MustNewConstMetric(swapUsagePeakDesc,
		prometheus.GaugeValue, float64(memState.SwapUsagePeak), containerName)
}
