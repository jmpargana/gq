GOBIN ?= $$(go env GOPATH)/bin

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml
	
visualise-coverage: check-coverage
	go tool cover -html=cover.out -o=cover.html
	open cover.html
	
tidy:
	go mod tidy -diff

test:
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

lint:
	golangci-lint run

build:
	go build -o ./bin/gq ./cmd/gq

.PHONY: tidy install-go-test-coverage test lint build
