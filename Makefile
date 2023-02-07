GOPATH:=$(shell go env GOPATH)
GOBIN:=$(shell pwd)/bin
CGO_ENABLED=0

DATE?=$(shell date -u "+%Y-%m-%d %H:%M:%S")
LDFLAGS=-s -w -X 'main.AppVersion=${app_version}' -X 'main.BuildDate=${DATE}'

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o bin/app -mod=readonly *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: lint
lint:
	golangci-lint run


