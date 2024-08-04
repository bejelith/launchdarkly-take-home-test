all: coverage test build

PHONY: server
server: 
	go build -o server ./main

test:
	go test -v -cover -coverprofile cover.out  ./...

coverage: test
	go tool cover -html=cover.out -o coverage.html

build: server
