package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// CommandRequest represents the JSON body for control API commands
type CommandRequest struct {
	Value string `json:"value"`
}

// CommandResponse represents the JSON response for control API commands
type CommandResponse struct {
	Command string `json:"command"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	router.Handler(http.MethodGet, *metricsPath, promhttp.HandlerFor(a.Prometheus.Reg, promhttp.HandlerOpts{}))
	router.HandlerFunc(http.MethodGet, "/healthz", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) })
	router.HandlerFunc(http.MethodPost, "/api/command/:command", a.handleCommand)
	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>axpert-gateway</title></head>
		<body>
		<h1>axpert-gateway</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>
		`))
	})

	return router
}

// handleCommand handles control API commands
func (a *Application) handleCommand(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	command := params.ByName("command")

	// Check if control API is enabled
	if !*controlEnabled {
		http.Error(w, "Control API is disabled", http.StatusForbidden)
		return
	}

	// Parse JSON body
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	log.Infof("Received command: %s with value: %s", command, req.Value)

	response := CommandResponse{
		Command: command,
		Status:  "success",
		Message: "Command received and processed",
	}

	// TODO: Implement actual command processing logic here
	// For now, just return a success response

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
