package operator

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
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
	OperatorName int = iota
	OperatorId
	OperatorCountry
	OperatorZip
	OperatorCity
	OperatorStreet
	OperatorService
)

var Operators map[string]*Operator = nil

func ReadOperators() error {
	operators := make(map[string]*Operator, 128)

	records, err := readFile("operator.csv")
	if err != nil {
		return err
	}

	for _, rec := range records {
		id := rec[OperatorId]
		operator, exists := operators[id]
		if !exists {
			operator = &Operator{
				id:        id,
				name:      rec[OperatorName],
				country:   rec[OperatorCountry],
				zip:       rec[OperatorZip],
				city:      rec[OperatorCity],
				street:    rec[OperatorStreet],
				servicees: make([]Service, 0),
			}
			operators[id] = operator
		}
		service := stringToService[rec[OperatorService]]
		operator.servicees = append(operator.servicees, service)
	}

	for _, op := range operators {
		slices.Sort(op.servicees)
	}

	return nil
}

func readFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", fileName, err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", fileName, err)
	}

	return lines[1:], nil
}
