package numbers

import (
	"fmt"
	"strconv"

	"github.com/SuSonicTH/gortr/data/util"
)

type NumberType struct {
	Name  string
	value string
}

var numberTypes = []NumberType{
	{"range", "ortsnetze"},
	{"geo", "geographisch"},
	{"network selection prefix", "Betreiberauswahl-Präfix"},
	{"corporate", "private Netze"},
	{"mobile", "mobile Rufnummern"},
	{"dial-up", "Dial-Up Internetzugänge"},
	{"location independent", "standortunabhängige Rufnummern"},
	{"converged service", "konvergente Dienste"},
	{"freephone", "tariffreie Dienste"},
	{"services with Ceeling", "Dienste mit geregelten Tarifobergrenzen"},
	{"event based service", "eventtarifierte Dienste"},
	{"SMS service", "SMS Dienste mit geregelten Tarifobergrenzen"},
	{"routing number", "Routingnummern"},
	{"value added service", "frei kalkulierbare Mehrwertdienste"},
	{"event based value added service", "eventtarifierte Mehrwertdienste"},
	{"dialer-program", "Dialer-Programme"},
}

var fileToNumberType map[string]*NumberType = make(map[string]*NumberType)
var nameToNumberType map[string]*NumberType = make(map[string]*NumberType)

type Number struct {
	NumberType  *NumberType
	Prefix      string
	Start       string
	End         string
	Operator_id string
	PfxFrom     string
	PfxTo       string
	Singles     []string
}

var numbers map[string]*Number = make(map[string]*Number, 0)
var minLen int = 100
var maxLen int = 0

func Load() error {
	initMaps()
	if err := loadGeo(); err != nil {
		return err
	}
	if err := loadNonGeo(); err != nil {
		return err
	}
	return nil
}

func Search(search string) (*Number, error) {
	search = util.Normalize(search)
	var length = min(len(search), maxLen)
	for i := length; i >= minLen; i-- {
		if number, ok := numbers[search[:i]]; ok {
			return number, nil
		}
	}
	return nil, fmt.Errorf("No matching number found for %s", search)
}

func initMaps() {
	for _, numberType := range numberTypes {
		fileToNumberType[numberType.value] = &numberType
		nameToNumberType[numberType.Name] = &numberType
	}
}

const (
	geoRegion int = iota
	geoRegionName
	geoFrom
	geoTo
	geoOperatorName
	geoOperatorId
)

func loadGeo() error {
	records, err := util.ReadFile("geo.csv")
	if err != nil {
		return err
	}
	geo := nameToNumberType["Geo"]
	for _, rec := range records {
		addNumber(geo, rec[geoRegion], rec[geoFrom], rec[geoTo], rec[geoOperatorId])
	}
	return nil
}

const (
	nonGeoType int = iota
	nonGeoPrefix
	nonGeoFrom
	nonGeoTo
	nonGeoOperatorName
	nonGeoOperatorId
)

func loadNonGeo() error {
	records, err := util.ReadFile("nongeo.csv")
	if err != nil {
		return err
	}
	for _, rec := range records {
		numberType, found := fileToNumberType[rec[nonGeoType]]
		if !found {
			return fmt.Errorf("Could not find number type %q", rec[nonGeoType])
		}
		addNumber(numberType, rec[nonGeoPrefix], rec[nonGeoFrom], rec[nonGeoTo], rec[nonGeoOperatorId])
	}
	return nil
}

func addNumber(numberType *NumberType, prefix, from, to, operator_id string) error {
	pfxFrom, pfxTo := getPrefix(from, to)

	numberFrom := prefix + pfxFrom
	numberTo := prefix + pfxTo

	if len(numberFrom) > maxLen {
		maxLen = len(numberFrom)
	}
	if len(numberFrom) < minLen {
		minLen = len(numberFrom)
	}

	singles, err := getSingles(numberFrom, numberTo)
	if err != nil {
		return err
	}

	num := Number{
		NumberType:  numberType,
		Prefix:      prefix,
		Start:       from,
		End:         to,
		Operator_id: operator_id,
		PfxFrom:     pfxFrom,
		PfxTo:       pfxTo,
		Singles:     singles,
	}
	for _, number := range singles {
		numbers[number] = &num
	}
	return nil
}

func getPrefix(from, to string) (pfxFrom, pfxTo string) {
	for i := len(from) - 1; i >= 1; i-- {
		if from[i] != '0' || to[i] != '9' {
			pfxFrom = from[:i+1]
			pfxTo = to[:i+1]
			return
		}
	}
	return
}

func getSingles(pfxFrom, pfxTo string) ([]string, error) {
	from, err := strconv.Atoi(pfxFrom)
	if err != nil {
		return nil, fmt.Errorf("could not convert '%s' to integer: %w", pfxFrom, err)
	}

	to, err := strconv.Atoi(pfxTo)
	if err != nil {
		return nil, fmt.Errorf("could not convert '%s' to integer: %w", pfxFrom, err)
	}

	numbers := make([]string, 0)
	for i := from; i <= to; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	return numbers, nil
}
