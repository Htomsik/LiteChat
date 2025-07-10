
.PHONY: backendBuild frontendBuild test
.DEFAULT_GOAL := backendBuild

backendBuild:
	go mod tidy
	go build -v ./cmd/apiServer
	swag init -g ./cmd/apiServer/main.go

frontendBuild:
	cd website && npm install && npm run build

test:
	go test -v -race -timeout 15s ./...

