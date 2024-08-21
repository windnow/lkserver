package lkserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	code int
	size int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	w.size = n
	return
}
