package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	logLevel       = flag.String("log.level", "info", "Log level for logging.")
	listenAddr     = flag.String("web.listen-address", ":8080", "The address to listen on for HTTP requests.")
	metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	interval       = flag.Int("axpert.interval", 30, "Interval in seconds for data polling.")
	metricsEnabled = flag.Bool("axpert.metrics", true, "Set to true to enable metrics collection.")
	controlEnabled = flag.Bool("axpert.control", false, "Set to true to enable control API.")
)

func main() {
	flag.Parse()

	if level, err := log.ParseLevel(*logLevel); err != nil {
		log.Fatalln(err)
	} else {
		log.SetLevel(level)
	}

	app := &Application{
		Prometheus: &Prometheus{
			Reg: createRegistry(),
		},
	}
	app.Prometheus.RegisterMetrics()

	log.Infoln("Initialising inverters connected through USB")
	invs, err := initInverters()
	if err != nil {
		log.Fatalln("failed to initialise inverters:", err)
	}
	app.Inverters = invs
	for _, inv := range app.Inverters {
		defer inv.Connector.Close()
	}

	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: app.Routes(),
	}

	if *metricsEnabled {
		go func() {
			startMetricsCollection(app, time.Duration(*interval)*time.Second)
		}()
	}

	log.Infoln("Starting axpert-gateway at:", *listenAddr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("error starting HTTP server:", err)
	}
}
