package get

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/SuSonicTH/gortr/dbload"
)

type Processing func(db *sql.DB, rows [][]string) error

type Source struct {
	name       string
	url        string
	processing Processing
}

type NumberType struct {
	name  string
	value string
}

var numberTypes = []NumberType{
	{"local area", "ortsnetze"},
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
var numberTypeNameToId map[string]string = make(map[string]string, len(numberTypes))
var numberTypeFileToId map[string]string = make(map[string]string, len(numberTypes))

var sources = []Source{
	{"operator", "https://data.rtr.at/api/v2/tables/tk-agg?mediaType=csv&unpaged=true", loadOperators},
	{"geo", "https://data.rtr.at/api/v2/tables/tn-geo?de&mediaType=csv&unpaged=true", loadGeo},
	{"nongeo", "https://data.rtr.at/api/v2/tables/tn-dienste?mediaType=csv&unpaged=true", loadNonGeo},
	// {"region", "https://data.rtr.at/api/v2/tables/tn-ortsnetze?mediaType=csv&unpaged=true"},
	// {"short", "https://data.rtr.at/api/v2/tables/tn-kurz?mediaType=csv&unpaged=true"},
	// {"param", "https://data.rtr.at/api/v2/tables/tn-skp?mediaType=csv&unpaged=true"},
}

var db_setup = []string{
	"CREATE TABLE operators(id INTEGER PRIMARY KEY, name, country, zip, city,street)",
	"CREATE TABLE number_type(id INTEGER PRIMARY KEY, name, file_name)",
	"CREATE TABLE ranges(id INTEGER PRIMARY KEY, type INTEGER, prefix, start, end, fk_operator INTEGER)",
	"CREATE TABLE singles(number PRIMARY KEY, fk_range INTEGER)",
	//"CREATE TABLE region(ortsnetzkennzahl,ortsnetzname)",
	//"CREATE TABLE short(rufnummernbereich,gebiet,rufnummer,betreiber,betreiberid)",
	//"CREATE TABLE param(parameter,wertvon,wertbis,betreiber,strasse,land,plz,ort,betreiberid)",
}

var operators map[string]string = make(map[string]string)

var ignoredRanges map[string]bool = make(map[string]bool)

var lastRangesId = 1

func FromRtr(db *sql.DB) error {
	databaseInit(db)
	initIgnoredRanges()

	for _, source := range sources {
		rows, err := downloadFile(source.url, source.name)
		if err != nil {
			return fmt.Errorf("could not download file %s: %w", source.name, err)
		}
		if err := source.processing(db, rows); err != nil {
			return err
		}
	}
	return nil
}

func databaseInit(db *sql.DB) error {
	for _, sql := range db_setup {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}

	ntInsert := make([][]string, 0, len(numberTypes))
	for i, nrType := range numberTypes {
		id := strconv.Itoa(i + 1)
		numberTypeNameToId[nrType.name] = id
		numberTypeFileToId[nrType.value] = id
		ntInsert = append(ntInsert, []string{id, nrType.name, nrType.value})
	}
	return dbload.BulkInsert(db, "number_type", ntInsert)
}

func initIgnoredRanges() {
	if data, err := os.ReadFile("ignoredRanges"); err == nil {
		for _, number := range strings.Split(string(data), "\n") {
			ignoredRanges[number] = true
		}
	}
}

func downloadFile(url string, name string) ([][]string, error) {
	fmt.Printf("Downloading %s... ", name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return nil, err
	}
	fmt.Println("OK")
	return rows[1:], nil
}

func loadOperators(db *sql.DB, rows [][]string) error {
	insert := make([][]string, 0)
	for _, operator := range rows {
		name := operator[0]
		if _, exists := operators[name]; !exists {
			id := operator[1]
			country := operator[2]
			zip := operator[3]
			city := operator[4]
			street := operator[5]
			insert = append(insert, []string{id, name, country, zip, city, street})
			operators[name] = id
		}
	}
	return dbload.BulkInsert(db, "operators", insert)
}

func loadGeo(db *sql.DB, rows [][]string) error {
	fmt.Printf("Reading geo... ")
	geo := numberTypeNameToId["geo"]

	ranges := make([][]string, 0, 70_000)
	singles := make([][]string, 0, 800_000)

	for _, row := range rows {
		if err := addRange(geo, row[0], row[2], row[3], row[5], &ranges, &singles); err != nil {
			return err
		}
	}
	fmt.Printf("OK read %d ranges and %d singles\n", len(ranges), len(singles))

	fmt.Printf("Inserting geo... ")
	if err := dbload.BulkInsert(db, "ranges", ranges); err != nil {
		return err
	}

	if err := dbload.BulkInsert(db, "singles", singles); err != nil {
		return err
	}

	lastRangesId = len(ranges) + 1
	fmt.Println("OK")
	return nil
}

func loadNonGeo(db *sql.DB, rows [][]string) error {
	fmt.Printf("Reading nongeo... ")

	ranges := make([][]string, 0, 20_000)
	singles := make([][]string, 0, 15_000_000)

	for _, row := range rows {
		nrType := numberTypeFileToId[row[0]]
		if err := addRange(nrType, row[1], row[2], row[3], row[5], &ranges, &singles); err != nil {
			return err
		}
	}
	fmt.Printf("OK read %d ranges and %d singles\n", len(ranges), len(singles))

	fmt.Printf("Inserting nongeo... ")
	if err := dbload.BulkInsert(db, "ranges", ranges); err != nil {
		return err
	}

	if err := dbload.BulkInsert(db, "singles", singles); err != nil {
		return err
	}

	fmt.Println("OK")
	return nil
}

func addRange(numberType, prefix, from, to, nop string, ranges *[][]string, singles *[][]string) error {
	pfxFrom, pfxTo := getPrefix(from, to)
	id := strconv.Itoa(len(*ranges) + lastRangesId)

	if _, found := ignoredRanges[prefix+from]; found {
		return nil
	}

	*ranges = append(*ranges, []string{id, numberType, prefix, from, to, nop})
	if err := addSingles(prefix+pfxFrom, prefix+pfxTo, id, singles); err != nil {
		if err == ErrDuplicateSingle {
			fmt.Fprintf(os.Stderr, "duplicated singles found for prefix %s from %s to %s\n", prefix, from, to)
		} else {
			return err
		}
	}
	return nil
}

func getPrefix(from, to string) (pfxFrom, pfxTo string) {
	for i := len(from) - 1; i >= 0; i-- {
		if from[i] != '0' || to[i] != '9' {
			pfxFrom = from[:i+1]
			pfxTo = to[:i+1]
			return
		}
	}
	return
}

var uniqueSingles map[int]bool = make(map[int]bool, 16_000_000)
var ErrDuplicateSingle = errors.New("duplicated single number")

func addSingles(pfxFrom, pfxTo, rangeId string, singles *[][]string) error {
	from, err := strconv.Atoi(pfxFrom)
	if err != nil {
		return fmt.Errorf("could not convert '%s' to integer: %w", pfxFrom, err)
	}
	to, err := strconv.Atoi(pfxTo)
	if err != nil {
		return fmt.Errorf("could not convert '%s' to integer: %w", pfxFrom, err)
	}

	for i := from; i <= to; i++ {
		if _, exists := uniqueSingles[i]; exists {
			return ErrDuplicateSingle
		} else {
			uniqueSingles[i] = true
		}
		*singles = append(*singles, []string{strconv.Itoa(i), rangeId})
	}
	return nil
}
