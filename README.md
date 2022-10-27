# Course Watch API
## Prerequisites
### Swag
Swag (https://github.com/swaggo/swag) is required to regenerate Swagger
files:

`go install github.com/swaggo/swag/cmd/swag@latest`

The initial version is built using swag v1.8.6
### gomock

`go install github.com/golang/mock/mockgen@v1.6.0`

Example usage on Windows: 

`mockgen -source=".\internal\delivery\http\auth.go" -destination=".\internal\delivery\http\mocks\mock_auth.go"`

Each file which requires mocks must be added to the `gen` target in Makefile