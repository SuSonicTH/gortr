package region

import (
	"fmt"

	"github.com/SuSonicTH/gortr/data/util"
)

type Region struct {
	Prefix string
	Name   string
}

const (
	indexPrefix int = iota
	indexName
)

var Regions map[string]Region
var minLen int = 99
var maxLen int = 0

func Read() error {
	Regions = make(map[string]Region, 128)

	records, err := util.ReadFile("region.csv")
	if err != nil {
		return err
	}

	for _, rec := range records {
		prefix := rec[indexPrefix]
		Regions[prefix] = Region{
			Prefix: prefix,
			Name:   rec[indexName],
		}
		if len(prefix) < minLen {
			minLen = len(prefix)
		}
		if len(prefix) > maxLen {
			maxLen = len(prefix)
		}

	}

	return nil
}

func Search(number string) (*Region, error) {
	if err := Read(); err != nil {
		return nil, err
	}
	number = util.Normalize(number)
	var length = min(len(number), maxLen)
	for i := length; i >= minLen; i-- {
		if region, ok := Regions[number[:i]]; ok {
			return &region, nil
		}
	}
	return nil, fmt.Errorf("No region found for %s", number)
}
