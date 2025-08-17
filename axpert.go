package main

import (
	"errors"
	"fmt"

	"github.com/marevers/energia/pkg/axpert"
	"github.com/marevers/energia/pkg/connector"
)

// Initialise any inverters connected through USB and return them
func initInverters() ([]*Inverter, error) {
	var invs []*Inverter

	crs, err := axpert.GetUSBInverters()
	if err != nil {
		return nil, errors.New("error: no Axpert inverters found")
	}

	for _, cr := range crs {
		inv := &Inverter{
			Connector: cr,
		}

		sn, err := axpert.SerialNo(cr)
		if err != nil {
			return nil, errors.New("error: failed to retrieve serial number")
		}
		inv.SerialNo = sn

		invs = append(invs, inv)
	}

	return invs, nil
}

// Takes an input and updates all relevant current settings using the contents of the input
func (i *Inverter) UpdateCurrentSettings(input any) error {
	if i.CurrentSettings == nil {
		i.CurrentSettings = &CurrentSettings{}
	}

	switch inp := input.(type) {
	case *axpert.ParallelInfo:
		if cs := mapChargeSource(inp.ACCharging); cs != "" {
			i.CurrentSettings.ChargeSource = cs
		} else {
			return fmt.Errorf("Unrecognized charge source: %s", inp.ACCharging)
		}
	case *axpert.RatingInfo:
		if osp := mapOutputSourcePriority(inp.OutputSourcePriority); osp != "" {
			i.CurrentSettings.OutputSourcePriority = osp
		} else {
			return fmt.Errorf("Unrecognized output source priority: %s", inp.OutputSourcePriority)
		}

		if csp := mapChargerSourcePriority(inp.ChargerSourcePriority); csp != "" {
			i.CurrentSettings.ChargerSourcePriority = csp
		} else {
			return fmt.Errorf("Unrecognized charger source priority: %s", inp.ChargerSourcePriority)
		}

		i.CurrentSettings.BatteryRechargeVoltage = inp.BatteryRechargeVoltage
		i.CurrentSettings.BatteryRedischargeVoltage = inp.BatteryRedischargeVoltage
		i.CurrentSettings.BatteryCutoffVoltage = inp.BatteryUnderVoltage
		i.CurrentSettings.BatteryFloatVoltage = inp.BatteryFloatVoltage
	case string:
		if dMode := mapDeviceMode(inp); dMode != "" {
			i.CurrentSettings.DeviceMode = dMode
		} else {
			return fmt.Errorf("Unrecognized device mode: %s", inp)
		}
	default:
		return fmt.Errorf("unknown input type: %T", inp)
	}

	return nil
}

// Sets the output source priority to either 'utility', 'solar' or 'sbu'
func setOutputSourcePriority(c connector.Connector, p string) error {
	var osp axpert.OutputSourcePriority

	switch p {
	case "utility":
		osp = axpert.OutputUtilityFirst
	case "solar":
		osp = axpert.OutputSolarFirst
	case "sbu":
		osp = axpert.OutputSBUFirst
	default:
		return fmt.Errorf("error: unrecognized output source priority: %s", p)
	}

	err := axpert.SetOutputSourcePriority(c, osp)
	if err != nil {
		return err
	}

	return nil
}

// Sets the charger source priority to either 'utilityfirst', 'solarfirst', 'solarandutility' or 'solar only'
func setChargerSourcePriority(c connector.Connector, p string) error {
	var csp axpert.ChargerSourcePriority

	switch p {
	case "utilityfirst":
		csp = axpert.ChargerUtilityFirst
	case "solarfirst":
		csp = axpert.ChargerSolarFirst
	case "solarandutility":
		csp = axpert.ChargerSolarAndUtility
	case "solaronly":
		csp = axpert.ChargerSolarOnly
	default:
		return fmt.Errorf("error: unrecognized charger source priority: %s", p)
	}

	err := axpert.SetChargerSourcePriority(c, csp)
	if err != nil {
		return err
	}

	return nil
}

// Sets the battery recharge voltage to a valid whole number (between 44 and 51 V).
func setBatteryRechargeVoltage(c connector.Connector, cs *CurrentSettings, v float32) error {
	switch {
	case (v < 44 || v > 51) || !isWhole(v): // Invalid value
		return fmt.Errorf("error: battery recharge voltage must be a whole number between 44 and 51 V")
	case v > cs.BatteryRedischargeVoltage: // Exceeds the redischarge voltage
		return fmt.Errorf("error: battery recharge voltage may not exceed redischarge voltage")
	case v > cs.BatteryFloatVoltage: // Exceeds the float voltage
		return fmt.Errorf("error: battery recharge voltage may not exceed float voltage")
	case v < cs.BatteryCutoffVoltage: // Lower than cutoff voltage
		return fmt.Errorf("error: battery recharge voltage may not be lower than cutoff voltage")
	}

	if err := axpert.SetBatteryRechargeVoltage(c, v); err != nil {
		return err
	}

	return nil
}

// Sets the battery redischarge voltage to a valid whole number (between 48 and 58 V).
func setBatteryRedischargeVoltage(c connector.Connector, cs *CurrentSettings, v float32) error {
	switch {
	case (v < 48 || v > 58) || !isWhole(v): // Invalid value
		return fmt.Errorf("error: battery redischarge voltage must be a whole number between 48 and 58 V")
	case v < cs.BatteryRechargeVoltage: // Lower than redischarge voltage
		return fmt.Errorf("error: battery redischarge voltage may not be lower than recharge voltage")
	case v > cs.BatteryFloatVoltage: // Exceeds the float voltage
		return fmt.Errorf("error: battery redischarge voltage may not exceed float voltage")
	case v < cs.BatteryCutoffVoltage: // Lower than cutoff voltage
		return fmt.Errorf("error: battery redischarge voltage may not be lower than cutoff voltage")
	}

	if err := axpert.SetBatteryRedischargeVoltage(c, v); err != nil {
		return err
	}

	return nil
}

// Sets the maximum AC charge current.
// TODO: Currently not working.
// func setMaxACChargeCurrent(c connector.Connector, cr uint8) error {
// 	// TODO: Add some kind of boundaries here

// 	err := axpert.SetMaxUtilityChargingCurrent(c, cr)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Returns true if f is a whole number.
func isWhole(f float32) bool {
	return f == float32(int(f))
}
