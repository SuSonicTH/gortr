package main

import (
	"github.com/SuSonicTH/gortr/data/operator"
	"github.com/SuSonicTH/gortr/get"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Get struct {
	} `cmd:"" help:"get data from rtr"`

	Search struct {
		Number []string `arg:"" name:"number" help:"number to search"`
	} `cmd:"" help:"search number"`
}

func main() {
	ctx := kong.Parse(&CLI)

	var err error
	switch ctx.Command() {
	case "get":
		err = get.Numbers()
	case "search <number>":
		operator.ReadOperators()
	default:
		panic(ctx.Command())
	}

	if err != nil {
		panic(err)
	}
}
