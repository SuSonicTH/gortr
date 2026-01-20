package data

import (
	"github.com/SuSonicTH/gortr/data/geo"
	"github.com/SuSonicTH/gortr/data/operator"
	"github.com/SuSonicTH/gortr/data/region"
)

var Operators map[string]*operator.Operator = nil
var Regions map[string]region.Region = nil
var Geo map[string]*geo.Geo = nil

func Read() error {
	var err error

	Operators, err = operator.Read()
	if err != nil {
		return err
	}

	Regions, err = region.Read()
	if err != nil {
		return err
	}

	Geo, err = geo.Read()
	if err != nil {
		return err
	}
	return nil
}
