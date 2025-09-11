package server

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// responseWriterWrapper captures status code for logging
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrap ResponseWriter to capture status
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)

		level := "✅"
		if rw.statusCode >= 400 {
			level = "❌"
		}

		log.Printf("%s [%s] %s | %d | %s",
			level,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			elapsed,
		)
	})
}
