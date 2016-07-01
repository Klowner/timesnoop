#!/bin/sh
webpack
go-bindata -prefix src/app -o src/app/bindata.go src/app/static
GOPATH=${PWD} go build -o bin/activitymonitor app && strip bin/activitymonitor
