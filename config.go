package main

import "github.com/marevers/energia/pkg/connector"

type Application struct {
	Config     *Config
	Prometheus *Prometheus
	Inverters  []*Inverter
}

type Inverter struct {
	Connector       *connector.USBConnector
	SerialNo        string
	CurrentSettings *CurrentSettings
}

type Config struct {
}
