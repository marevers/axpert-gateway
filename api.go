package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Represents the JSON body for control API commands
type CommandRequest struct {
	Value    string `json:"value"`
	SerialNo string `json:"serialno"`
}

// Represents the JSON body for settings requests
type SettingsRequest struct {
	SerialNo string `json:"serialno"`
}

// Represents the JSON body for settings responses
type SettingsResponse struct {
	SerialNo string          `json:"serialno"`
	Settings CurrentSettings `json:"settings"`
}

// Represents the JSON response for control API commands
type CommandResponse struct {
	Command string `json:"command"`
	Value   string `json:"value"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Represents an inverter for the API
type InverterInfo struct {
	SerialNo string `json:"serialno"`
}

// Represents the response for listing inverters
type InvertersResponse struct {
	Inverters []InverterInfo `json:"inverters"`
	Count     int            `json:"count"`
}

// Defines the signature for command handler functions
type CommandHandler func(app *Application, req CommandRequest) error

// Maps command names to their handler functions
var commandHandlers = map[string]CommandHandler{
	"setOutputPriority":  handleSetOutputPriority,
	"setChargerPriority": handleSetChargerPriority,
	// "setMaxChargeCurrent": handleSetMaxChargeCurrent,
}

// Handles control API commands
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

	log.Infof("Received command: %s with value: %s for serialno: %s", command, req.Value, req.SerialNo)

	// Look up command handler
	handler, exists := commandHandlers[command]
	if !exists {
		log.Errorf("Unknown command: %s", command)
		http.Error(w, "Unknown command", http.StatusBadRequest)
		return
	}

	// Execute command
	if err := handler(a, req); err != nil {
		log.Errorf("Command execution failed: %v", err)
		response := CommandResponse{
			Command: command,
			Value:   req.Value,
			Status:  "error",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CommandResponse{
		Command: command,
		Value:   req.Value,
		Status:  "success",
		Message: "Command executed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Finds an inverter by its serial number
func findInverterBySerial(app *Application, serialNo string) (*Inverter, error) {
	for _, inv := range app.Inverters {
		if inv.SerialNo == serialNo {
			return inv, nil
		}
	}
	return nil, fmt.Errorf("inverter with serial number %s not found", serialNo)
}

// Handles listing all available inverters
func (a *Application) handleListInverters(w http.ResponseWriter, r *http.Request) {
	inverters := make([]InverterInfo, 0, len(a.Inverters))

	for _, inv := range a.Inverters {
		inverters = append(inverters, InverterInfo{
			SerialNo: inv.SerialNo,
		})
	}

	response := InvertersResponse{
		Inverters: inverters,
		Count:     len(inverters),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Failed to encode inverters response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handles retrieving current settings for an inverter
func (a *Application) handleGetCurrentSettings(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body
	var req SettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	log.Infof("Retrieving current settings for inverter with serialno '%s'", req.SerialNo)

	inv, err := findInverterBySerial(a, req.SerialNo)
	if err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Use cached current settings from the inverter struct
	if inv.CurrentSettings == nil {
		log.Errorf("Current settings not available for %s (may not have been collected yet)", req.SerialNo)
		http.Error(w, "Current settings not available - please wait for next metrics collection cycle", http.StatusServiceUnavailable)
		return
	}

	settings := *inv.CurrentSettings
	response := SettingsResponse{
		SerialNo: inv.SerialNo,
		Settings: settings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Failed to encode settings response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Sets the output source priority for a specific inverter
func handleSetOutputPriority(app *Application, req CommandRequest) error {
	log.Infof("Setting output priority to: %s for inverter: %s", req.Value, req.SerialNo)

	inv, err := findInverterBySerial(app, req.SerialNo)
	if err != nil {
		return err
	}

	return setOutputSourcePriority(inv.Connector, req.Value)
}

// Sets the charger source priority for a specific inverter
func handleSetChargerPriority(app *Application, req CommandRequest) error {
	log.Infof("Setting charger priority to: %s for inverter: %s", req.Value, req.SerialNo)

	inv, err := findInverterBySerial(app, req.SerialNo)
	if err != nil {
		return err
	}

	return setChargerSourcePriority(inv.Connector, req.Value)
}

// Sets the maximum AC charge current for a specific inverter
// func handleSetMaxChargeCurrent(app *Application, req CommandRequest) error {
// 	log.Infof("Setting max charge current to: %s for inverter: %s", req.Value, req.SerialNo)

// 	// Convert string value to uint8
// 	current, err := strconv.ParseUint(req.Value, 10, 8)
// 	if err != nil {
// 		return fmt.Errorf("invalid current value: %s", req.Value)
// 	}

// 	inv, err := findInverterBySerial(app, req.SerialNo)
// 	if err != nil {
// 		return err
// 	}

// 	return setMaxACChargeCurrent(inv.Connector, uint8(current))
// }
