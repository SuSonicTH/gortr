#!/bin/sh
go build -ldflags "-s -w" -trimpath
upx --lzma gortr

