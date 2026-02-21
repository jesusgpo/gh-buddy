BINARY_NAME = gh-buddy
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -X github.com/jesusgpo/gh-buddy/cmd.version=$(VERSION)"

.PHONY: build install clean test lint release

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build
	gh extension install .

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

test:
	go test ./... -v

lint:
	golangci-lint run ./...

release:
	goreleaser release --clean

release-snapshot:
	goreleaser release --snapshot --clean
