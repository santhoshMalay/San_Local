.DEFAULT_GOAL := build

.PHONY: swag
swag:
	swag init -g internal/apiserver/apiserver.go -o api/swagger

.PHONY: build
build:
	go mod download && CGO_ENABLED=0 go build -o ./.bin/apiserver ./cmd/apiserver

.PHONY: gen
gen:
	mockgen -source=internal/delivery/http/auth.go -destination=internal/delivery/http/mocks/mock_auth.go

