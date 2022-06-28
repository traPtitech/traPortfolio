TEST_DB_USER := root
TEST_DB_PASS := password
TEST_DB_HOST := 127.0.0.1
TEST_DB_PORT := 3307
TEST_DB_NAME := portfolio
TEST_MARIADB_DSN := mariadb://${TEST_DB_USER}:${TEST_DB_PASS}@${TEST_DB_HOST}:${TEST_DB_PORT}/${TEST_DB_NAME}

GOFILES=$(wildcard *.go **/*.go)
BINARY=./bin/traPortfolio
GO_RUN := ${BINARY} --db-user ${TEST_DB_USER} --db-pass ${TEST_DB_PASS} --db-host ${TEST_DB_HOST} --db-port ${TEST_DB_PORT} --db-name ${TEST_DB_NAME}
GOTEST_FLAGS := -v -cover -race

GOLANGCI_LINT_VERSION := latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

TBLS_VERSION := latest
TBLS := TBLS_DSN=${TEST_MARIADB_DSN} go run github.com/k1LoW/tbls@${TBLS_VERSION}

SPECTRAL_VERSION := latest
SPECTRAL := docker run --rm -it -w /tmp -v $$PWD:/tmp stoplight/spectral:${SPECTRAL_VERSION}

.PHONY: ${shell egrep -o ^[a-zA-Z_-]+: ./Makefile | sed 's/://'}

all: clean mod build

clean:
	@${RM} ${BINARY}
	@go clean

mod:
	@go mod tidy

build: ${GOFILES}
	@go build -o ${BINARY}

check: all lint test-all db-lint openapi-lint

test: ${GOFILES}
	go test ${GOTEST_FLAGS} $$(go list ./... | grep -v "integration_tests")

test-all: ${GOFILES}
	go test ${GOTEST_FLAGS} ./...

test-integration-db: ${GOFILES}
	go test ${GOTEST_FLAGS} ./integration_tests/...

lint:
	@${GOLANGCI_LINT} run --fix ./...

go-gen:
	@go generate -x ./...

migrate: ${BINARY} # require test-db
	@${GO_RUN} --migrate

db-gen-docs: migrate
	@${RM} -rf ./docs/dbschema
	@${TBLS} doc

db-lint: migrate
	@${TBLS} lint

up-test-db:
	@TEST_DB_PORT=${TEST_DB_PORT} ./dev/bin/up-test-db.sh

rm-test-db:
	@./dev/bin/down-test-db.sh

openapi-lint:
	@${SPECTRAL} lint ./docs/swagger/traPortfolio.v1.yaml
