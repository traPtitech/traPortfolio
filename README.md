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

### Quick start with DevContainer (Recommended)

If you use VSCode, you can use [DevContainer](https://code.visualstudio.com/docs/remote/containers) to develop traPortfolio.
See [./.devcontainer/README.md](./.devcontainer/README.md) for more details.

### Requirements (for local development)

- Bash
- make
- Docker
- Docker Compose
- Go
- MySQL

### Start docker container (with Docker Compose)

```bash
docker compose up
```

Tips: You can change the configuration by rewriting [./dev/config.yaml](./dev/config.yaml)

Now you can access to

- <http://localhost:1323> for backend server.
- <http://localhost:3001> for adminer
  - username: `root`
  - password: `password`
  - database: `portfolio`
  - port: `3306`

### Run locally

**NOTE**: Make sure `backend` container is stopped because of port conflicts

```bash
make up-db-container
go run main.go -c ./dev/config.yaml
```

Tips: You can change the configuration by

- Specifying it with flags (Run `go run main.go --help`)
- Rewriting [./dev/config.yaml](dev/config.yaml)

### Generate DB docs

```bash
make db-gen-docs
```

### Run linters

#### DB linter (tbls)

```bash
make db-lint
```

#### OpenAPI linter (spectral)

```bash
make openapi-lint
```

#### Go linter (golangci-lint)

```bash
make lint
```

### Run tests

#### Unit tests

```bash
make test
```

#### Integration tests

```bash
make test-integration
```

If you want to test both of them, run the following command.

```bash
make test-all
```
