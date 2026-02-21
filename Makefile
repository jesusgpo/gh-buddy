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

test:
	go test ./... -v

lint:
	golangci-lint run ./...

release:
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64  .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64  .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
