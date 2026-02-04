package get

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

type Source struct {
	table string
	url   string
}

var sources = []Source{
	{"operator", "https://data.rtr.at/api/v2/tables/tk-agg?mediaType=csv&unpaged=true"},
	{"geo", "https://data.rtr.at/api/v2/tables/tn-geo?de&mediaType=csv&unpaged=true"},
	{"nongeo", "https://data.rtr.at/api/v2/tables/tn-dienste?mediaType=csv&unpaged=true"},
	{"region", "https://data.rtr.at/api/v2/tables/tn-ortsnetze?mediaType=csv&unpaged=true"},
	{"short", "https://data.rtr.at/api/v2/tables/tn-kurz?mediaType=csv&unpaged=true"},
	{"param", "https://data.rtr.at/api/v2/tables/tn-skp?mediaType=csv&unpaged=true"},
}

var db_setup = []string{
	//id INTEGER PRIMARY KEY, name, country, zip, city,street,
	"CREATE TABLE operator(betreiber,betreiberid,land,plz,ort,strasse,dienst,anzeigedatum,dienstaufnahme)",
	"CREATE TABLE geo(ortsnetzkennzahl,ortsnetzname,rufnummernbeginn,rufnummernende,betreiber,betreiberid)",
	"CREATE TABLE nongeo(rufnummernbereich,bereichskennzahl,rufnummernbeginn,rufnummernende,betreiber,betreiberid)",
	"CREATE TABLE region(ortsnetzkennzahl,ortsnetzname)",
	"CREATE TABLE short(rufnummernbereich,gebiet,rufnummer,betreiber,betreiberid)",
	"CREATE TABLE param(parameter,wertvon,wertbis,betreiber,strasse,land,plz,ort,betreiberid)",
}

func Numbers(db *sql.DB) error {
	for _, sql := range db_setup {
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}

	for _, source := range sources {
		rows, err := downloadFile(db, source.url, source.table)
		if err != nil {
			return fmt.Errorf("could not download file %s: %w", source.table, err)
		}
		BulkInsert(db, source.table, rows[1:])
	}
	return nil
}

func downloadFile(db *sql.DB, url string, table string) ([][]string, error) {
	fmt.Printf("Downloading %s... ", table)
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

const BATCH_SIZE int = 1000

func BulkInsert(db *sql.DB, table string, rows [][]string) error {
	fmt.Printf("Inserting into %s... ", table)
	columns := len(rows[0])
	v := make([]string, columns)
	for i := range v {
		v[i] = "?"
	}
	valueString := fmt.Sprintf("(%s),", strings.Join(v, ","))

	from := 0
	for to := BATCH_SIZE; to < len(rows); to += BATCH_SIZE {
		if err := insertBatch(db, table, rows[from:to], valueString); err != nil {
			return err
		}
		from = to
	}
	if from < len(rows) {
		if err := insertBatch(db, table, rows[from:], valueString); err != nil {
			return err
		}
	}
	fmt.Println("OK")
	return nil
}

func insertBatch(db *sql.DB, table string, rows [][]string, valueString string) error {
	values := strings.Repeat(valueString, len(rows))
	values = values[:len(values)-1]

	args := make([]any, 0, len(rows)*len(rows[0]))
	for _, row := range rows {
		for _, col := range row {
			args = append(args, col)
		}
	}
	stmt := fmt.Sprintf("INSERT INTO %s VALUES %s", table, values)
	_, err := db.Exec(stmt, args...)
	return err
}
