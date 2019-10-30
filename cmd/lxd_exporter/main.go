package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	lxd "github.com/lxc/lxd/client"
	"github.com/nieltg/lxd_exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "staging-UNVERSIONED"

	port = kingpin.Arg(
		"port", "Provide the port to listen on").Default("9472").Int16()
)

func main() {
	logger := log.New(os.Stderr, "lxd_exporter: ", log.LstdFlags)

	kingpin.Version(version)
	kingpin.Parse()

	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		logger.Fatalf("Unable to contact LXD server: %s", err)
		return
	}

	prometheus.MustRegister(metrics.NewCollector(logger, server))
	http.Handle("/metrics", promhttp.Handler())

	serveAddress := fmt.Sprintf(":%d", *port)
	log.Print("Server listening on ", serveAddress)
	log.Fatal(http.ListenAndServe(serveAddress, nil))
}
