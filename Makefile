# Template for a standard Makefile in a go project, pretty much just works
 
GO ?= go
GOLINT ?= golint
GOPATH := $(GOPATH)
 
default: test
 
build:
	$(GO) build
 
test: build fmt lint
	$(GO) test -v
 
# Requires go lint tool, install w/ `go get github.com/golang/lint/golint`
lint:
	$(GOLINT) .
 
fmt:
	$(GO) fmt
 
clean:
	$(GO) clean
