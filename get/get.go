package get

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Source struct {
	fileName string
	url      string
}

var sources = []Source{
	{"geo.csv", "https://data.rtr.at/api/v2/tables/tn-geo?de&mediaType=csv&unpaged=true"},
	{"region.csv", "https://data.rtr.at/api/v2/tables/tn-ortsnetze?mediaType=csv&unpaged=true"},
	{"nongeo.csv", "https://data.rtr.at/api/v2/tables/tn-dienste?mediaType=csv&unpaged=true"},
	{"short.csv", "https://data.rtr.at/api/v2/tables/tn-kurz?mediaType=csv&unpaged=true"},
	{"param.csv", "https://data.rtr.at/api/v2/tables/tn-skp?mediaType=csv&unpaged=true"},
	{"operator.csv", "https://data.rtr.at/api/v2/tables/tk-agg?mediaType=csv&unpaged=true"},
}

func Numbers() error {
	for _, source := range sources {
		err := downloadFile(source.url, source.fileName)
		if err != nil {
			return fmt.Errorf("could not download file %s: %w", source.fileName, err)
		}
	}
	return nil
}

func downloadFile(url string, filepath string) error {
	fmt.Printf("Downloading %s... ", filepath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("OK")
	return nil
}
