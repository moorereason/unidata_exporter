package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	namespace = "unidata"
)

// Flags
var (
	addr            = flag.String("listen-address", ":9777", "The address to listen on for HTTP requests.")
	metricsPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	udtbin          = flag.String("udtbin", "/opt/ud82/bin", "Path to UDTBIN.")
	timeoutDuration = flag.Duration("timeout", 2*time.Second, "Check timeout duration.")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(newUnidataCollector(*udtbin))

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html><head><title>Colleague Exporter</title></head><body><h1>Colleague Exporter</h1><p><a href="` + *metricsPath + `">Metrics</a></p></body></html>`))
		if err != nil {
			log.Errorf("error writing to HTTP response writer: %v", err)
		}
	})

	log.Infoln("Starting unidata_exporter (version=0.1.0)")
	log.Infoln("UDTBIN is", *udtbin)
	log.Infoln("Listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error in HTTP server: %s", err)
	}
}
