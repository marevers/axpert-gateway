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

// Updates a single setting in the settings cache
func (i *Inverter) UpdateCurrentSettings(settingName, value string) error {
	if i.CurrentSettings == nil {
		i.CurrentSettings = &CurrentSettings{}
	}

	switch settingName {
	case "outputSourcePriority":
		i.CurrentSettings.OutputSourcePriority = value
	case "chargerSourcePriority":
		i.CurrentSettings.ChargerSourcePriority = value
	case "deviceMode":
		i.CurrentSettings.DeviceMode = value
	case "chargeSource":
		i.CurrentSettings.ChargeSource = value
	default:
		return fmt.Errorf("unknown setting name: %s", settingName)
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
