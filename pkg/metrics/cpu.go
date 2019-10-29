package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

var cpuUsageDesc = prometheus.NewDesc("lxd_container_cpu_usage",
	"Container CPU Usage in Seconds",
	[]string{"container_name"}, nil,
)

func (collector *collector) collectCPUMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	cpuState lxdapi.ContainerStateCPU,
) {
	ch <- prometheus.MustNewConstMetric(
		cpuUsageDesc, prometheus.GaugeValue, float64(cpuState.Usage), containerName)
}
