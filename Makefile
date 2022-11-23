.DEFAULT_GOAL := build

.PHONY: swag
swag:
	swag init --parseDependency --dir internal/apiserver -g apiserver.go -o api/swagger

.PHONY: build
build:
	go mod download && CGO_ENABLED=0 go build -o ./.bin/apiserver ./cmd/apiserver

.PHONY: gen
gen:
	go generate ./...

