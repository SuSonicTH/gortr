package geo

import (
	"fmt"
	"strconv"

	"github.com/SuSonicTH/gortr/data/util"
)

type Geo struct {
	region      string
	start       string
	end         string
	operator_id string
	pfxFrom     string
	pfxTo       string
	numbers     []string
}

const (
	indexRegion int = iota
	indexRegionName
	indexFrom
	indexTo
	indexOperatorName
	indexOperatorId
)

func Read() (map[string]*Geo, error) {
	regions := make(map[string]*Geo, 128)

	records, err := util.ReadFile("geo.csv")
	if err != nil {
		return nil, err
	}

	for _, rec := range records {
		_ = rec
		region := rec[indexRegion]
		from := rec[indexFrom]
		to := rec[indexTo]
		pfxFrom, pfxTo := getPrefix(from, to)
		numbers, err := getNumbers(region+pfxFrom, region+pfxTo)
		if err != nil {
			return nil, err
		}
		geo := Geo{
			region:      region,
			start:       from,
			end:         to,
			operator_id: rec[indexOperatorId],
			pfxFrom:     pfxFrom,
			pfxTo:       pfxTo,
			numbers:     numbers,
		}
		for _, number := range numbers {
			regions[number] = &geo
		}
	}
	return regions, nil
}

func getPrefix(from, to string) (pfxFrom, pfxTo string) {
	for i := len(from) - 1; i > 1; i-- {
		if from[i] != '0' || to[i] != '9' {
			pfxFrom = from[:i+1]
			pfxTo = to[:i+1]
			return
		}
	}
	return
}

func getNumbers(pfxFrom, pfxTo string) ([]string, error) {
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
