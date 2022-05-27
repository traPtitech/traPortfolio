TEST_DB_PORT := 3307

TBLS_VERSION := 1.49.6

GOFILES=$(wildcard *.go **/*.go)

BINARY=./bin/traPortfolio

TEST_INTEGRATION_TAGS := "integration db"

GOLANGCI_LINT_VERSION := latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

.PHONY: all
all: clean mod build

.PHONY: clean
clean:
	@$(RM) $(BINARY)
	@go clean

.PHONY: mod
mod:
	@go mod tidy

.PHONY: build
build: $(GOFILES)
	@go build -o $(BINARY)

.PHONY: test
test: $(GOFILES)
	go test -v -cover -race ./...

.PHONY: test-all
test-all: $(GOFILES)
	go test -v -cover -race -tags=$(TEST_INTEGRATION_TAGS) ./...

.PHONY: test-integration-db
test-integration-db: $(GOFILES)
	go test -v -cover -race -tags=$(TEST_INTEGRATION_TAGS) ./integration_tests/...

.PHONY: lint
lint:
	@${GOLANGCI_LINT} run --fix ./...

.PHONY: up-test-db
up-test-db:
	@TEST_DB_PORT=$(TEST_DB_PORT) ./dev/bin/up-test-db.sh

.PHONY: db-gen-docs
db-gen-docs:
	# @./dev/bin/db-gen-docs.sh $(TEST_DB_PORT) $(TBLS_VERSION)
	@if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	DB_HOST=localhost DB_PORT=$(TEST_DB_PORT) go run main.go -migrate true
	docker run -u $$(id -u) --rm --net=host -e TBLS_DSN="mariadb://root:password@127.0.0.1:$(TEST_DB_PORT)/portfolio" -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) doc

.PHONY: db-lint
db-lint:
	DB_HOST=localhost DB_PORT=$(TEST_DB_PORT) go run main.go -migrate true
	docker run --rm --net=host -e TBLS_DSN="mariadb://root:password@127.0.0.1:$(TEST_DB_PORT)/portfolio" -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) lint

.PHONY: rm-test-db
rm-test-db:
	@./dev/bin/down-test-db.sh

.PHONY: go-gen
go-gen:
	@go generate ./...
