FROM golang:latest

WORKDIR "/api"

COPY ./ .
RUN go build -v ./cmd/Server
CMD ["./apiServer"]