package http

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestMiddleware logs HTTP requests without query strings, avoiding accidental token exposure.
func RequestMiddleware(next http.Handler) http.Handler {
	return RequestMiddlewareWithLogger(next, log.Logger)
}

// RequestMiddlewareWithLogger logs HTTP requests without query strings using the provided logger.
func RequestMiddlewareWithLogger(next http.Handler, logger zerolog.Logger) http.Handler {
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

		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lrw.statusCode).
			Dur("duration", time.Since(start)).
			Msg("handled HTTP request")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	mu            sync.Mutex
	statusCode    int
	headerWritten bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.mu.Lock()
	if lrw.headerWritten {
		lrw.mu.Unlock()
		return
	}
	lrw.statusCode = code
	lrw.headerWritten = true
	lrw.mu.Unlock()
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	lrw.mu.Lock()
	if !lrw.headerWritten {
		lrw.mu.Unlock()
		lrw.WriteHeader(http.StatusOK)
	} else {
		lrw.mu.Unlock()
	}
	return lrw.ResponseWriter.Write(data)
}

func (lrw *loggingResponseWriter) Flush() {
	lrw.mu.Lock()
	defer lrw.mu.Unlock()
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

func (lrw *loggingResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	if rf, ok := lrw.ResponseWriter.(io.ReaderFrom); ok {
		lrw.mu.Lock()
		if !lrw.headerWritten {
			lrw.mu.Unlock()
			lrw.WriteHeader(http.StatusOK)
		} else {
			lrw.mu.Unlock()
		}
		return rf.ReadFrom(r)
	}
	return io.Copy(lrw, r)
}
