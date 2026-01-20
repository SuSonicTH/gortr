package region

import (
	"github.com/SuSonicTH/gortr/data/util"
)

type Region struct {
	prefix string
	name   string
}

const (
	indexPrefix int = iota
	indexName
)

func Read() (map[string]Region, error) {
	regions := make(map[string]Region, 128)

	records, err := util.ReadFile("region.csv")
	if err != nil {
		return nil, err
	}

	for _, rec := range records {
		prefix := rec[indexPrefix]
		regions[prefix] = Region{
			prefix: prefix,
			name:   rec[indexName],
		}
	}

	return regions, nil
}
