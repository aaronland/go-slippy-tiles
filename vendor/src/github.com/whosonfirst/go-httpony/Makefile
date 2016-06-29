CWD=$(shell pwd)
GOPATH := $(CWD)/vendor:$(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-httpony; then rm -rf src/github.com/whosonfirst/go-httpony; fi
	mkdir -p src/github.com/whosonfirst/go-httpony
	cp httpony.go src/github.com/whosonfirst/go-httpony/
	cp -r cors src/github.com/whosonfirst/go-httpony/
	cp -r tls src/github.com/whosonfirst/go-httpony/
	cp -r rewrite src/github.com/whosonfirst/go-httpony/
	cp -r crypto src/github.com/whosonfirst/go-httpony/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps fmt bin

deps:   self

vendor: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +

fmt:
	go fmt cmd/*.go
	go fmt *.go
	go fmt cors/*.go
	go fmt tls/*.go
	go fmt rewrite/*.go
	go fmt crypto/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/echo-pony cmd/echo-pony.go
