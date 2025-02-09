.PHONY: build test

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-chaos.version=${VERSION}'" -o conduit-connector-chaos cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools.go
	@go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -tI % go install %
	@go mod tidy

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: fmt
fmt:
	gofumpt -l -w .
	gci write --skip-generated  .

.PHONY: generate
generate:
	go generate ./...
	conn-sdk-cli readmegen -w

