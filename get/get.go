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
	{"geo", "https://data.rtr.at/api/v2/tables/tn-geo?de&mediaType=csv&unpaged=true"},
	{"region", "https://data.rtr.at/api/v2/tables/tn-ortsnetze?mediaType=csv&unpaged=true"},
	{"nongeo", "https://data.rtr.at/api/v2/tables/tn-dienste?mediaType=csv&unpaged=true"},
	{"short", "https://data.rtr.at/api/v2/tables/tn-kurz?mediaType=csv&unpaged=true"},
	{"param", "https://data.rtr.at/api/v2/tables/tn-skp?mediaType=csv&unpaged=true"},
	{"operator", "https://data.rtr.at/api/v2/tables/tk-agg?mediaType=csv&unpaged=true"},
}

func Numbers(db *sql.DB) error {
	for _, source := range sources {
		err := downloadFile(db, source.url, source.table)
		if err != nil {
			return fmt.Errorf("could not download file %s: %w", source.table, err)
		}
	}
	return nil
}

func downloadFile(db *sql.DB, url string, table string) error {
	fmt.Printf("Downloading %s... ", table)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return err
	}
	db.Exec("create table " + table + "(" + strings.Join(rows[0], ",") + ")")
	BulkInsert(db, table, rows[1:])
	fmt.Println("OK")
	return nil
}

const BATCH_SIZE int = 1000

func BulkInsert(db *sql.DB, table string, rows [][]string) error {
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
