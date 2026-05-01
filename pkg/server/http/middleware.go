package http

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// RequestMiddleware logs HTTP requests without query strings, avoiding accidental token exposure.
func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == HealthEndpoint {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		log.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lrw.statusCode).
			Dur("duration", time.Since(start)).
			Msg("handled HTTP request")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.headerWritten {
		return
	}
	lrw.statusCode = code
	lrw.headerWritten = true
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	if !lrw.headerWritten {
		lrw.WriteHeader(http.StatusOK)
	}
	return lrw.ResponseWriter.Write(data)
}

func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}
