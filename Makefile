
.PHONY: backendBuild frontendBuild test
.DEFAULT_GOAL := start

start: frontendBuild backendBuild killBackend
	./apiServer.exe

backendBuild:
	go mod tidy
	go build -v ./cmd/apiServer

killBackend:
	-taskkill /IM apiServer.exe /F 2>NUL || exit 0

backendBuildSwag:
	swag init -g ./cmd/apiServer/main.go

frontendBuild:
	cd website && npm install && npm run build

test:
	go test -v -race -timeout 15s ./...

