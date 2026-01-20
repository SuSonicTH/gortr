package operator

import (
	"slices"

	"github.com/SuSonicTH/gortr/data/util"
)

type Service int

const (
	Public Service = iota
	LandLine
	Mobile
	FixedISP
	MobileISP
	Data
)

var stringToService = map[string]Service{
	"Öffentliche Kommunikationsnetze":                                                Public,
	"Fester nummerngebundener interpersoneller Kommunikationsdienst (NB-ICS fest)":   LandLine,
	"Mobiler nummerngebundener interpersoneller Kommunikationsdienst (NB-ICS mobil)": Mobile,
	"Fester Internetzugangsdienst (IAS fest)":                                        FixedISP,
	"Datenübertragungsdienste":                                                       Data,
	"Mobiler Internetzugangsdienst (IAS mobil)":                                      MobileISP,
}

type Operator struct {
	id        string
	name      string
	country   string
	zip       string
	city      string
	street    string
	servicees []Service
}

const (
	indexName int = iota
	indexId
	indexCountry
	indexZip
	indexCity
	indexStreet
	indexService
)

func Read() (map[string]*Operator, error) {
	operators := make(map[string]*Operator, 128)

	records, err := util.ReadFile("operator.csv")
	if err != nil {
		return nil, err
	}

	for _, rec := range records {
		id := rec[indexId]
		operator, exists := operators[id]
		if !exists {
			operator = &Operator{
				id:        id,
				name:      rec[indexName],
				country:   rec[indexCountry],
				zip:       rec[indexZip],
				city:      rec[indexCity],
				street:    rec[indexStreet],
				servicees: make([]Service, 0),
			}
			operators[id] = operator
		}
		service := stringToService[rec[indexService]]
		operator.servicees = append(operator.servicees, service)
	}

	for _, op := range operators {
		slices.Sort(op.servicees)
	}

	return operators, nil
}
