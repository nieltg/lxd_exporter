package metrics

import (
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (collector *collector) collectNetworkMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	networkStates map[string]lxdapi.ContainerStateNetwork,
) {
	for interfaceName, state := range networkStates {
		networkMetrics := map[string]int64{
			"BytesReceived":   state.Counters.BytesReceived,
			"BytesSent":       state.Counters.BytesSent,
			"PacketsReceived": state.Counters.PacketsReceived,
			"PacketsSent":     state.Counters.PacketsSent,
		}

		for metricName, value := range networkMetrics {
			ch <- prometheus.MustNewConstMetric(
				networkUsageDesc, prometheus.GaugeValue, float64(value),
				containerName, interfaceName, metricName)
		}
	}
}
