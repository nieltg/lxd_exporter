package metrics

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	lxdapi "github.com/lxc/lxd/shared/api"
	mockclient "github.com/nieltg/lxd_exporter/test/mock_client"
	"github.com/prometheus/client_golang/prometheus"
)

func drain(ch <-chan prometheus.Metric) {
	for range ch {
	}
}

func Example_collector_Collect_containerNamesError() {
	controller := gomock.NewController(nil)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return(nil, fmt.Errorf("fail")).AnyTimes()

	ch := make(chan prometheus.Metric)
	go drain(ch)
	NewCollector(logger, server).Collect(ch)
	close(ch)
	// Output:
	// lxd_exporter: Can't query container names: fail
}

func Test_collector_Collect_queryContainerState(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0").Return(
		&lxdapi.ContainerState{}, "", nil)

	ch := make(chan prometheus.Metric)
	go drain(ch)
	NewCollector(logger, server).Collect(ch)
	close(ch)
}

func Test_collector_Collect_queryContainerStates(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return([]string{
		"box0",
		"box1",
	}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0").Return(
		&lxdapi.ContainerState{}, "", nil)
	server.EXPECT().GetContainerState("box1").Return(
		&lxdapi.ContainerState{}, "", nil)

	ch := make(chan prometheus.Metric)
	go drain(ch)
	NewCollector(logger, server).Collect(ch)
	close(ch)
}

func Example_collector_Collect_containerStateError() {
	controller := gomock.NewController(nil)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState(gomock.Any()).Return(
		nil, "", fmt.Errorf("fail")).AnyTimes()

	ch := make(chan prometheus.Metric)
	go drain(ch)
	NewCollector(logger, server).Collect(ch)
	close(ch)
	// Output:
	// lxd_exporter: Can't query container state for `box0`: fail
}
