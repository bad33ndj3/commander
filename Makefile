# Ironically use a Makefile for a Makefile replacement :)

.PHONY: example example-help example-ac-help test lint

example:
	@echo "\nStarting AC with temperature 18 and fan speed 3:"
	@go run ./example/main.go ac --temperature 18 --fanspeed 3
	@echo "\nStarting engine quietly:"
	@go run ./example/main.go start --quiet
	@echo "\nChecking status:"
	@go run ./example/main.go status

example-help:
	@go run ./example/main.go help

example-ac-help:
	@go run ./example/main.go help ac

test:
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func=coverage.txt

lint:
	@golangci-lint run ./...

# Run all checks
check: lint test