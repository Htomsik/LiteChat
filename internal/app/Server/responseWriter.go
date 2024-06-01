package Server

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// responseWriter
type responseWriter struct {
	http.ResponseWriter
	code int
}

func (writer *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {

	h, ok := writer.ResponseWriter.(http.Hijacker)

	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}

func (writer *responseWriter) WriteHeader(code int) {
	writer.code = code
	writer.ResponseWriter.WriteHeader(code)
}
