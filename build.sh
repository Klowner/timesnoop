#!/bin/sh
node_modules/.bin/webpack \
	&& bin/go-bindata -prefix src/app -o src/app/bindata.go src/app/static \
	&& GOPATH=${PWD} go build -o bin/timesnoop app && strip bin/timesnoop
