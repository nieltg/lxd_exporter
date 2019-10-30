package metrics

import (
	"log"

	lxd "github.com/lxc/lxd/client"
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector collects metrics to be sent to Prometheus.
type collector struct {
	logger *log.Logger
	server lxd.InstanceServer
}

// NewCollector creates a new collector with logger and LXD connection.
func NewCollector(
	logger *log.Logger, server lxd.InstanceServer) prometheus.Collector {
	return &collector{logger: logger, server: server}
}

// Describe fills given channel with metrics descriptor.
func (collector *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cpuUsageDesc
	ch <- memUsageDesc
	ch <- memUsagePeakDesc
	ch <- swapUsageDesc
	ch <- swapUsagePeakDesc
	ch <- processCountDesc
	ch <- containerPIDDesc
	ch <- runningStatusDesc
	ch <- diskUsageDesc
	ch <- networkUsageDesc
}

// Collect fills given channel with metrics data.
func (collector *collector) Collect(ch chan<- prometheus.Metric) {
	containerNames, err := collector.server.GetContainerNames()
	if err != nil {
		collector.logger.Printf("Can't query container names: %s", err)
		return
	}

	for _, containerName := range containerNames {
		state, _, err := collector.server.GetContainerState(containerName)
		if err != nil {
			collector.logger.Printf(
				"Can't query container state for `%s`: %s", containerName, err)
			continue
		}

		collector.collectContainerMetrics(ch, containerName, state)
	}
}

func (collector *collector) collectContainerMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	state *lxdapi.ContainerState,
) {
	collector.collectCPUMetrics(ch, containerName, state.CPU)
	collector.collectMemoryMetrics(ch, containerName, state.Memory)
	collector.collectProcessMetrics(ch, containerName, state)
	collector.collectRunningStatusMetrics(ch, containerName, state.Status)
	collector.collectDiskMetrics(ch, containerName, state.Disk)
	collector.collectNetworkMetrics(ch, containerName, state.Network)
}
