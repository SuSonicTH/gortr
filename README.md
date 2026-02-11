# GoRTR
a simple tool to download numbering data from the austrian [RTR](https://www.rtr.at/) with search functionality

This is mainly a lerning excercise to get into [go](https://go.dev/)

## installation

If you have go installed and a working CGO environmant (see below) you can use following command to install the tool
```
go install github.com/SuSonicTH/gortr
```

As this tool uses on [go-sqlite3](https://github.com/mattn/go-sqlite3) a working c-compiler setup for go is needed and CGO_ENABLED=1 has to be set else the compilation might succseed but the resulting binary will print an error.
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
        serach for a matching number
```
