package metrics

import (
	"log"

	lxd "github.com/lxc/lxd/client"
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

// Describe ...
func (collector *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cpuUsageDesc
	ch <- memUsageDesc
	ch <- memUsagePeakDesc
	ch <- swapUsageDesc
	ch <- swapUsagePeakDesc
}

// Collect ...
func (collector *collector) Collect(ch chan<- prometheus.Metric) {
	names, err := collector.server.GetContainerNames()
	if err != nil {
		collector.logger.Printf("Can't query container names: %s", err)
		return
	}

	for _, name := range names {
		state, _, err := collector.server.GetContainerState(name)
		if err != nil {
			collector.logger.Printf(
				"Can't query container state for `%s`: %s", name, err)
			break
		}

		ch <- prometheus.MustNewConstMetric(
			cpuUsageDesc, prometheus.GaugeValue, float64(state.CPU.Usage), name)
		ch <- prometheus.MustNewConstMetric(
			memUsageDesc, prometheus.GaugeValue, float64(state.Memory.Usage), name)
		ch <- prometheus.MustNewConstMetric(
			memUsagePeakDesc, prometheus.GaugeValue, float64(state.Memory.UsagePeak),
			name)
		ch <- prometheus.MustNewConstMetric(
			swapUsageDesc, prometheus.GaugeValue, float64(state.Memory.SwapUsage),
			name)
		ch <- prometheus.MustNewConstMetric(
			swapUsagePeakDesc, prometheus.GaugeValue, float64(
				state.Memory.SwapUsagePeak), name)
	}
}
