package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/SuSonicTH/gortr/get"
	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE = "./gortr.sqlite3"

func main() {
	pRefresh := flag.Bool("refresh", false, "get data from rtr.at")
	pSearch := flag.String("search", "", "serach for a matching number")
	pRegion := flag.String("region", "", "serach for a matching region")
	pListRegions := flag.Bool("listRegions", false, "list all regions with name")

	flag.Parse()

	if *pRefresh {
		refresh()
	}

	db := openDb()
	defer db.Close()

	if *pSearch != "" {
		searchNumber(db, *pSearch)
		return
	}

	if *pRegion != "" {
		searchRegion(db, *pRegion)
		return
	}

	if *pListRegions {
		listRegions(db)
		return
	}

	if !*pRefresh {
		fmt.Printf("no argument given\n")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
}

func openDb() *sql.DB {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		fmt.Println("No database found, loading data from RTR")
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
		if rangeId, single, err := getSingle(db, number[:i]); err == nil {
			printNumber(db, search, rangeId, single)
			return
		} else if err != ErrorNotFound {
			panic(err)
		}
	}
	fmt.Printf("%s not found\n", search)
}

func Normalize(number string) string {
	num := strings.Trim(number, " \t\r\n")
	num = strings.TrimPrefix(num, "0043")
	num = strings.TrimPrefix(num, "+43")
	num = strings.TrimPrefix(num, "0")
	return num
}

var ErrorNotFound = errors.New("Number not found")

func getSingle(db *sql.DB, search string) (rangeId string, single string, retErr error) {
	rows, err := db.Query("select fk_range, number from singles where number = ?", search)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		retErr = rows.Scan(&rangeId, &single)
	} else {
		retErr = ErrorNotFound
	}
	return
}

const numberFormatStringRange = `searched      %s

number        %s
prefix        %s
range start   %s
range end     %s

operator      %s
              %s
              %s %s %s
              
`

func printNumber(db *sql.DB, search string, rangeId string, single string) (retErr error) {
	rows, err := db.Query(`
		select r.prefix, r.start, r.end, 
		       o.name, o.street, o.country, o.zip, o.city
		from ranges r,
		     operators o
		where r.id = ?
		  and r.fk_operator = o.id`, rangeId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var prefix, start, end, operatorName, street, country, zip, city string
		retErr = rows.Scan(&prefix, &start, &end, &operatorName, &street, &zip, &city, &country)
		fmt.Printf(numberFormatStringRange, search, single, prefix, start, end, operatorName, street, country, zip, city)
	}
	return
}

func searchRegion(db *sql.DB, search string) {
	number := Normalize(search)

	for i := len(number); i > 0; i-- {
		if name, err := getRegion(db, number[:i]); err == nil {
			fmt.Printf("%s,%s\n", number[:i], name)
			return
		} else if err != ErrorNotFound {
			panic(err)
		}
	}
	fmt.Printf("%s not found\n", search)
}

func listRegions(db *sql.DB) {
	rows, err := db.Query("select prefix,name from regions order by 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var prefix, name string
		if err := rows.Scan(&prefix, &name); err != nil {
			panic(err)
		}
		fmt.Printf("%s,%s\n", prefix, name)
	}
}

func getRegion(db *sql.DB, search string) (name string, retErr error) {
	rows, err := db.Query("select name from regions where prefix = ?", search)
	if err != nil {
		retErr = err
		return
	}
	defer rows.Close()

	if rows.Next() {
		retErr = rows.Scan(&name)
	} else {
		retErr = ErrorNotFound
	}
	return
}
