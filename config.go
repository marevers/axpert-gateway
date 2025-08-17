package main

import (
	"sync"

	"github.com/marevers/energia/pkg/connector"
)

// Represents the application root
type Application struct {
	Prometheus *Prometheus
	Inverters  []*Inverter
}

// Represents an inverter
type Inverter struct {
	Connector       *connector.USBConnector
	SerialNo        string
	CurrentSettings *CurrentSettings
	mu              sync.Mutex
}

// Represents the current inverter settings
type CurrentSettings struct {
	OutputSourcePriority      string  `json:"outputSourcePriority"`
	ChargerSourcePriority     string  `json:"chargerSourcePriority"`
	DeviceMode                string  `json:"deviceMode"`
	ChargeSource              string  `json:"chargeSource"`
	BatteryRechargeVoltage    float32 `json:"batteryRechargeVoltage"`
	BatteryRedischargeVoltage float32 `json:"batteryRedischargeVoltage"`
	BatteryCutoffVoltage      float32 `json:"batteryCutoffVoltage"`
	BatteryFloatVoltage       float32 `json:"batteryFloatVoltage"`
}
