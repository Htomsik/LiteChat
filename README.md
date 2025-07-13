

<h2 align="center">LiteChat</h3>


<h3 align="center">

Simple text chat on websockets on GO + Vue

</h3>

## Required
- npm (for backend)
- Make (if you want to use build scripts)

##  Usage

#### Build With Make

``` bash
# Full Build + Autostart
Make start 

## Full Build
MAKE build 

# Or if you want build only front/back
MAKE backendBuild
MAKE frontendBuild  # REQUIRED IF WANT HOST LIKE SPA
```
#### Manual build

``` bash

## Backend
go mod tidy
go build -v ./cmd/apiServer

## Frontend
cd website && npm install && npm run build # REQUIRED IF WANT HOST LIKE SPA
```


### Startup
``` bash
./apiServer # For back
cd website && npm run dev # If you want only front
```
