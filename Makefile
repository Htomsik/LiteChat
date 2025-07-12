
.PHONY: backendBuild frontendBuild test
.DEFAULT_GOAL := start

start: build killBackend
	./apiServer.exe

build: frontendBuild backendBuild

backendBuild:
	go mod tidy
	go build -v ./cmd/apiServer

killBackend:
	-taskkill /IM apiServer.exe /F 2>NUL || exit 0

backendInitSwagger:
	swag init -g ./cmd/apiServer/main.go

backendTests:
	go test -v -race -timeout 15s ./...

frontendBuild:
	cd website && npm install && npm run build


