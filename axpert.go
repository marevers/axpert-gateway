package main

import (
	"errors"

	"github.com/marevers/energia/pkg/axpert"
)

// Initialise any inverters connected through USB and return them.
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
