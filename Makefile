# Ironically use a Makefile for a Makefile replacement :)

.PHONY: example
example:
	@echo "\n INFO: Starting AC with temperature 18 and fan speed 3:"
	@go run ./example/main.go ac --temperature 18 --fanspeed 3
	@echo "\n INFO: Starting engine quietly:"
	@go run ./example/main.go start --quiet false
	@echo "\n INFO: Checking status:"
	@go run ./example/main.go status

example-help:
	@go run ./example/main.go help

example-ac-help:
	@go run ./example/main.go help ac