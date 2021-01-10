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

```
$ docker-compose up
```

Now you can access to 

- `http://localhost:1323` for backend server.
- `http://localhost:3001` for Adminer
  - username: `root`
  - password: `password`
  - database: `portfolio`

### Set up local develop environment without docker

#### Requirements

- go 1.15
- mysql

1. Make sure MySQL is running
2. Run the following command in the project root

```
$ DB_HOST=localhost go run main.go
```

if you want to change DB port, set the `DB_PORT` environment variable.

### Rebuild

```
$ doker-compose up --build
```

### Up test DB

```
$ make up-test-db
```

### Remove test DB

```
$ make rm-test-db
```

### Generate DB docs

1. Make sure test MySQL container is running
2. `make db-gen-docs`

### Lint

DB Lint(you need docker)

```
$ make db-lint
```

golangci-lint

```
$ make golangci-lint
```

golangci-lint with docker
```
$ make docker-golangci-lint
```