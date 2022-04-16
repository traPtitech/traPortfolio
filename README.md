# traPortfolio

traPortfolio backend

## Develop environment

### Set up local develop environment with docker

#### Requirements

- docker
- docker-compose
- make
- bash

1. Run the following command in the project root

```bash
docker-compose up
```

Now you can access to

- `http://localhost:1323` for backend server.
- `http://localhost:3001` for Adminer
  - username: `root`
  - password: `password`
  - database: `portfolio`

### Set up local develop environment without docker

#### Requirements

- go 1.17
- mysql

1. Make sure MySQL is running
2. Run the following command in the project root

```bash
DB_HOST=localhost go run main.go
```

if you want to change DB port, set the `DB_PORT` environment variable.

### Rebuild

```bash
docker-compose up --build
```

### Up test DB

```bash
make up-test-db
```

### Remove test DB

```bash
make rm-test-db
```

### Generate DB docs

Make sure test MySQL container is running

```bash
make db-gen-docs
```

### Lint

DB Lint(you need docker)

```bash
make db-lint
```

golangci-lint(you need docker)

```bash
make golangci-lint
```
