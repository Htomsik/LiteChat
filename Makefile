.PHONY: build
start:
	swag init -g .\cmd\apiServer\main.go
	go mod tidy
	go build -v ./cmd/apiServer
	.\apiServer.exe

build:
	swag init -g .\cmd\apiServer\main.go
	go mod tidy
	go build -v ./cmd/apiServer


.DEFAULT_GOAL := build
