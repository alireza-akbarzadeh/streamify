swagger:
	swag init --parseDependency --parseInternal -g cmd/main.go

APP_NAME=streamify
BIN_DIR=bin
SRC_DIR=.
SQLC_DIR=sqlc.yml
MIGRATION_DIR=internal/sql/schema
DB_URL=postgres://postgres:postgres@localhost:5432/streamify

COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

.PHONY: build run migrate generate dev \
	test test-v test-race test-cover test-cover-html clean \
	vendor tidy-vendor build-vendor test-vendor

# Build the app
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) cmd/main.go

# Run the app (ensure build first)
run: build
	$(BIN_DIR)/$(APP_NAME)

# Run DB migrations
migrate:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

# Generate sqlc code
generate:
	sqlc generate -f $(SQLC_DIR)

# Development mode (build + run)
dev: run

# Test targets
test:
	go test ./...

test-v:
	go test -v ./...

test-race:
	go test -race ./...

test-cover:
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_OUT) ./...

test-cover-html: test-cover
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

clean:
	rm -rf $(COVERAGE_DIR) $(BIN_DIR)

# Sync vendor directory with go.mod/go.sum
vendor:
	go mod vendor

# Tidy and vendor in one step
tidy-vendor:
	go mod tidy
	go mod vendor

# Build using vendor directory
build-vendor:
	@mkdir -p $(BIN_DIR)
	GOFLAGS=-mod=vendor go build -o $(BIN_DIR)/$(APP_NAME) cmd/main.go

# Test using vendor directory
test-vendor:
	GOFLAGS=-mod=vendor go test ./...
