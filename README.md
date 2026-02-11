# GoRTR
a simple tool to download numbering data from the austrian [RTR](https://www.rtr.at/) with search functionality

This is mainly a lerning excercise to get into [go](https://go.dev/)

## Installation

### from the release page
just download the binary for your platform from the [releaeses](https://github.com/SuSonicTH/gortr/releases) and execute it, put it on your path for convinence.

### from source
If you have go installed and a working CGO environmant (see below) you can use following command to install the tool
```
go install github.com/SuSonicTH/gortr
```

As this tool uses on [go-sqlite3](https://github.com/mattn/go-sqlite3) a working c-compiler setup for go is needed and CGO_ENABLED=1 has to be set else the compilation will succseed but the resulting binary will print an error as soon as the sqlite module is used.
This should work out of the box on linux and some other OSes but not on Windows.

If you are on windows (or want to cross compile) you can try my [GOZ](https://github.com/SuSonicTH/goz) go wraper that handles (cross-)compilation of CGO packages.
execute:
```
go install github.com/SuSonicTH/goz
goz install github.com/SuSonicTH/gortr
```

## Usage
```
Usage of gortr:
  -listLocalAreas
        list all local areas with name
  -listParameter string
        list given parameter type use ALL for every parameter
  -listParameterTypes
        list all parameter types to use in -listParameter
  -localArea string
        serach for a matching local area
  -refresh
        get data from rtr.at
  -search string
        serach for a matching number (in national format with/without leading zero or international format with +43 or 0043)
```

On the first run, or whenever you execute it with the argument `--refresh`, the tool gets all the numbering data from [RTR](https://www.rtr.at/) and saves the data in you home directory under `.gortr/gortr.sqlite3`
The downloaded data is just a few MB but the data in the sqlite database stored locally is *almost 700MB* because of the huge number of single numbers in all the ranges.
