package main

import (
	"fmt"
	"time"

	"github.com/marevers/energia/pkg/axpert"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"

	log "github.com/sirupsen/logrus"
)

// mapDeviceMode converts device mode to string
func mapDeviceMode(mode string) string {
	switch mode {
	case "P": // PowerOn
		return "poweron"
	case "S": // StandBy
		return "standby"
	case "L": // Utility
		return "utility"
	case "B": // Battery
		return "battery"
	case "F": // Fault
		return "fault"
	case "H": // PowerSaving
		return "powersaving"
	default:
		return ""
	}
}

// mapChargeSource converts AC charge metric to string
func mapChargeSource(acCharge bool) string {
	if acCharge {
		return "utility"
	}

	return "solar"
}

// mapOutputSourcePriority converts axpert output source priority to string
func mapOutputSourcePriority(priority axpert.OutputSourcePriority) string {
	switch priority {
	case axpert.OutputUtilityFirst:
		return "utility"
	case axpert.OutputSolarFirst:
		return "solar"
	case axpert.OutputSBUFirst:
		return "sbu"
	default:
		return ""
	}
}

// mapChargerSourcePriority converts axpert charger source priority to string
func mapChargerSourcePriority(priority axpert.ChargerSourcePriority) string {
	switch priority {
	case axpert.ChargerUtilityFirst:
		return "utilityfirst"
	case axpert.ChargerSolarFirst:
		return "solarfirst"
	case axpert.ChargerSolarAndUtility:
		return "solarandutility"
	case axpert.ChargerSolarOnly:
		return "solaronly"
	default:
		return ""
	}
}

const (
	// LabelSerialNumber represents the inverter serial number
	LabelSerialNumber = "serialno"

	// Namespace is the metrics prefix
	Namespace = "axpert"
)

var (
	// Labels are the static labels that come with every metric
	labels = []string{
		LabelSerialNumber,
	}
)

type Prometheus struct {
	Reg     *prometheus.Registry
	Metrics struct {
		// Device status
		GridFrequencyVec *prometheus.GaugeVec
		GridVoltageVec   *prometheus.GaugeVec

		PvInputVoltage1Vec *prometheus.GaugeVec
		PvInputVoltage2Vec *prometheus.GaugeVec
		PvInputVoltage3Vec *prometheus.GaugeVec
		PvInputCurrent1Vec *prometheus.GaugeVec
		PvInputCurrent2Vec *prometheus.GaugeVec
		PvInputCurrent3Vec *prometheus.GaugeVec

		AcOutputVoltageVec       *prometheus.GaugeVec
		AcOutputFrequencyVec     *prometheus.GaugeVec
		AcOutputApparentPowerVec *prometheus.GaugeVec
		AcOutputActivePowerVec   *prometheus.GaugeVec

		OutputLoadPercentVec *prometheus.GaugeVec
		HeatSinkTempVec      *prometheus.GaugeVec

		BatVoltageVec       *prometheus.GaugeVec
		BatCapacityVec      *prometheus.GaugeVec
		BatChgCurrentVec    *prometheus.GaugeVec
		BatDischgCurrentVec *prometheus.GaugeVec

		ChargeOnVec     *prometheus.GaugeVec
		SCCChargeOn1Vec *prometheus.GaugeVec
		SCCChargeOn2Vec *prometheus.GaugeVec
		SCCChargeOn3Vec *prometheus.GaugeVec

		// Parallel device information
		LineLossVec   *prometheus.GaugeVec
		LoadOnVec     *prometheus.GaugeVec
		ACChargeOnVec *prometheus.GaugeVec

		// Rating information
		OutputSourcePrioVec       *prometheus.GaugeVec
		ChargerSourcePrioVec      *prometheus.GaugeVec
		MaxACChargerCurrentVec    *prometheus.GaugeVec
		BatteryRechgVoltageVec    *prometheus.GaugeVec
		BatteryRedischgVoltageVec *prometheus.GaugeVec
		BatteryUnderVoltageVec    *prometheus.GaugeVec
		BatteryFloatVoltageVec    *prometheus.GaugeVec

		// Statuses
		OverloadVec *prometheus.GaugeVec

		// Device mode
		DeviceModeVec *prometheus.GaugeVec

		// Output mode
		OutputModeVec *prometheus.GaugeVec

		// Scrape error
		ScrapeError prometheus.Gauge
	}
}

func createRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return reg
}

func (p *Prometheus) RegisterMetrics() {
	// Device status

	p.Metrics.GridFrequencyVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid_frequency",
		Namespace: Namespace,
		Help:      "Grid frequency in herz",
	}, labels)

	p.Metrics.GridVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid_voltage",
		Namespace: Namespace,
		Help:      "Grid voltage",
	}, labels)

	p.Metrics.PvInputVoltage1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput1_voltage",
		Namespace: Namespace,
		Help:      "PV input 1 voltage",
	}, labels)

	p.Metrics.PvInputVoltage2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput2_voltage",
		Namespace: Namespace,
		Help:      "PV input 2 voltage",
	}, labels)

	p.Metrics.PvInputVoltage3Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput3_voltage",
		Namespace: Namespace,
		Help:      "PV input 3 voltage",
	}, labels)

	p.Metrics.PvInputCurrent1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput1_current",
		Namespace: Namespace,
		Help:      "PV input 1 current in amps",
	}, labels)

	p.Metrics.PvInputCurrent2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput2_current",
		Namespace: Namespace,
		Help:      "PV input 2 current in amps",
	}, labels)

	p.Metrics.PvInputCurrent3Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "pvinput3_current",
		Namespace: Namespace,
		Help:      "PV input 3 current in amps",
	}, labels)

	p.Metrics.AcOutputVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput_voltage",
		Namespace: Namespace,
		Help:      "AC output voltage",
	}, labels)

	p.Metrics.AcOutputFrequencyVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput_frequency",
		Namespace: Namespace,
		Help:      "AC output frequency in herz",
	}, labels)

	p.Metrics.AcOutputApparentPowerVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput_apparent_power",
		Namespace: Namespace,
		Help:      "AC output apparent power in volt-amps",
	}, labels)

	p.Metrics.AcOutputActivePowerVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput_active_power",
		Namespace: Namespace,
		Help:      "AC output active power in watts",
	}, labels)

	p.Metrics.OutputLoadPercentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "output_load_percent",
		Namespace: Namespace,
		Help:      "Output load in percentage",
	}, labels)

	p.Metrics.HeatSinkTempVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "heatsink_temperature",
		Namespace: Namespace,
		Help:      "Heatsink temperature in celsius",
	}, labels)

	p.Metrics.BatVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_voltage",
		Namespace: Namespace,
		Help:      "Battery voltage",
	}, labels)

	p.Metrics.BatCapacityVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_capacity_percent",
		Namespace: Namespace,
		Help:      "Battery capacity in percentage",
	}, labels)

	p.Metrics.BatChgCurrentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_charge_current",
		Namespace: Namespace,
		Help:      "Battery charge current in amps",
	}, labels)

	p.Metrics.BatDischgCurrentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_discharge_current",
		Namespace: Namespace,
		Help:      "Battery discharge current in amps",
	}, labels)

	p.Metrics.ChargeOnVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "chargeon",
		Namespace: Namespace,
		Help:      "Returns 1 if battery  is being charged",
	}, labels)

	p.Metrics.SCCChargeOn1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sccchargeon1",
		Namespace: Namespace,
		Help:      "Returns 1 if battery is being charged with solar power 1",
	}, labels)

	p.Metrics.SCCChargeOn2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sccchargeon2",
		Namespace: Namespace,
		Help:      "Returns 1 if battery is being charged with solar power 2",
	}, labels)

	p.Metrics.SCCChargeOn3Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sccchargeon3",
		Namespace: Namespace,
		Help:      "Returns 1 if battery is being charged with solar power 3",
	}, labels)

	// Parallel device information

	p.Metrics.LineLossVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "lineloss",
		Namespace: Namespace,
		Help:      "Returns 1 if utility line is offline",
	}, labels)

	p.Metrics.LoadOnVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "loadon",
		Namespace: Namespace,
		Help:      "Returns 1 if output has load",
	}, labels)

	p.Metrics.ACChargeOnVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acchargeon",
		Namespace: Namespace,
		Help:      "Returns 1 if battery is being charged with utility power",
	}, labels)

	// Rating info

	p.Metrics.OutputSourcePrioVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "output_sourcepriority",
		Namespace: Namespace,
		Help:      "Shows the output source priority - 0: Utility first, 1: Solar first, 2: SBU first",
	}, labels)

	p.Metrics.ChargerSourcePrioVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "charger_sourcepriority",
		Namespace: Namespace,
		Help:      "Shows the charger source priority - 0: Utility first, 1: Solar first, 2: Solar and utility, 3: Solar only",
	}, labels)

	p.Metrics.MaxACChargerCurrentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "charger_maxcurrent",
		Namespace: Namespace,
		Help:      "Max AC charging current in amps",
	}, labels)

	p.Metrics.BatteryRechgVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_recharge_voltage",
		Namespace: Namespace,
		Help:      "Battery recharge voltage",
	}, labels)

	p.Metrics.BatteryRedischgVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_redischarge_voltage",
		Namespace: Namespace,
		Help:      "Battery redischarge voltage",
	}, labels)

	p.Metrics.BatteryUnderVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_cutoff_voltage",
		Namespace: Namespace,
		Help:      "Battery under / cutoff voltage",
	}, labels)

	p.Metrics.BatteryFloatVoltageVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "battery_float_voltage",
		Namespace: Namespace,
		Help:      "Battery float voltage",
	}, labels)

	// Statuses

	p.Metrics.OverloadVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "overload",
		Namespace: Namespace,
		Help:      "Returns 1 if system is overloaded",
	}, labels)

	// Device mode

	p.Metrics.DeviceModeVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "devicemode",
		Namespace: Namespace,
		Help:      "Shows the device mode - 0: PowerOnMode, 1: StandbyMode, 2: LineMode, 3: BatteryMode, 4: FaultMode, 5: PowerSavingMode",
	}, labels)

	// Output mode
	// TODO: There seem to be more output modes than the ones below

	p.Metrics.OutputModeVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "outputmode",
		Namespace: Namespace,
		Help:      "Shows the output mode - 0: SingleMachine, 1: Parallel, 2: Phase1, 3: Phase2, 4: Phase3",
	}, labels)

	// Scrape error

	p.Metrics.ScrapeError = promauto.With(p.Reg).NewGauge(prometheus.GaugeOpts{
		Name:      "scrape_error",
		Namespace: Namespace,
		Help:      "Returns 1 if the last scrape failed",
	})
}

func convertBoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}

	return 0.0
}

func (a *Application) CalculateMetrics() {
	scrapeErr := false

	for _, inv := range a.Inverters {
		log.Infof("Starting metrics retrieval from device with serialno '%s'", inv.SerialNo)

		inv.mu.Lock()
		defer inv.mu.Unlock()

		var labelValues []string

		labelValues = append(
			labelValues,
			inv.SerialNo,
		)

		// Device status
		dsp, err := axpert.DeviceGeneralStatus(inv.Connector)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve device general status from from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("device general status:")
			log.Debugf("%+v", dsp)

			a.Prometheus.Metrics.GridFrequencyVec.WithLabelValues(labelValues...).Set(float64(dsp.GridFrequency))
			a.Prometheus.Metrics.GridVoltageVec.WithLabelValues(labelValues...).Set(float64(dsp.GridVoltage))
			a.Prometheus.Metrics.PvInputVoltage1Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputVoltage1))
			a.Prometheus.Metrics.PvInputVoltage2Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputVoltage2))
			a.Prometheus.Metrics.PvInputVoltage3Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputVoltage3))
			a.Prometheus.Metrics.PvInputCurrent1Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputCurrent1))
			a.Prometheus.Metrics.PvInputCurrent2Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputCurrent2))
			a.Prometheus.Metrics.PvInputCurrent3Vec.WithLabelValues(labelValues...).Set(float64(dsp.PVInputCurrent3))
			a.Prometheus.Metrics.AcOutputVoltageVec.WithLabelValues(labelValues...).Set(float64(dsp.ACOutputVoltage))
			a.Prometheus.Metrics.AcOutputFrequencyVec.WithLabelValues(labelValues...).Set(float64(dsp.ACOutputFrequency))
			a.Prometheus.Metrics.AcOutputApparentPowerVec.WithLabelValues(labelValues...).Set(float64(dsp.ACOutputApparentPower))
			a.Prometheus.Metrics.AcOutputActivePowerVec.WithLabelValues(labelValues...).Set(float64(dsp.ACOutputActivePower))
			a.Prometheus.Metrics.OutputLoadPercentVec.WithLabelValues(labelValues...).Set(float64(dsp.OutputLoadPercent))
			a.Prometheus.Metrics.HeatSinkTempVec.WithLabelValues(labelValues...).Set(float64(dsp.HeatSinkTemperature))
			a.Prometheus.Metrics.BatVoltageVec.WithLabelValues(labelValues...).Set(float64(dsp.BatteryVoltage))
			a.Prometheus.Metrics.BatCapacityVec.WithLabelValues(labelValues...).Set(float64(dsp.BatteryCapacity))
			a.Prometheus.Metrics.BatChgCurrentVec.WithLabelValues(labelValues...).Set(float64(dsp.BatteryChargingCurrent))
			a.Prometheus.Metrics.BatDischgCurrentVec.WithLabelValues(labelValues...).Set(float64(dsp.BatteryDischargeCurrent))
			a.Prometheus.Metrics.ChargeOnVec.WithLabelValues(labelValues...).Set(convertBoolToFloat(dsp.ChargingOn))
			a.Prometheus.Metrics.SCCChargeOn1Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(dsp.SCC1ChargingOn))
			a.Prometheus.Metrics.SCCChargeOn2Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(dsp.SCC2ChargingOn))
			a.Prometheus.Metrics.SCCChargeOn3Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(dsp.SCC3ChargingOn))
		}

		// Parallel device information
		pi, err := axpert.ParallelDeviceInfo(inv.Connector, 0)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve parallel device info from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("parallel device information:")
			log.Debugf("%+v", pi)

			if err := inv.UpdateCurrentSettings(pi); err != nil {
				log.Errorf("error: failed to update current settings for device with serialno '%s': %s", err)
			}

			a.Prometheus.Metrics.LoadOnVec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pi.LoadOn))
			a.Prometheus.Metrics.LineLossVec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pi.LineLoss))
			a.Prometheus.Metrics.ACChargeOnVec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pi.ACCharging))
		}

		// Rating info
		ri, err := axpert.DeviceRatingInfo(inv.Connector)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve rating info from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("rating information:")
			log.Debugf("%+v", ri)

			a.Prometheus.Metrics.OutputSourcePrioVec.WithLabelValues(labelValues...).Set(float64(ri.OutputSourcePriority))
			a.Prometheus.Metrics.ChargerSourcePrioVec.WithLabelValues(labelValues...).Set(float64(ri.ChargerSourcePriority))
			a.Prometheus.Metrics.MaxACChargerCurrentVec.WithLabelValues(labelValues...).Set(float64(ri.MaxACChargingCurrent))
			a.Prometheus.Metrics.BatteryRechgVoltageVec.WithLabelValues(labelValues...).Set(float64(ri.BatteryRechargeVoltage))
			a.Prometheus.Metrics.BatteryRedischgVoltageVec.WithLabelValues(labelValues...).Set(float64(ri.BatteryRedischargeVoltage))
			a.Prometheus.Metrics.BatteryUnderVoltageVec.WithLabelValues(labelValues...).Set(float64(ri.BatteryUnderVoltage))
			a.Prometheus.Metrics.BatteryFloatVoltageVec.WithLabelValues(labelValues...).Set(float64(ri.BatteryFloatVoltage))

			if err := inv.UpdateCurrentSettings(ri); err != nil {
				log.Errorf("error: failed to update current settings for device with serialno '%s': %s", err)
			}
		}

		// Statuses
		warnOverload := false

		wns, err := axpert.WarningStatus(inv.Connector)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve warnings from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("wns:")
			log.Debugf("%+v", wns)
			for _, wn := range wns {
				if wn == axpert.WarnOverload {
					warnOverload = true
				}
			}

			a.Prometheus.Metrics.OverloadVec.WithLabelValues(labelValues...).Set(convertBoolToFloat(warnOverload))
		}

		// Device mode
		md, err := axpert.DeviceMode(inv.Connector)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve device mode from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("device mode:", md)

			mode, err := parseDeviceMode(md)
			if err != nil {
				scrapeErr = true
				log.Errorf("error: failed to parse device mode from device with serialno '%s': %s", inv.SerialNo, err)
			} else {
				if err := inv.UpdateCurrentSettings(md); err != nil {
					log.Errorf("error: failed to update current settings for device with serialno '%s': %s", err)
				}

				a.Prometheus.Metrics.DeviceModeVec.WithLabelValues(labelValues...).Set(mode)
			}
		}

		// Output mode
		om, err := axpert.DeviceOutputMode(inv.Connector)
		if err != nil {
			scrapeErr = true
			log.Errorf("error: failed to retrieve device output mode from device with serialno '%s'", inv.SerialNo)
		} else {
			log.Debugln("device output mode:", om)

			a.Prometheus.Metrics.OutputModeVec.WithLabelValues(labelValues...).Set(float64(om))
		}

		log.Infof("Finished metrics retrieval from device with serialno '%s'", inv.SerialNo)
	}

	if scrapeErr {
		a.Prometheus.Metrics.ScrapeError.Set(1)
		return
	}

	a.Prometheus.Metrics.ScrapeError.Set(0)
	return
}

func parseDeviceMode(m string) (float64, error) {
	switch m {
	case "P": // PowerOn
		return 0.0, nil
	case "S": // StandBy
		return 1.0, nil
	case "L": // Utility
		return 2.0, nil
	case "B": // Battery
		return 3.0, nil
	case "F": // Fault
		return 4.0, nil
	case "H": // PowerSaving
		return 5.0, nil
	default:
		return -1.0, fmt.Errorf("error: unknown device mode: %s", m)
	}
}

func startMetricsCollection(a *Application, t time.Duration) {
	tck := time.NewTicker(t)
	defer tck.Stop()

	for {
		a.CalculateMetrics()
		<-tck.C
	}
}
