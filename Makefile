SHELL := /bin/bash
export GOBIN := $(CWD)/.bin
NAME=schema-registry-statistics

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o ${NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${NAME}-linux main.go

.PHONY: clean
clean:
	go clean
	rm ${NAME}-darwin
	rm ${NAME}-linux

.PHONY: test
test:
	go test -v ./...

.PHONY: test_coverage
test_coverage:
	go test -v ./... -coverprofile=coverage.out

.PHONY: dep
dep:
	go mod download

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: safe
safe:
	go vet
	go test -race .
	go build -race .

.PHONY: test_race
test_race:
	go run -race . --bootstrap "localhost:9092" \
	--topic "payments-topic" \
	--group "TEST_GROUP" \
	--tls --cert "ca.pem" \
	--user "USERNAME" \
	--password "PASSWORD" \
	--oldest --verbose

.PHONY: lint
lint:
	golangci-lint run