test:
	go test ./...

lint:
	golangci-lint run

build:
	go build -o ./bin/gq ./cmd/gq

.PHONY: test lint build
