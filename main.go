package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/SuSonicTH/gortr/get"
	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE = "./gortr.sqlite3"

func main() {
	showHelp := true
	pRefresh := flag.Bool("refresh", false, "get data from rtr.at")
	pSearch := flag.String("search", "", "serach for a matching number")

	flag.Parse()

	if *pRefresh {
		showHelp = false
		refresh()
	}

	if *pSearch != "" {
		showHelp = false
		db := openDb()
		defer db.Close()

		searchNumber(db, *pSearch)
	}

	if showHelp {
		fmt.Printf("no argument given\n")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
}

func openDb() *sql.DB {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		refresh()
	}

	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		panic(err)
	}

	return db
}

func refresh() {
	os.Remove(DB_FILE)

	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec("PRAGMA synchronous = OFF")
	db.Exec("PRAGMA journal_mode = OFF")
	db.Exec("PRAGMA locking_mode = EXCLUSIVE")

	if err := get.FromRtr(db); err != nil {
		panic(err)
	}
}

func searchNumber(db *sql.DB, search string) {
	number := Normalize(search)

	for i := len(number); i > 0; i-- {
		if single := getSingle(db, number[:i]); single != nil {
			println(*single)
		}
	}

}

func Normalize(number string) string {
	num := strings.Trim(number, " \t\r\n")
	num = strings.TrimPrefix(num, "0043")
	num = strings.TrimPrefix(num, "+43")
	num = strings.TrimPrefix(num, "0")
	return num
}

func getSingle(db *sql.DB, number string) *int {
	rows, err := db.Query("select fk_range from singles where number = ?", number)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var rangeId int
		err = rows.Scan(&rangeId)
		if err != nil {
			panic(err)
		}
		return &rangeId
	}
	return nil
}
