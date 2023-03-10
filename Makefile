.PHONY: build test

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-chaos.version=${VERSION}'" -o conduit-connector-chaos cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...
