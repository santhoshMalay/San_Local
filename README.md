# Course Watch API
## Prerequisites
### Swag
Swag (https://github.com/swaggo/swag) is required to regenerate Swagger
files:

`go install github.com/swaggo/swag/cmd/swag@latest`

The current version is built using swag v1.8.7

In order to rebuild swagger files run:

`swag init --parseDependency --dir internal/apiserver -g apiserver.go -o api/swagger`

Alternatively, use `swag` target in the Makefile

### gomock

`go install github.com/golang/mock/mockgen@v1.6.0`

Each go file which is expected to yield mocks must contain a go `generate` directive, e.g.:

`//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go`

This example is from `internal/delivery/http/auth.go`

With the directive in place, run `go generate ./...` or use the `gen` target in the Makefile whenever mocks need to be updated 

Note: To run manually for a specific file from Windows console use the following syntax 

`mockgen -source=".\internal\delivery\http\auth.go" -destination=".\internal\delivery\http\mocks\mock_auth.go"`

### postgres

To run a database in docker:

`docker run -e POSTGRES_PASSWORD=$env:POSTGRES_PASSWORD -p 6543:5432 -d --rm postgres:15.1`

Applying migrations with [migrate](https://github.com/golang-migrate/migrate) tool.

*Installation on Windows:*

Using scoop

`$ scoop install migrate`

With go toolchain

```
$ go get -u -d github.com/golang-migrate/migrate/cmd/migrate
$ cd $env:GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate
$ git checkout $TAG  # e.g. v4.1.0
$ go install -tags 'pgx' github.com/golang-migrate/migrate/v4/cmd/migrate@$TAG
```

*Basic usage:*

`$ migrate -source file://path/to/migrations -database pgx://postgres:$env:POSTGRES_PASSWORD@localhost:6543/postgres up`