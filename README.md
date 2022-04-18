# traPortfolio

[![GitHub release](https://img.shields.io/github/release/traPtitech/traPortfolio.svg)](https://GitHub.com/traPtitech/traPortfolio/releases/) [![CI](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml) [![Build image](https://github.com/traPtitech/traPortfolio/actions/workflows/release.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/release.yaml) [![codecov](https://codecov.io/gh/traPtitech/traPortfolio/branch/main/graph/badge.svg?token=2HB6P7RUX8)](https://codecov.io/gh/traPtitech/traPortfolio) [![swagger](https://img.shields.io/badge/swagger-docs-brightgreen)](https://apis.trap.jp/?urls.primaryName=traPortfolio)

- Backend
  - [traPtitech/traPortfolio](https://github.com/traPtitech/traPortfolio) (this repository)
- Frontend
  - [traPtitech/traPortfolio-UI](https://github.com/traPtitech/traPortfolio-UI)
  - [traPtitech/traPortfolio-Dashboard](https://github.com/traPtitech/traPortfolio-Dashboard)

## Develop environment

If you want to contribute to traPortfolio, then follow these pages.

- [Architecture memo (in Japanese)](./docs/architecture.md)
- [API schema](./docs/swagger/traPortfolio.v1.yaml)
- [DB schema](./docs/dbschema)

### Requirements

- bash
- make
- docker
- docker-compose
- go 1.17
- mysql

### Start docker container (with docker-compose)

```bash
docker-compose up
```

Now you can access to

- `http://localhost:1323` for backend server.
- `http://localhost:3001` for Adminer
  - username: `root`
  - password: `password`
  - database: `portfolio`
  - port: `3306`

(Optional) By running the server with a `--migrate` flag, sample data will be inserted into the database.

### Set Up test DB (with docker, port:3307)

```bash
make up-test-db
```

(Optional) By running the server with a `--migrate` flag, sample data will be inserted into the database.

### Remove test DB

```bash
make rm-test-db
```

### Run locally

Make sure test MySQL container is running.

```bash
go run main.go --db-user root --db-pass password --db-port 3307 --db-host localhost --db-name portfolio
```

### Generate DB docs

Make sure test MySQL container is running.

```bash
make db-gen-docs
```

### Run linters

#### DB linter (tbls)

Make sure test MySQL container is running.

```bash
make db-lint
```

#### Go linter (golangci-lint)

```bash
make golangci-lint
```
