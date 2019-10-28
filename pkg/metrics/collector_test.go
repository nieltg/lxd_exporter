package metrics

import (
	"fmt"
	"log"
	"os"

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
