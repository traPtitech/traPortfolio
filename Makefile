TEST_DB_PORT := 3306
TBLS_VERSION := 1.38.3

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint run

.PHONY: up-test-db
up-test-db:
	@TEST_DB_PORT=$(TEST_DB_PORT) ./dev/bin/up-test-db.sh

.PHONY: db-gen-docs
db-gen-docs:
	@./dev/bin/db-gen-docs.sh $(TEST_DB_PORT) $(TBLS_VERSION)

.PHONY: rm-test-db
rm-test-db:
	@./dev/bin/down-test-db.sh