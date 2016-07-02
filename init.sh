#!/bin/sh
GOPATH=$(pwd) go get -v app
GOPATH=$(pwd) go get -u github.com/jteeuwen/go-bindata/...
rm bin/app
npm install
