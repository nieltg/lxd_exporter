package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectProcessMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	state *lxdapi.ContainerState,
) {
	ch <- prometheus.MustNewConstMetric(
		processCountDesc, prometheus.GaugeValue, float64(state.Processes),
		containerName)
	ch <- prometheus.MustNewConstMetric(
		containerPIDDesc, prometheus.GaugeValue, float64(state.Pid), containerName)
}
