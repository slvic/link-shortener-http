# Link Shortener

## Requirements

The project has been tested on Go 1.17. It uses some new features such as `go install`, `go:embed` (for third party
dependencies like Swagger UI). Other binary dependencies will be downloaded to the `bin` folder.

## First steps

### Third party dependencies

To download third party dependencies (e.g. buf, proto-gen-go, goose, Swagger UI) run `make deps`.

### Build project

Run `make build`. It will download Swagger UI if it does not exist and build `sinkapi/main.go` file.

### Lint before commit

Run `make lint`. It will check *.proto files with [buf](https://buf.build/) and *.go files with golangci-lint

### Set up dependent services

Run `docker compose up` to run PostgreSQL.

### Run tests

```bash
# Run unit tests
make test

# Run integration tests
# Paste your connection string. 
# If you ran docker compose up, you can copy-paste the code below.
export DSN="user=postgres password=postgres database=api sslmode=disable" 

# run migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DSN make migrate

# run integration tests
make test-integration
``` 

### Create migration

```bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DSN make migration
```

The new migration will be added to the migration folder. Do not forget to rename it.

### Apply migrations

```bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DSN make migrate
```
