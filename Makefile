CWD=$(shell pwd)
VENDORGOPATH := $(CWD)/vendor:$(CWD)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/thisisaaronland/go-slippy-tiles; then rm -rf src/github.com/thisisaaronland/go-slippy-tiles; fi
	mkdir -p src/github.com/thisisaaronland/go-slippy-tiles/provider
	mkdir -p src/github.com/thisisaaronland/go-slippy-tiles/cache
	cp slippytiles.go src/github.com/thisisaaronland/go-slippy-tiles/
	cp provider/*.go src/github.com/thisisaaronland/go-slippy-tiles/provider/
	cp cache/*.go src/github.com/thisisaaronland/go-slippy-tiles/cache/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps deps bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/jtacoma/uritemplates"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-httpony"

bin:	self fmt
	@GOPATH=$(GOPATH) go build -o bin/tile-proxy cmd/tile-proxy.go

fmt:
	go fmt cmd/*.go
	go fmt cache/*.go
	go fmt provider/*.go
	go fmt *.go

