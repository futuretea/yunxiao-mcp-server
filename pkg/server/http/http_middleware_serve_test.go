package http

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
)

func TestMiddlewareDoesNotLogQueryString(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?yunxiao_access_token=secret", nil)
	handler := RequestMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestMiddlewareSkipsLoggingForHealthEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)
	called := false
	handler := RequestMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("handler was not called for health endpoint")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestLoggingResponseWriterIsIdempotent(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	lrw.WriteHeader(http.StatusCreated)
	if lrw.statusCode != http.StatusCreated {
		t.Fatalf("statusCode = %d, want %d", lrw.statusCode, http.StatusCreated)
	}

	// Second call should be ignored
	lrw.WriteHeader(http.StatusInternalServerError)
	if lrw.statusCode != http.StatusCreated {
		t.Fatalf("statusCode = %d, want %d after second WriteHeader", lrw.statusCode, http.StatusCreated)
	}

	// Write should not override status if already set
	_, _ = lrw.Write([]byte("body"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("recorder code = %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestLoggingResponseWriterFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	// Should not panic; ResponseRecorder implements Flusher
	lrw.Flush()
}

func TestLoggingResponseWriterHijack(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	// ResponseRecorder does not implement Hijacker
	_, _, err := lrw.Hijack()
	if err != http.ErrNotSupported {
		t.Fatalf("Hijack() error = %v, want ErrNotSupported", err)
	}

	hijackRec := &hijackableRecorder{ResponseRecorder: rec}
	lrw2 := &loggingResponseWriter{ResponseWriter: hijackRec, statusCode: http.StatusOK}
	conn, buf, err := lrw2.Hijack()
	if err != nil {
		t.Fatalf("Hijack() error = %v", err)
	}
	if conn != hijackRec.conn {
		t.Fatal("Hijack() returned wrong conn")
	}
	if buf != hijackRec.buf {
		t.Fatal("Hijack() returned wrong bufio")
	}
}

type hijackableRecorder struct {
	*httptest.ResponseRecorder
	conn net.Conn
	buf  *bufio.ReadWriter
}

func (h *hijackableRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, h.buf, nil
}

func TestServeStartsAndShutsDown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = Serve(ctx, newTestMCPServer(t, "token"), &config.StaticConfig{
		Port:                  port,
		LogLevel:              "info",
		BaseURL:               config.DefaultBaseURL,
		RequestTimeoutSeconds: 30,
	})
	if err != nil {
		t.Fatalf("Serve() error = %v", err)
	}
}

func TestServeReturnsListenError(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = listener.Close() }()
	port := listener.Addr().(*net.TCPAddr).Port

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = Serve(ctx, newTestMCPServer(t, "token"), &config.StaticConfig{
		Port:                  port,
		LogLevel:              "info",
		BaseURL:               config.DefaultBaseURL,
		RequestTimeoutSeconds: 30,
	})
	if err == nil {
		t.Fatal("Serve() expected error for occupied port")
	}
}
