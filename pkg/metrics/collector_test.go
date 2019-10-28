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
