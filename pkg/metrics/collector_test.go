package metrics

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mockclient "github.com/nieltg/lxd_exporter/test/mock_client"
	"github.com/prometheus/client_golang/prometheus"
)

func Example_collector_Collect_containerNamesError() {
	controller := gomock.NewController(nil)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return(nil, fmt.Errorf("fail")).AnyTimes()

	NewCollector(logger, server).Collect(make(chan prometheus.Metric))
	// Output:
	// lxd_exporter: Can't query container names: fail
}

func Test_collector_Collect_queryContainerState(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	logger := log.New(os.Stdout, "lxd_exporter: ", 0)
	server := mockclient.NewMockInstanceServer(controller)
	server.EXPECT().GetContainerNames().Return([]string{"box0"}, nil).AnyTimes()
	server.EXPECT().GetContainerState("box0")

	NewCollector(logger, server).Collect(make(chan prometheus.Metric))
}
