package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/SuSonicTH/gortr/data/numbers"
	"github.com/SuSonicTH/gortr/data/region"
	"github.com/SuSonicTH/gortr/get"
	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE = "./gortr.sqlite3"

func main() {
	showHelp := true
	pRefresh := flag.Bool("refresh", false, "get data from rtr.at")
	pRegion := flag.String("region", "", "match given number to a region")
	pSearch := flag.String("search", "", "serach for a matching number")

	flag.Parse()

	if *pRefresh {
		showHelp = false
		refresh()
	}

	if *pRegion != "" {
		showHelp = false
		searchReagon(*pRegion)
	}

	if *pSearch != "" {
		showHelp = false
		db := openDb()
		defer db.Close()

		//searchNumber(db,*pSearch)
		searchNumber(*pSearch)
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

	if err := get.Numbers(db); err != nil {
		panic(err)
	}
}

func searchReagon(search string) {
	reg, err := region.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("prefix: 0%s\n", reg.Prefix)
	fmt.Printf("name:   %s\n", reg.Name)
}

func searchNumber(search string) {
	if err := numbers.Load(); err != nil {
		panic(err)
	}

	number, err := numbers.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	if number.PfxFrom == number.PfxTo {
		fmt.Printf("number: 0%s%s\n", number.Prefix, number.PfxFrom)
	} else {
		fmt.Printf("number: 0%s%s - 0%s%s\n", number.Prefix, number.PfxFrom, number.Prefix, number.PfxTo)

	}
	fmt.Printf("type: %s\n", number.NumberType.Name)
	fmt.Printf("operator: %s - %s\n", number.Operator.Id, number.Operator.Name)
}
