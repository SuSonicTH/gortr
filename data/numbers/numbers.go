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
	{"Range", ""},
	{"Geo", ""},
	{"", "private Netze"},
	{"", "mobile Rufnummern"},
	{"", "Dial-Up Internetzugänge"},
	{"", "standortunabhängige Rufnummern"},
	{"", "konvergente Dienste"},
	{"", "tariffreie Dienste"},
	{"", "Dienste mit geregelten Tarifobergrenzen"},
	{"", "eventtarifierte Dienste"},
	{"", "SMS Dienste mit geregelten Tarifobergrenzen"},
	{"", "Routingnummern"},
	{"", "frei kalkulierbare Mehrwertdienste"},
	{"", "eventtarifierte Mehrwertdienste"},
	{"", "Dialer-Programme"},
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

func initMaps() {
	for _, numberType := range numberTypes {
		if numberType.value != "" {
			fileToNumberType[numberType.value] = &numberType
		}
		if numberType.Name != "" {
			nameToNumberType[numberType.Name] = &numberType
		}
	}
}

func Load() {
	initMaps()
	loadGeo()
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
