package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SuSonicTH/gortr/data/numbers"
	"github.com/SuSonicTH/gortr/data/region"
	"github.com/SuSonicTH/gortr/get"
)

func main() {
	showHelp := true
	pRefresh := flag.Bool("refresh", false, "get data from rtr.at")
	pRegion := flag.String("region", "", "match given number to a region")
	pSearch := flag.String("search", "", "serach for a matching number")

	flag.Parse()

	if *pRefresh {
		showHelp = false
		get.Numbers()
	}
	if *pRegion != "" {
		showHelp = false
		reg, err := region.Search(*pRegion)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("prefix: 0%s\n", reg.Prefix)
			fmt.Printf("name:   %s\n", reg.Name)
		}
	}

	if *pSearch != "" {
		showHelp = false
		if err := numbers.Load(); err != nil {
			panic(err)
		}
		number, err := numbers.Search(*pSearch)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%+v\n", *number)
			//fmt.Printf("prefix: 0%s\n", reg.Prefix)
			//fmt.Printf("name:   %s\n", reg.Name)
		}
	}

	if showHelp {
		fmt.Printf("no argument given\n")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
}
