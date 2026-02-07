package dbload

import (
	"database/sql"
	"fmt"
	"strings"
)

const maxParameters = 32766

func BulkInsert(db *sql.DB, table string, rows [][]string) error {
	columns := len(rows[0])
	v := make([]string, columns)
	for i := range v {
		v[i] = "?"
	}
	valueString := fmt.Sprintf("(%s),", strings.Join(v, ","))

	from := 0
	batchSize := maxParameters / columns
	for to := batchSize; to < len(rows); to += batchSize {
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
