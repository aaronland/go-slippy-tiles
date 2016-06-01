CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

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

build:	rmdeps bin

deps:
	@GOPATH=$(shell pwd) go get -u "github.com/jtacoma/uritemplates"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-httpony"

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/tile-proxy cmd/tile-proxy.go

vendor: rmdeps deps
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

fmt:
	go fmt cmd/*.go
	go fmt cache/*.go
	go fmt provider/*.go
	go fmt *.go

