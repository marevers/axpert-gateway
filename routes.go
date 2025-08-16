package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	router.Handler(http.MethodGet, *metricsPath, promhttp.HandlerFor(a.Prometheus.Reg, promhttp.HandlerOpts{}))
	router.HandlerFunc(http.MethodGet, "/healthz", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) })
	router.HandlerFunc(http.MethodPost, "/api/command/:command", a.handleCommand)
	router.HandlerFunc(http.MethodGet, "/api/inverters", a.handleListInverters)
	router.HandlerFunc(http.MethodPost, "/api/settings", a.handleGetCurrentSettings)
	router.ServeFiles("/control/*filepath", http.Dir("frontend/"))

	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>axpert-gateway</title></head>
		<body>
		<h1>axpert-gateway</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		<p><a href="/control/">Control Interface</a></p>
		</body>
		</html>
		`))
	})

	return router
}
