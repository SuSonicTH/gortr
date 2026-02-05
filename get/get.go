package get

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

var sources = []Source{
	{"operator", "https://data.rtr.at/api/v2/tables/tk-agg?mediaType=csv&unpaged=true", loadOperators},
	{"geo", "https://data.rtr.at/api/v2/tables/tn-geo?de&mediaType=csv&unpaged=true", loadGeo},
	//{"nongeo", "https://data.rtr.at/api/v2/tables/tn-dienste?mediaType=csv&unpaged=true", loadNonGeo},
	// {"region", "https://data.rtr.at/api/v2/tables/tn-ortsnetze?mediaType=csv&unpaged=true"},
	// {"short", "https://data.rtr.at/api/v2/tables/tn-kurz?mediaType=csv&unpaged=true"},
	// {"param", "https://data.rtr.at/api/v2/tables/tn-skp?mediaType=csv&unpaged=true"},
}

var db_setup = []string{
	"CREATE TABLE operator(id INTEGER PRIMARY KEY, name, country, zip, city,street)",
	"CREATE TABLE number_type(id INTEGER PRIMARY KEY, name, file_name)",
	"CREATE TABLE ranges(id INTEGER PRIMARY KEY, type INTEGER, prefix, start, end, fk_operator INTEGER, range_from, range_to)",
	"CREATE TABLE singles(number PRIMARY KEY, type INTEGER, fk_range INTEGER)",
	//"CREATE TABLE region(ortsnetzkennzahl,ortsnetzname)",
	//"CREATE TABLE short(rufnummernbereich,gebiet,rufnummer,betreiber,betreiberid)",
	//"CREATE TABLE param(parameter,wertvon,wertbis,betreiber,strasse,land,plz,ort,betreiberid)",
}

var operators map[string]string = make(map[string]string)

func Numbers(db *sql.DB) error {
	databaseInit(db)

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
		ntInsert = append(ntInsert, []string{id, nrType.name, nrType.value})
	}
	return dbload.BulkInsert(db, "number_type", ntInsert)
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
	return dbload.BulkInsert(db, "operator", insert)
}

func loadGeo(db *sql.DB, rows [][]string) error {
	fmt.Printf("Inserting geo... ")
	geo := numberTypeNameToId["geo"]

	ranges := make([][]string, 0, 70000)
	singles := make([][]string, 0, 10_000_000)

	for _, row := range rows {
		if err := addRange(geo, row[0], row[2], row[3], row[5], &ranges, &singles); err != nil {
			return err
		}
	}

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
	id := strconv.Itoa(len(*ranges) + 1)

	if prefix == "7242" && from == "931000" {
		return nil
	}

	*ranges = append(*ranges, []string{id, numberType, prefix, from, to, nop, pfxFrom, pfxTo})
	if err := addSingles(prefix+pfxFrom, prefix+pfxTo, numberType, id, singles); err != nil {
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

var uniqueSingles map[int]bool = make(map[int]bool, 1000000)
var ErrDuplicateSingle = errors.New("duplicated single number")

func addSingles(pfxFrom, pfxTo, numberType string, rangeId string, singles *[][]string) error {
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
		*singles = append(*singles, []string{strconv.Itoa(i), numberType, rangeId})
	}
	return nil
}

func loadNonGeo(db *sql.DB, rows [][]string) error {
	return nil
}
