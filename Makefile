DB_USER := root
DB_PASS := password
DB_HOST := 127.0.0.1
DB_PORT := 3307
DB_NAME := portfolio
MARIADB_DSN := mariadb://${DB_USER}:${DB_PASS}@${DB_HOST}:$(DB_PORT)/${DB_NAME}

GOFILES=$(wildcard *.go **/*.go)

BINARY=./bin/traPortfolio
GO_RUN := ${BINARY} --db-user ${DB_USER} --db-pass ${DB_PASS} --db-host ${DB_HOST} --db-port ${DB_PORT} --db-name ${DB_NAME}

TEST_INTEGRATION_TAGS := "integration db"

GOLANGCI_LINT_VERSION := latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

TBLS_VERSION := latest
TBLS := TBLS_DSN=${MARIADB_DSN} go run github.com/k1LoW/tbls@${TBLS_VERSION}

SPECTRAL_VERSION := latest
SPECTRAL := docker run --rm -it -w /tmp -v $$PWD:/tmp stoplight/spectral:${SPECTRAL_VERSION}

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

.PHONY: check
check: all lint test-all db-lint openapi-lint

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

.PHONY: go-gen
go-gen:
	@go generate -x ./...

.PHONY: migrate
migrate: ${BINARY} # require test-db
	@${GO_RUN} --migrate

.PHONY: db-gen-docs
db-gen-docs: migrate
	@${RM} -rf ./docs/dbschema
	@${TBLS} doc

.PHONY: db-lint
db-lint: migrate
	@${TBLS} lint

.PHONY: up-test-db
up-test-db:
	@TEST_DB_PORT=$(DB_PORT) ./dev/bin/up-test-db.sh

.PHONY: rm-test-db
rm-test-db:
	@./dev/bin/down-test-db.sh

.PHONY: openapi-lint
openapi-lint:
	@${SPECTRAL} lint ./docs/swagger/traPortfolio.v1.yaml
