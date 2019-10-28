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

// Describe ...
func (collector *collector) Describe(c chan<- *prometheus.Desc) {
}

// Collect ...
func (collector *collector) Collect(c chan<- prometheus.Metric) {
	names, err := collector.server.GetContainerNames()
	if err != nil {
		collector.logger.Printf("Can't query container names: %s", err)
		return
	}

	for _, name := range names {
		_, _, err = collector.server.GetContainerState(name)
		collector.logger.Printf(
			"Can't query container state for `%s`: %s", name, err)
	}
}
