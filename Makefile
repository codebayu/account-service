# Makefile for Account Service

.PHONY: run test test-coverage coverage-html migrate tidy clean

# Variables
COVERAGE_FILE = coverage.out
TARGET_PACKAGES = ./internal/handler/...,./internal/repository/...,./internal/service/...,./internal/utils/...

# Run the API server
run:
	go run ./cmd/api

# Run all unit tests
test:
	go test -v ./...

# Run unit tests with filtered coverage (Target: >80%)
test-coverage:
	go test -v -coverpkg=$(TARGET_PACKAGES) -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

# Open the visual coverage report in browser
coverage-html: test-coverage
	go tool cover -html=$(COVERAGE_FILE)

# Run database migrations
migrate:
	go run ./cmd/migrate

# Tidy up Go modules
tidy:
	go mod tidy

# Clean up generated files
clean:
	rm -f $(COVERAGE_FILE)
	rm -f api
	rm -f migrate
