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

const (
	// LabelSerialNumber represents the inverter serial number
	LabelSerialNumber = "serialno"

	// LabelSource represents the charge/load source
	LabelSource = "source"

	//LabelWorkMode represents work mode
	LabelWorkMode = "mode"

	// Namespace is the metrics prefix
	Namespace = "axpert"
)

var (
	// Labels are the static labels that come with every metric
	labels = []string{
		LabelSerialNumber,
	}

	labelsSource = []string{
		LabelSerialNumber,
		LabelSource,
	}

	labelsWorkMode = []string{
		LabelSerialNumber,
		LabelWorkMode,
	}
)

type Prometheus struct {
	Reg     *prometheus.Registry
	Metrics struct {
		GridFrequency1Vec *prometheus.GaugeVec
		GridFrequency2Vec *prometheus.GaugeVec
		GridVoltage1Vec   *prometheus.GaugeVec
		GridVoltage2Vec   *prometheus.GaugeVec

		PvInputVoltage1Vec *prometheus.GaugeVec
		PvInputVoltage2Vec *prometheus.GaugeVec
		PvInputCurrent1Vec *prometheus.GaugeVec
		PvInputCurrent2Vec *prometheus.GaugeVec

		AcOutputVoltage1Vec       *prometheus.GaugeVec
		AcOutputVoltage2Vec       *prometheus.GaugeVec
		AcOutputFrequency1Vec     *prometheus.GaugeVec
		AcOutputFrequency2Vec     *prometheus.GaugeVec
		AcOutputApparentPower1Vec *prometheus.GaugeVec
		AcOutputApparentPower2Vec *prometheus.GaugeVec
		AcOutputActivePower1Vec   *prometheus.GaugeVec
		AcOutputActivePower2Vec   *prometheus.GaugeVec

		OutputLoadPercent1Vec *prometheus.GaugeVec
		OutputLoadPercent2Vec *prometheus.GaugeVec

		BatVoltageVec       *prometheus.GaugeVec
		BatCapacityVec      *prometheus.GaugeVec
		BatChgCurrentVec    *prometheus.GaugeVec
		BatDischgCurrentVec *prometheus.GaugeVec

		TotalPvInputPowerVec          *prometheus.GaugeVec
		TotalOutputLoadPercentVec     *prometheus.GaugeVec
		TotalBatChgCurrentVec         *prometheus.GaugeVec
		TotalAcOutputApparentPowerVec *prometheus.GaugeVec
		TotalAcOutputActivePowerVec   *prometheus.GaugeVec

		ChargeSourceVec *prometheus.GaugeVec
		LoadSourceVec   *prometheus.GaugeVec
		WorkModeVec     *prometheus.GaugeVec

		HasLoad1Vec     *prometheus.GaugeVec
		HasLoad2Vec     *prometheus.GaugeVec
		ACChargeOn1Vec  *prometheus.GaugeVec
		ACChargeOn2Vec  *prometheus.GaugeVec
		ChargeOnVec     *prometheus.GaugeVec
		SCCChargeOn1Vec *prometheus.GaugeVec
		SCCChargeOn2Vec *prometheus.GaugeVec
		LineLoss1Vec    *prometheus.GaugeVec
		LineLoss2Vec    *prometheus.GaugeVec
		OverloadVec     *prometheus.GaugeVec

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
	// Grid

	p.Metrics.GridFrequency1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid1_frequency",
		Namespace: Namespace,
		Help:      "Grid 1 frequency in herz",
	}, labels)

	p.Metrics.GridFrequency2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid2_frequency",
		Namespace: Namespace,
		Help:      "Grid 2 frequency in herz",
	}, labels)

	p.Metrics.GridVoltage1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid1_voltage",
		Namespace: Namespace,
		Help:      "Grid 1 voltage",
	}, labels)

	p.Metrics.GridVoltage2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "grid2_voltage",
		Namespace: Namespace,
		Help:      "Grid 2 voltage",
	}, labels)

	// PV input

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

	// AC output

	p.Metrics.AcOutputVoltage1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput1_voltage",
		Namespace: Namespace,
		Help:      "AC output 1 voltage",
	}, labels)

	p.Metrics.AcOutputVoltage2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput2_voltage",
		Namespace: Namespace,
		Help:      "AC output 2 voltage",
	}, labels)

	p.Metrics.AcOutputFrequency1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput1_frequency",
		Namespace: Namespace,
		Help:      "AC output 1 frequency in herz",
	}, labels)

	p.Metrics.AcOutputFrequency2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput2_frequency",
		Namespace: Namespace,
		Help:      "AC output 2 frequency in herz",
	}, labels)

	p.Metrics.AcOutputApparentPower1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput1_apparent_power",
		Namespace: Namespace,
		Help:      "AC output 1 apparent power in volt-amps",
	}, labels)

	p.Metrics.AcOutputApparentPower2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput2_apparent_power",
		Namespace: Namespace,
		Help:      "AC output 2 apparent power in volt-amps",
	}, labels)

	p.Metrics.AcOutputActivePower1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput1_active_power",
		Namespace: Namespace,
		Help:      "AC output 1 active power in watts",
	}, labels)

	p.Metrics.AcOutputActivePower2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acoutput2_active_power",
		Namespace: Namespace,
		Help:      "AC output 2 active power in watts",
	}, labels)

	// Output load

	p.Metrics.OutputLoadPercent1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "output1_load_percent",
		Namespace: Namespace,
		Help:      "Output 1 load in percentage",
	}, labels)

	p.Metrics.OutputLoadPercent2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "output2_load_percent",
		Namespace: Namespace,
		Help:      "Output 2 load in percentage",
	}, labels)

	// Battery

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

	// Totals

	p.Metrics.TotalPvInputPowerVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_pvinput_power",
		Namespace: Namespace,
		Help:      "Total PV input power in watts",
	}, labels)

	p.Metrics.TotalOutputLoadPercentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_output_load_percent",
		Namespace: Namespace,
		Help:      "Total output load in percentage",
	}, labels)

	p.Metrics.TotalBatChgCurrentVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_battery_charge_current",
		Namespace: Namespace,
		Help:      "Total battery charge current in amps",
	}, labels)

	p.Metrics.TotalAcOutputApparentPowerVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_acoutput_apparent_power",
		Namespace: Namespace,
		Help:      "Total AC output apparent power in volt-amps",
	}, labels)

	p.Metrics.TotalAcOutputActivePowerVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_acoutput_active_power",
		Namespace: Namespace,
		Help:      "Total AC output active power in watts",
	}, labels)

	// Charge / Load source

	p.Metrics.ChargeSourceVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "charge_source",
		Namespace: Namespace,
		Help:      "Charge source",
	}, labelsSource)

	p.Metrics.LoadSourceVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "load_source",
		Namespace: Namespace,
		Help:      "Load source",
	}, labelsSource)

	// Work mode

	p.Metrics.WorkModeVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "work_mode",
		Namespace: Namespace,
		Help:      "Work mode",
	}, labelsWorkMode)

	// Boolean statuses

	p.Metrics.HasLoad1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hasload1",
		Namespace: Namespace,
		Help:      "Returns 1 if output 1 has load",
	}, labels)

	p.Metrics.HasLoad2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "hasload2",
		Namespace: Namespace,
		Help:      "Returns 1 if output 2 has load",
	}, labels)

	p.Metrics.ACChargeOn1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acchargeon1",
		Namespace: Namespace,
		Help:      "Returns 1 if line 1 is being charged with utility power",
	}, labels)

	p.Metrics.ACChargeOn2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "acchargeon2",
		Namespace: Namespace,
		Help:      "Returns 1 if line 2 is being charged with utility power",
	}, labels)

	p.Metrics.SCCChargeOn1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sccchargeon1",
		Namespace: Namespace,
		Help:      "Returns 1 if line 1 is being charged with solar power",
	}, labels)

	p.Metrics.SCCChargeOn2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "sccchargeon2",
		Namespace: Namespace,
		Help:      "Returns 1 if line 2 is being charged with solar power",
	}, labels)

	p.Metrics.LineLoss1Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "lineloss1",
		Namespace: Namespace,
		Help:      "Returns 1 if utility line 1 is offline",
	}, labels)

	p.Metrics.LineLoss2Vec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "lineloss2",
		Namespace: Namespace,
		Help:      "Returns 1 if utility line 2 is offline",
	}, labels)

	p.Metrics.OverloadVec = promauto.With(p.Reg).NewGaugeVec(prometheus.GaugeOpts{
		Name:      "overload",
		Namespace: Namespace,
		Help:      "Returns 1 if system is overloaded",
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

func (a *Application) CalculateMetrics() error {
	a.Prometheus.Metrics.ScrapeError.Set(0)

	for _, inv := range a.Inverters {
		pdi1, err := axpert.ParallelDeviceInfo(inv.Connector, 0)
		if err != nil {
			return fmt.Errorf("error: failed to retrieve parallel device info from device 1 with serialno '%s'", inv.SerialNo)
		}
		pdi2, err := axpert.ParallelDeviceInfo(inv.Connector, 1)
		if err != nil {
			return fmt.Errorf("error: failed to retrieve parallel device info from device 1 with serialno '%s'", inv.SerialNo)
		}
		log.Infof("Retrieved metrics from device with serialno '%s'", inv.SerialNo)

		var labelValues []string

		labelValues = append(
			labelValues,
			inv.SerialNo,
		)

		// err := p.Session.GetWorkInfo()
		// if err != nil {
		// 	a.Prometheus.Metrics.ScrapeError.Set(1)
		// 	return err
		// }

		// Standard metrics

		a.Prometheus.Metrics.GridFrequency1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.GridFrequency))
		a.Prometheus.Metrics.GridFrequency2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.GridFrequency))
		a.Prometheus.Metrics.GridVoltage1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.GridVoltage))
		a.Prometheus.Metrics.GridVoltage2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.GridVoltage))

		a.Prometheus.Metrics.PvInputVoltage1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.PV1InputVoltage))
		a.Prometheus.Metrics.PvInputVoltage2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.PV1InputVoltage))
		a.Prometheus.Metrics.PvInputCurrent1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.PV1InputCurrent))
		a.Prometheus.Metrics.PvInputCurrent2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.PV1InputCurrent))

		a.Prometheus.Metrics.AcOutputVoltage1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.ACOutputVoltage))
		a.Prometheus.Metrics.AcOutputVoltage2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.ACOutputVoltage))
		a.Prometheus.Metrics.AcOutputFrequency1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.ACOutputFrequency))
		a.Prometheus.Metrics.AcOutputFrequency2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.ACOutputFrequency))
		a.Prometheus.Metrics.AcOutputApparentPower1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.ACOutputApparentPower))
		a.Prometheus.Metrics.AcOutputApparentPower2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.ACOutputApparentPower))
		a.Prometheus.Metrics.AcOutputActivePower1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.ACOutputActivePower))
		a.Prometheus.Metrics.AcOutputActivePower2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.ACOutputActivePower))

		a.Prometheus.Metrics.OutputLoadPercent1Vec.WithLabelValues(labelValues...).Set(float64(pdi1.OutputLoadPercent))
		a.Prometheus.Metrics.OutputLoadPercent2Vec.WithLabelValues(labelValues...).Set(float64(pdi2.OutputLoadPercent))

		a.Prometheus.Metrics.BatVoltageVec.WithLabelValues(labelValues...).Set(float64(pdi1.BatteryVoltage))
		a.Prometheus.Metrics.BatCapacityVec.WithLabelValues(labelValues...).Set(float64(pdi1.BatteryCapacity))
		a.Prometheus.Metrics.BatChgCurrentVec.WithLabelValues(labelValues...).Set(float64(pdi1.BatteryChargingCurrent))
		a.Prometheus.Metrics.BatDischgCurrentVec.WithLabelValues(labelValues...).Set(float64(pdi1.BatteryDischargeCurrent))

		// a.Prometheus.Metrics.TotalPvInputPowerVec.WithLabelValues(labelValues...).Set(p.Session.WorkInfo.TotalPvInputPower)
		// a.Prometheus.Metrics.TotalOutputLoadPercentVec.WithLabelValues(labelValues...).Set(p.Session.WorkInfo.TotalOutputLoadPercent)
		// a.Prometheus.Metrics.TotalBatChgCurrentVec.WithLabelValues(labelValues...).Set(p.Session.WorkInfo.TotalBatChgCurrent)
		a.Prometheus.Metrics.TotalAcOutputApparentPowerVec.WithLabelValues(labelValues...).Set(float64(pdi1.TotalACOutputApparentPower))
		a.Prometheus.Metrics.TotalAcOutputActivePowerVec.WithLabelValues(labelValues...).Set(float64(pdi1.TotalOutputActivePower))

		// Boolean statuses

		a.Prometheus.Metrics.HasLoad1Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi1.LoadOn))
		a.Prometheus.Metrics.HasLoad2Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi2.LoadOn))
		a.Prometheus.Metrics.ACChargeOn1Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi1.ACCharging))
		a.Prometheus.Metrics.ACChargeOn2Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi2.ACCharging))
		a.Prometheus.Metrics.SCCChargeOn1Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi1.SCC1Charging))
		a.Prometheus.Metrics.SCCChargeOn2Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi2.SCC1Charging))
		a.Prometheus.Metrics.LineLoss1Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi1.LineLoss))
		a.Prometheus.Metrics.LineLoss2Vec.WithLabelValues(labelValues...).Set(convertBoolToFloat(pdi2.LineLoss))
		// a.Prometheus.Metrics.OverloadVec.WithLabelValues(labelValues...).Set(convertBoolToFloat())

		// Named statuses

		// var labelValuesChargeSource []string = append(
		// 	labelValues,
		// 	p.Session.WorkInfo.ChargeSource,
		// )

		// a.Prometheus.Metrics.ChargeSourceVec.Reset()
		// a.Prometheus.Metrics.ChargeSourceVec.WithLabelValues(labelValuesChargeSource...).Set(1)

		// var labelValuesLoadSource []string = append(
		// 	labelValues,
		// 	p.Session.WorkInfo.LoadSource,
		// )

		// a.Prometheus.Metrics.LoadSourceVec.Reset()
		// a.Prometheus.Metrics.LoadSourceVec.WithLabelValues(labelValuesLoadSource...).Set(1)

		// var labelValuesWorkMode []string = append(
		// 	labelValues,
		// 	p.Session.WorkInfo.WorkMode,
		// )

		// a.Prometheus.Metrics.WorkModeVec.Reset()
		// a.Prometheus.Metrics.WorkModeVec.WithLabelValues(labelValuesWorkMode...).Set(1)
	}

	return nil
}

func startMetricsCollection(a *Application, t time.Duration) {
	tck := time.NewTicker(t)
	defer tck.Stop()

	for {
		err := a.CalculateMetrics()
		if err != nil {
			log.Warnln(err)
		}

		<-tck.C
	}
}
