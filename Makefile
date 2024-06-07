.PHONY: build

build:
	go mod tidy
	go build -v ./cmd/apiServer


swag:
	swag init -g ./cmd/apiServer/main.go

.PHONY: test
test:
	go test -v -race -timeout 15s ./...


.DEFAULT_GOAL := build
