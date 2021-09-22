TEST_DB_PORT := 3307
TBLS_VERSION := 1.38.3

GOFILES=$(wildcard *.go **/*.go)

BINARY=./bin/traPortfolio

.PHONY: all
all: clean build

.PHONY: test
test: $(GOFILES) 
	go test -v -cover -race ./...

.PHONY: build
build: $(GOFILES)
	go build -o $(BINARY)

.PHONY: clean
clean: ## Cleanup working directory
	@$(RM) $(BINARY)
	@go clean


.PHONY: golangci-lint
golangci-lint:
	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run

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
	docker run -u $$(id -u) --rm --net=host -e TBLS_DSN="mysql://root:password@127.0.0.1:$(TEST_DB_PORT)/portfolio" -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) doc

.PHONY: db-lint
db-lint:
	DB_HOST=localhost DB_PORT=$(TEST_DB_PORT) go run main.go -migrate true
	docker run --rm --net=host -e TBLS_DSN="mysql://root:password@127.0.0.1:$(TEST_DB_PORT)/portfolio" -v $$PWD:/work k1low/tbls:$(TBLS_VERSION) lint

.PHONY: rm-test-db
rm-test-db:
	@./dev/bin/down-test-db.sh

.PHONY: go-gen
go-gen:
	@go generate ./...