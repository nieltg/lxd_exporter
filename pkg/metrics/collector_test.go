package metrics

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	lxd "github.com/lxc/lxd/client"
	lxdapi "github.com/lxc/lxd/shared/api"
	mockclient "github.com/nieltg/lxd_exporter/test/mock_client"
	"github.com/nieltg/prom-example-testutil/pkg/testutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func prepare(t gomock.TestReporter) (
	controller *gomock.Controller,
	logger *log.Logger,
	server *mockclient.MockInstanceServer,
) {
	controller = gomock.NewController(t)
	logger = log.New(os.Stdout, "lxd_exporter: ", 0)
	server = mockclient.NewMockInstanceServer(controller)
	return
}

func collect(logger *log.Logger, server lxd.InstanceServer) {
	ch := make(chan prometheus.Metric)
	go func() {
		for range ch {
		}
	}()

	NewCollector(logger, server).Collect(ch)
	close(ch)
}

func Example_collector_Collect_containerNamesError() {
	controller, logger, server := prepare(nil)
	defer controller.Finish()
	server.EXPECT().GetContainerNames().Return(nil, fmt.Errorf("fail")).AnyTimes()

	collect(logger, server)
	// Output:
	// lxd_exporter: Can't query container names: fail
}

func Test_collector_Collect_queryContainerState(t *testing.T) {
	controller, logger, server := prepare(t)
	defer controller.Finish()
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0").Return(
		&lxdapi.ContainerState{}, "", nil)

	collect(logger, server)
}

func Test_collector_Collect_queryContainerStates(t *testing.T) {
	controller, logger, server := prepare(t)
	defer controller.Finish()
	server.EXPECT().GetContainerNames().Return([]string{
		"box0",
		"box1",
	}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0").Return(
		&lxdapi.ContainerState{}, "", nil)
	server.EXPECT().GetContainerState("box1").Return(
		&lxdapi.ContainerState{}, "", nil)

	collect(logger, server)
}

func Example_collector_Collect_containerStateError() {
	controller, logger, server := prepare(nil)
	defer controller.Finish()
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState(gomock.Any()).Return(
		nil, "", fmt.Errorf("fail")).AnyTimes()

	collect(logger, server)
	// Output:
	// lxd_exporter: Can't query container state for `box0`: fail
}

func Test_collector_Describe(t *testing.T) {
	controller, logger, server := prepare(t)
	defer controller.Finish()

	containsValueChannel := make(chan bool)
	ch := make(chan *prometheus.Desc)
	go func() {
		containsValue := false
		for range ch {
			if ch != nil {
				containsValue = true
			}
		}
		containsValueChannel <- containsValue
	}()
	NewCollector(logger, server).Describe(ch)
	close(ch)

	assert.True(t, <-containsValueChannel)
}

func prepareSingle(t gomock.TestReporter, state *lxdapi.ContainerState) (
	controller *gomock.Controller,
	logger *log.Logger,
	server *mockclient.MockInstanceServer,
) {
	controller, logger, server = prepare(nil)
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0").Return(state, "", nil)
	return
}

func collectAndPrint(
	logger *log.Logger,
	server lxd.InstanceServer,
	names ...string,
) {
	collector := NewCollector(logger, server)
	testutil.CollectAndPrint(collector, names...)
}

func Example_collector_cpuUsage() {
	controller, logger, server := prepareSingle(nil, &lxdapi.ContainerState{
		CPU: lxdapi.ContainerStateCPU{
			Usage: 98,
		},
	})
	defer controller.Finish()

	collectAndPrint(logger, server, "lxd_container_cpu_usage")
	// Output:
	// # HELP lxd_container_cpu_usage Container CPU Usage in Seconds
	// # TYPE lxd_container_cpu_usage gauge
	// lxd_container_cpu_usage{container_name="box0"} 98
}

func Example_collector_memUsage() {
	controller, logger, server := prepareSingle(nil, &lxdapi.ContainerState{
		Memory: lxdapi.ContainerStateMemory{
			Usage: 30,
		},
	})
	defer controller.Finish()

	collectAndPrint(logger, server, "lxd_container_mem_usage")
	// Output:
	// # HELP lxd_container_mem_usage Container Memory Usage
	// # TYPE lxd_container_mem_usage gauge
	// lxd_container_mem_usage{container_name="box0"} 30
}

func Example_collector_memUsagePeak() {
	controller, logger, server := prepareSingle(nil, &lxdapi.ContainerState{
		Memory: lxdapi.ContainerStateMemory{
			UsagePeak: 70,
		},
	})
	defer controller.Finish()

	collectAndPrint(logger, server, "lxd_container_mem_usage_peak")
	// Output:
	// # HELP lxd_container_mem_usage_peak Container Memory Usage Peak
	// # TYPE lxd_container_mem_usage_peak gauge
	// lxd_container_mem_usage_peak{container_name="box0"} 70
}

func Example_collector_swapUsage() {
	controller, logger, server := prepareSingle(nil, &lxdapi.ContainerState{
		Memory: lxdapi.ContainerStateMemory{
			SwapUsage: 10,
		},
	})
	defer controller.Finish()

	collectAndPrint(logger, server, "lxd_container_swap_usage")
	// Output:
	// # HELP lxd_container_swap_usage Container Swap Usage
	// # TYPE lxd_container_swap_usage gauge
	// lxd_container_swap_usage{container_name="box0"} 10
}

func Example_collector_swapUsagePeak() {
	controller, logger, server := prepareSingle(nil, &lxdapi.ContainerState{
		Memory: lxdapi.ContainerStateMemory{
			SwapUsagePeak: 20,
		},
	})
	defer controller.Finish()

	collectAndPrint(logger, server, "lxd_container_swap_usage_peak")
	// Output:
	// # HELP lxd_container_swap_usage_peak Container Swap Usage Peak
	// # TYPE lxd_container_swap_usage_peak gauge
	// lxd_container_swap_usage_peak{container_name="box0"} 20
}
