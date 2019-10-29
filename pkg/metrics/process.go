package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	processCountDesc = prometheus.NewDesc("lxd_container_process_count",
		"Container number of process Running",
		[]string{"container_name"}, nil,
	)
	containerPIDDesc = prometheus.NewDesc("lxd_container_pid",
		"Container PID",
		[]string{"container_name"}, nil,
	)
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
