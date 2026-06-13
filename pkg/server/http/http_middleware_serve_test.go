package http

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/rs/zerolog"
)

func TestMiddlewareDoesNotLogQueryString(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.DebugLevel)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?yunxiao_access_token=secret", nil)
	handler := RequestMiddlewareWithLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}), logger)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	logOutput := buf.String()
	if !strings.Contains(logOutput, `"path":"/test"`) {
		t.Fatalf("log output missing expected path: %s", logOutput)
	}
	if strings.Contains(logOutput, "secret") {
		t.Fatalf("log output leaked query string: %s", logOutput)
	}
}

func TestMiddlewareSkipsLoggingForHealthEndpoint(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.DebugLevel)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)
	called := false
	handler := RequestMiddlewareWithLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}), logger)

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("handler was not called for health endpoint")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no log output for health endpoint, got: %s", buf.String())
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

type flushableRecorder struct {
	*httptest.ResponseRecorder
	flushCalled bool
}

func (f *flushableRecorder) Flush() {
	f.flushCalled = true
	f.ResponseRecorder.Flush()
}

func TestLoggingResponseWriterFlush(t *testing.T) {
	rec := &flushableRecorder{ResponseRecorder: httptest.NewRecorder()}
	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	lrw.Flush()
	if !rec.flushCalled {
		t.Fatal("Flush() did not delegate to underlying Flusher")
	}
}

func TestLoggingResponseWriterHijack(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	// ResponseRecorder does not implement Hijacker
	_, _, err := lrw.Hijack()
	if err != http.ErrNotSupported {
		t.Fatalf("Hijack() error = %v, want ErrNotSupported", err)
	}

	dummyConn := &net.TCPConn{}
	dummyBuf := bufio.NewReadWriter(bufio.NewReader(strings.NewReader("")), bufio.NewWriter(io.Discard))
	hijackRec := &hijackableRecorder{ResponseRecorder: rec, conn: dummyConn, buf: dummyBuf}
	lrw2 := &loggingResponseWriter{ResponseWriter: hijackRec, statusCode: http.StatusOK}
	conn, buf, err := lrw2.Hijack()
	if err != nil {
		t.Fatalf("Hijack() error = %v", err)
	}
	if conn != dummyConn {
		t.Fatal("Hijack() returned wrong conn")
	}
	if buf != dummyBuf {
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

type readerFromRecorder struct {
	*httptest.ResponseRecorder
	readFromCalled bool
}

func (r *readerFromRecorder) ReadFrom(src io.Reader) (int64, error) {
	r.readFromCalled = true
	return io.Copy(r.ResponseRecorder, src)
}

func TestLoggingResponseWriterReadFromPreservesOptimization(t *testing.T) {
	t.Run("delegates to io.ReaderFrom when available", func(t *testing.T) {
		rec := &readerFromRecorder{ResponseRecorder: httptest.NewRecorder()}
		lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

		src := strings.NewReader("body")
		n, err := lrw.ReadFrom(src)
		if err != nil {
			t.Fatalf("ReadFrom() error = %v", err)
		}
		if n != 4 {
			t.Fatalf("ReadFrom() n = %d, want 4", n)
		}
		if !rec.readFromCalled {
			t.Fatal("ReadFrom() did not delegate to underlying io.ReaderFrom")
		}
		if rec.Body.String() != "body" {
			t.Fatalf("body = %q, want body", rec.Body.String())
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("falls back to io.Copy when io.ReaderFrom is absent", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

		src := strings.NewReader("fallback")
		n, err := lrw.ReadFrom(src)
		if err != nil {
			t.Fatalf("ReadFrom() error = %v", err)
		}
		if n != 8 {
			t.Fatalf("ReadFrom() n = %d, want 8", n)
		}
		if rec.Body.String() != "fallback" {
			t.Fatalf("body = %q, want fallback", rec.Body.String())
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

type readyListener struct {
	net.Listener
	ready chan struct{}
	once  sync.Once
}

func (l *readyListener) Accept() (net.Conn, error) {
	l.once.Do(func() { close(l.ready) })
	return l.Listener.Accept()
}

func TestServeStartsAndShutsDown(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"default-org"}]`))
	}))
	defer ts.Close()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Logf("close listener: %v", err)
		}
	}()

	addr := fmt.Sprintf("http://%s%s", listener.Addr().String(), HealthEndpoint)
	rl := &readyListener{Listener: listener, ready: make(chan struct{})}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- ServeListener(ctx, newTestMCPServerWithBaseURL(t, "token", ts.URL), &config.StaticConfig{
			LogLevel:              "info",
			RequestTimeoutSeconds: 30,
		}, rl)
	}()

	select {
	case <-rl.ready:
	case err := <-serveErr:
		t.Fatalf("Serve() exited before ready: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("server did not become ready within 5s")
	}

	resp, err := http.Get(addr)
	if err != nil {
		t.Fatalf("GET %s: %v", addr, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	_ = resp.Body.Close()

	cancel()
	select {
	case err := <-serveErr:
		if err != nil {
			t.Fatalf("Serve() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve() did not return after shutdown")
	}
}

type errorListener struct {
	addr net.Addr
	err  error
}

func (l *errorListener) Accept() (net.Conn, error) { return nil, l.err }
func (l *errorListener) Close() error              { return nil }
func (l *errorListener) Addr() net.Addr            { return l.addr }

func TestServeReturnsListenError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := Serve(ctx, newTestMCPServer(t, "token"), &config.StaticConfig{
		Port:                  -1,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
	})
	if err == nil {
		t.Fatal("Serve() expected error for invalid port")
	}
}

func TestServeListenerReturnsServeError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	listener := &errorListener{
		addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345},
		err:  errors.New("mock listener error"),
	}

	err := ServeListener(ctx, newTestMCPServer(t, "token"), &config.StaticConfig{
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
	}, listener)
	if err == nil {
		t.Fatal("ServeListener() expected error")
	}
	if !strings.Contains(err.Error(), "mock listener error") {
		t.Fatalf("ServeListener() error = %v, want mock listener error", err)
	}
}
