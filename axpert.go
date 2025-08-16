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

// Represents the current inverter settings
type CurrentSettings struct {
	OutputSourcePriority  string `json:"outputSourcePriority"`
	ChargerSourcePriority string `json:"chargerSourcePriority"`
}

// Returns the current inverter settings
func getCurrentSettings(c connector.Connector) (CurrentSettings, error) {
	ri, err := axpert.DeviceRatingInfo(c)
	if err != nil {
		return CurrentSettings{}, err
	}

	var osp string
	switch ri.OutputSourcePriority {
	case axpert.OutputUtilityFirst:
		osp = "utility"
	case axpert.OutputSolarFirst:
		osp = "solar"
	case axpert.OutputSBUFirst:
		osp = "sbu"
	default:
		return CurrentSettings{}, fmt.Errorf("error: unrecognized output source priority: %s", ri.OutputSourcePriority)
	}

	var csp string
	switch ri.ChargerSourcePriority {
	case axpert.ChargerUtilityFirst:
		csp = "utilityfirst"
	case axpert.ChargerSolarFirst:
		csp = "solarfirst"
	case axpert.ChargerSolarAndUtility:
		csp = "solarandutility"
	case axpert.ChargerSolarOnly:
		csp = "solaronly"
	default:
		return CurrentSettings{}, fmt.Errorf("error: unrecognized charger source priority: %s", ri.ChargerSourcePriority)
	}

	return CurrentSettings{
		OutputSourcePriority:  osp,
		ChargerSourcePriority: csp,
	}, nil

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
