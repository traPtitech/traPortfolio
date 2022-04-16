# traPortfolio

[![GitHub release](https://img.shields.io/github/release/traPtitech/traPortfolio.svg)](https://GitHub.com/traPtitech/traPortfolio/releases/) [![CI](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml) [![Build image](https://github.com/traPtitech/traPortfolio/actions/workflows/release.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/release.yaml) [![codecov](https://codecov.io/gh/traPtitech/traPortfolio/branch/main/graph/badge.svg?token=2HB6P7RUX8)](https://codecov.io/gh/traPtitech/traPortfolio) [![swagger](https://img.shields.io/badge/swagger-docs-brightgreen)](https://apis.trap.jp/?urls.primaryName=traPortfolio)

- Backend
  - [traPtitech/traPortfolio](https://github.com/traPtitech/traPortfolio) (this repository)
- Frontend
  - [traPtitech/traPortfolio-UI](https://github.com/traPtitech/traPortfolio-UI)
  - [traPtitech/traPortfolio-Dashboard](https://github.com/traPtitech/traPortfolio-Dashboard)

## Develop environment

### Requirements

- bash
- make
- docker
- docker-compose
- go 1.17
- mysql

### Set up with docker-compose

#### Start docker container

Run the following command in the project root

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

(Optional) After running the following command, sample data will be inserted into the database

```bash
go run main.go --migrate --db-user root --db-pass password --db-port 3306 --db-host localhost --db-name portfolio
```

#### Rebuild docker container

Run the following command in the project root

```bash
docker-compose up --build
```

### Set up without docker-compose

#### Set Up test DB (with port 3307)

Run the following command in the project root

```bash
make up-test-db
```

(Optional) After running the following command, sample data will be inserted into the database

```bash
go run main.go --migrate --db-user root --db-pass password --db-port 3307 --db-host localhost --db-name portfolio
```

#### Remove test DB

Run the following command in the project root

```bash
make rm-test-db
```

#### Run locally

Make sure test MySQL container is running,
and run the following command in the project root

```bash
go run main.go --db-user root --db-pass password --db-port 3307 --db-host localhost --db-name portfolio
```

### Generate DB docs

Make sure test MySQL container is running,
and run the following command in the project root

```bash
make db-gen-docs
```

### Run linters

#### DB linter (tbls)

Make sure test MySQL container is running,
and run the following command in the project root

```bash
make db-lint
```

#### Go linter (golangci-lint)

Run the following command in the project root

```bash
make golangci-lint
```
