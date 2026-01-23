test:
	go test ./...

lint:
	golangci-lint run

build:
	go build ./...

.PHONY: test lint build
