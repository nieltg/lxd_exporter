package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	diskUsageDesc = prometheus.NewDesc("lxd_container_disk_usage",
		"Container Disk Usage",
		[]string{"container_name", "disk_device"}, nil,
	)
	networkUsageDesc = prometheus.NewDesc("lxd_container_network_usage",
		"Container Network Usage",
		[]string{"container_name", "interface", "operation"}, nil,
	)
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
