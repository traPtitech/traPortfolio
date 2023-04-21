GOFILES=$(wildcard *.go **/*.go)
BINARY=./bin/traPortfolio
GO_RUN := ${BINARY} -c ./dev/config.yaml
GOTEST_FLAGS := -v -cover -race

GOLANGCI_LINT_VERSION := latest
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

TBLS_VERSION := latest
TBLS := go run github.com/k1LoW/tbls@${TBLS_VERSION}

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

# `test` is an alias for `test-unit`
test: ${GOFILES} test-unit

test-unit: ${GOFILES}
	go test ${GOTEST_FLAGS} $$(go list ./... | grep -v "integration_tests")

test-integration: ${GOFILES}
	go test ${GOTEST_FLAGS} ./integration_tests/...

test-all: ${GOFILES}
	go test ${GOTEST_FLAGS} ./...

lint:
	@${GOLANGCI_LINT} run --fix ./...

go-gen:
	@go generate -x ./...

migrate: ${BINARY}
	@${GO_RUN} --only-migrate

db-gen-docs: migrate
	@${RM} -rf ./docs/dbschema
	@${TBLS} doc

db-lint: migrate
	@${TBLS} lint

openapi-lint:
	@${SPECTRAL} lint ./docs/swagger/traPortfolio.v1.yaml
