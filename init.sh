#!/bin/sh
GOPATH=$(pwd) go get -v app
rm bin/app
