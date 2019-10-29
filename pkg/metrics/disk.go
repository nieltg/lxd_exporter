package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectDiskMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	diskStates map[string]lxdapi.ContainerStateDisk,
) {
	for diskName, state := range diskStates {
		ch <- prometheus.MustNewConstMetric(diskUsageDesc,
			prometheus.GaugeValue, float64(state.Usage), containerName, diskName)
	}
}
