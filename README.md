# Course Watch API
## Prerequisites
### Swag
Swag (https://github.com/swaggo/swag) is required to regenerate Swagger
files:

`go install github.com/swaggo/swag/cmd/swag@latest`

The initial version is built using swag v1.8.6
### gomock

`go install github.com/golang/mock/mockgen@v1.6.0`

Each go file which is expected to yield mocks must contain a go `generate` directive, e.g.:

`//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go`

This example is from `internal/delivery/http/auth.go`

With the directive in place, run `go generate ./...` or use the `gen` target in the Makefile whenever mocks need to be updated 

Note: To run manually for a specific file from Windows console use the following syntax 

`mockgen -source=".\internal\delivery\http\auth.go" -destination=".\internal\delivery\http\mocks\mock_auth.go"`
