export GOPATH := $(shell pwd)
export PATH := $(GOPATH)/bin:$(PATH)

PREFIX    := $(GOPATH)
HEAD = `git rev-parse --short HEAD`

all: build

build: LDFLAGS   += $(shell GOPATH=${GOPATH})
build: clean
	@echo "--> Building Gow Server..."
	@mkdir -p $(PREFIX)/bin/
	go build -v -o $(PREFIX)/bin/gow  --ldflags "-w -s -X main.git=$(HEAD)" gow.go
	@chmod 755 $(PREFIX)/bin/gow

clean:
	@echo "--> Cleaning..."
	@go clean
	@rm -f $(PREFIX)/bin/gow

fmt:
	go fmt ./...
	go vet ./...

.PHONY: clean fmt