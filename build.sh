#!/bin/sh
go build -ldflags "-s -w" -trimpath

FILE="gortr"
if [ -f "gortr.exe" ]; then
    FILE="gortr.exe"
fi

upx --lzma $FILE
