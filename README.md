# traPortfolio

[![GitHub release](https://img.shields.io/github/release/traPtitech/traPortfolio.svg?logo=github)](https://GitHub.com/traPtitech/traPortfolio/releases/) [![CI](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/main.yaml) [![Build image](https://github.com/traPtitech/traPortfolio/actions/workflows/image.yaml/badge.svg)](https://github.com/traPtitech/traPortfolio/actions/workflows/image.yaml) [![codecov](https://codecov.io/gh/traPtitech/traPortfolio/branch/main/graph/badge.svg?token=2HB6P7RUX8)](https://codecov.io/gh/traPtitech/traPortfolio) [![OpenAPI](https://img.shields.io/badge/OpenAPI-apis.trap.jp-6BA539?logo=openapiinitiative)](https://apis.trap.jp/?urls.primaryName=traPortfolio)

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

### Quick start with DevContainer

If you use VSCode, you can use [DevContainer](https://code.visualstudio.com/docs/remote/containers) to develop traPortfolio.
See [./.devcontainer/README.md](./.devcontainer/README.md) for more details.

### Start docker container (with Docker Compose)

```bash
docker compose up
```

or

```bash
# enable live reload
docker compose watch
```

Now you can access to

- <http://localhost:1323> for backend server.
- <http://localhost:3001> for adminer
  - username: `root`
  - password: `password`
  - database: `portfolio`
  - port: `3306`

## Tasks

Usable tasks are below.

> [!TIP]
> You can use `xc` to run the following tasks easily.
> See <https://xcfile.dev> for more details.
>
> ```bash
> go install github.com/joerdav/xc/cmd/xc@latest
> ```

### gen

Generate code.

```bash
go generate -x ./...
```

### lint

Run linter (golangci-lint).

```bash
go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix ./...
```

### test:unit

Run unit tests.

```bash
go test -v -cover -race ./internal/...
```

### test:integration

Run integration tests.

```bash
go test -v -cover -race ./integration_tests/...
```

### test:all

Run all tests.

Requires: test:unit, test:integration

RunDeps: async

### db:migrate

Migrate the database.

```bash
# TODO: use environment variables for config
docker compose run --build --entrypoint "/traPortfolio -c /opt/traPortfolio/config.yaml --db-host mysql --only-migrate" backend
```

### db:gen-docs

Generate database schema documentation with tbls.

Requires: db:migrate

```bash
rm -rf ./docs/dbschema
go run github.com/k1LoW/tbls@latest doc
```

### db:lint

Lint the database schema with tbls.

Requires: db:migrate

```bash
go run github.com/k1LoW/tbls@latest lint
```

### openapi:lint

Lint the OpenAPI schema with Spectral.

```bash
docker run --rm -it -w /tmp -v $PWD:/tmp stoplight/spectral:latest lint ./docs/swagger/traPortfolio.v1.yaml
```
