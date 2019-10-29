package metrics

import (
	"log"

	lxd "github.com/lxc/lxd/client"
	lxdapi "github.com/lxc/lxd/shared/api"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector ...
type collector struct {
	logger *log.Logger
	server lxd.InstanceServer
}

// NewCollector ...
func NewCollector(
	logger *log.Logger, server lxd.InstanceServer) prometheus.Collector {
	return &collector{logger: logger, server: server}
}

var cpuUsageDesc = prometheus.NewDesc("lxd_container_cpu_usage",
	"Container CPU Usage in Seconds",
	[]string{"container_name"}, nil,
)
var memUsageDesc = prometheus.NewDesc("lxd_container_mem_usage",
	"Container Memory Usage",
	[]string{"container_name"}, nil,
)
var memUsagePeakDesc = prometheus.NewDesc("lxd_container_mem_usage_peak",
	"Container Memory Usage Peak",
	[]string{"container_name"}, nil,
)
var swapUsageDesc = prometheus.NewDesc("lxd_container_swap_usage",
	"Container Swap Usage",
	[]string{"container_name"}, nil,
)
var swapUsagePeakDesc = prometheus.NewDesc("lxd_container_swap_usage_peak",
	"Container Swap Usage Peak",
	[]string{"container_name"}, nil,
)
var processCountDesc = prometheus.NewDesc("lxd_container_process_count",
	"Container number of process Running",
	[]string{"container_name"}, nil,
)
var containerPIDDesc = prometheus.NewDesc("lxd_container_pid",
	"Container PID",
	[]string{"container_name"}, nil,
)
var runningStatusDesc = prometheus.NewDesc("lxd_container_running_status",
	"Container Running Status",
	[]string{"container_name"}, nil,
)
var diskUsageDesc = prometheus.NewDesc("lxd_container_disk_usage",
	"Container Disk Usage",
	[]string{"container_name", "disk_device"}, nil,
)
var networkUsageDesc = prometheus.NewDesc("lxd_container_network_usage",
	"Container Network Usage",
	[]string{"container_name", "interface", "operation"}, nil,
)

// Describe ...
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

// Collect ...
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
			break
		}

		collector.collectContainerMetrics(ch, containerName, state)
	}
}

func (collector *collector) collectContainerMetrics(
	ch chan<- prometheus.Metric,
	containerName string,
	state *lxdapi.ContainerState,
) {
	ch <- prometheus.MustNewConstMetric(
		cpuUsageDesc, prometheus.GaugeValue, float64(state.CPU.Usage), containerName)
	ch <- prometheus.MustNewConstMetric(
		memUsageDesc, prometheus.GaugeValue, float64(state.Memory.Usage), containerName)
	ch <- prometheus.MustNewConstMetric(
		memUsagePeakDesc, prometheus.GaugeValue, float64(state.Memory.UsagePeak),
		containerName)
	ch <- prometheus.MustNewConstMetric(
		swapUsageDesc, prometheus.GaugeValue, float64(state.Memory.SwapUsage),
		containerName)
	ch <- prometheus.MustNewConstMetric(
		swapUsagePeakDesc, prometheus.GaugeValue, float64(
			state.Memory.SwapUsagePeak), containerName)
	ch <- prometheus.MustNewConstMetric(
		processCountDesc, prometheus.GaugeValue, float64(state.Processes), containerName)
	ch <- prometheus.MustNewConstMetric(
		containerPIDDesc, prometheus.GaugeValue, float64(state.Pid), containerName)

	runningStatus := 0
	if state.Status == "Running" {
		runningStatus = 1
	}
	ch <- prometheus.MustNewConstMetric(
		runningStatusDesc, prometheus.GaugeValue, float64(runningStatus), containerName)

	for diskName, diskState := range state.Disk {
		ch <- prometheus.MustNewConstMetric(diskUsageDesc,
			prometheus.GaugeValue, float64(diskState.Usage), containerName, diskName)
	}

	for ethName, netState := range state.Network {
		networkMetrics := map[string]int64{
			"BytesReceived":   netState.Counters.BytesReceived,
			"BytesSent":       netState.Counters.BytesSent,
			"PacketsReceived": netState.Counters.PacketsReceived,
			"PacketsSent":     netState.Counters.PacketsSent,
		}

		for metricName, value := range networkMetrics {
			ch <- prometheus.MustNewConstMetric(networkUsageDesc,
				prometheus.GaugeValue, float64(value), containerName, ethName, metricName)
		}
	}
}
