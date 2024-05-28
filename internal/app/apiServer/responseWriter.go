package apiServer

import "net/http"

// responseWriter
type responseWriter struct {
	http.ResponseWriter
	code int
}

func (writer *responseWriter) WriteHeader(code int) {
	writer.code = code
	writer.ResponseWriter.WriteHeader(code)
}
