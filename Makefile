# Ironically use a Makefile for a Makefile replacement :)

.PHONY: example
example:
	@echo "Showing general help:"
	@go run ./example/main.go help
	@echo "\nShowing help for AC command:"
	@go run ./example/main.go help ac
	@echo "\nStarting AC with temperature 18 and fan speed 3:"
	@go run ./example/main.go ac --temperature 18 --fanspeed 3
	@echo "\nStarting engine quietly:"
	@go run ./example/main.go engine start --quiet
	@echo "\nChecking status:"
	@go run ./example/main.go status

example-help:
	@go run ./example/main.go help

example-ac-help:
	@go run ./example/main.go help ac