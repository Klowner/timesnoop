#!/bin/sh
GOPATH=${PWD} go build -o bin/activitymonitor app && strip bin/activitymonitor
