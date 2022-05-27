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

.PHONY: ${shell egrep -o ^[a-zA-Z_-]+: ./Makefile | sed 's/://'}

all: clean mod build

clean:
	@$(RM) $(BINARY)
	@go clean

mod:
	@go mod tidy

build: $(GOFILES)
	@go build -o $(BINARY)

check: all lint test-all db-lint openapi-lint

test: $(GOFILES)
	go test -v -cover -race ./...

test-all: $(GOFILES)
	go test -v -cover -race -tags=$(TEST_INTEGRATION_TAGS) ./...

test-integration-db: $(GOFILES)
	go test -v -cover -race -tags=$(TEST_INTEGRATION_TAGS) ./integration_tests/...

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
	@TEST_DB_PORT=$(DB_PORT) ./dev/bin/up-test-db.sh

rm-test-db:
	@./dev/bin/down-test-db.sh

openapi-lint:
	@${SPECTRAL} lint ./docs/swagger/traPortfolio.v1.yaml
