package http

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	mcpserver "github.com/futuretea/yunxiao-mcp-server/pkg/server/mcp"
)

func newTestMCPServer(t *testing.T, accessToken string) *mcpserver.Server {
	return newTestMCPServerWithBaseURL(t, accessToken, config.DefaultBaseURL)
}

func newTestMCPServerWithBaseURL(t *testing.T, accessToken, baseURL string) *mcpserver.Server {
	t.Helper()

	server, err := mcpserver.NewServer(mcpserver.Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               baseURL,
		AccessToken:           accessToken,
		LogLevel:              "info",
		RequestTimeoutSeconds: 30,
		ReadOnly:              true,
	}})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	return server
}

func TestHealthEndpoint(t *testing.T) {
	handler := NewHandler(newTestMCPServer(t, "token"), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.String() != "healthy" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestHealthEndpointDoesNotRequireStartupAccessToken(t *testing.T) {
	handler := NewHandler(newTestMCPServer(t, ""), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.String() != "healthy" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestHealthEndpointReturnsUnhealthy(t *testing.T) {
	handler := NewHandler(mcpserver.NewTestServer(nil, nil), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if rec.Body.String() != "unhealthy" {
		t.Fatalf("body = %q, want unhealthy", rec.Body.String())
	}
}

func TestStreamableMCPEndpointIsMounted(t *testing.T) {
	handler := NewHandler(newTestMCPServer(t, "token"), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, MCPEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatal("/mcp endpoint is not mounted")
	}
}

func TestSSEMessageEndpointIsMounted(t *testing.T) {
	handler := NewHandler(newTestMCPServer(t, "token"), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, SSEMessageEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatal("/message endpoint is not mounted")
	}
}

func TestSSEEndpointCarriesQueryTokenToMessageEndpoint(t *testing.T) {
	testServer := httptest.NewServer(NewHandler(newTestMCPServer(t, ""), &http.Server{}, ""))
	defer testServer.Close()

	resp, err := testServer.Client().Get(testServer.URL + SSEEndpoint + "?yunxiao_access_token=query-token")
	if err != nil {
		t.Fatalf("GET /sse: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	var event strings.Builder
	for !strings.Contains(event.String(), "data: ") {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read SSE endpoint event: %v", err)
		}
		event.WriteString(line)
	}

	if !strings.Contains(event.String(), "yunxiao_access_token=query-token") {
		t.Fatalf("SSE endpoint event = %q", event.String())
	}
}

func TestStreamableMCPUsesRequestAccessTokenForToolCall(t *testing.T) {
	tokenCh := make(chan string, 1)
	pathCh := make(chan string, 1)
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/platform/organizations" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"org-1","name":"Test Org"}]`))
			return
		}
		tokenCh <- r.Header.Get("x-yunxiao-token")
		pathCh <- r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"user-1"}`))
	}))
	defer apiServer.Close()

	testServer := httptest.NewServer(NewHandler(
		newTestMCPServerWithBaseURL(t, "default-token", apiServer.URL),
		&http.Server{},
		"",
	))
	defer testServer.Close()

	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "get_current_user",
		},
	})
	if err != nil {
		t.Fatalf("marshal JSON-RPC request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, testServer.URL+MCPEndpoint, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-yunxiao-token", "request-token")

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("POST /mcp: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d body = %s", resp.StatusCode, body)
	}

	select {
	case got := <-tokenCh:
		if got != "request-token" {
			t.Fatalf("x-yunxiao-token = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Yunxiao API was not called")
	}

	select {
	case got := <-pathCh:
		if got != "/oapi/v1/platform/users:me" {
			t.Fatalf("path = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Yunxiao API path was not captured")
	}

	if !strings.Contains(string(body), "user-1") {
		t.Fatalf("body = %s", body)
	}
}

func TestSSEMessageUsesQueryAccessTokenForToolCall(t *testing.T) {
	tokenCh := make(chan string, 1)
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oapi/v1/platform/organizations" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"org-1","name":"Test Org"}]`))
			return
		}
		tokenCh <- r.Header.Get("x-yunxiao-token")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"user-1"}`))
	}))
	defer apiServer.Close()

	testServer := httptest.NewServer(NewHandler(
		newTestMCPServerWithBaseURL(t, "default-token", apiServer.URL),
		&http.Server{},
		"",
	))
	defer testServer.Close()

	sseResp, err := testServer.Client().Get(testServer.URL + SSEEndpoint + "?yunxiao_access_token=query-token")
	if err != nil {
		t.Fatalf("GET /sse: %v", err)
	}
	defer sseResp.Body.Close()
	if sseResp.StatusCode != http.StatusOK {
		t.Fatalf("SSE status = %d", sseResp.StatusCode)
	}

	messageEndpoint := readSSEDataLine(t, sseResp.Body)
	if !strings.Contains(messageEndpoint, "yunxiao_access_token=query-token") {
		t.Fatalf("message endpoint = %q", messageEndpoint)
	}

	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "get_current_user",
		},
	})
	if err != nil {
		t.Fatalf("marshal JSON-RPC request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, testServer.URL+messageEndpoint, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	messageResp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("POST /message: %v", err)
	}
	defer messageResp.Body.Close()
	if messageResp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(messageResp.Body)
		t.Fatalf("message status = %d body = %s", messageResp.StatusCode, body)
	}

	select {
	case got := <-tokenCh:
		if got != "query-token" {
			t.Fatalf("x-yunxiao-token = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Yunxiao API was not called")
	}
}

func readSSEDataLine(t *testing.T, body io.Reader) string {
	t.Helper()

	reader := bufio.NewReader(body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read SSE event: %v", err)
		}
		if data, ok := strings.CutPrefix(line, "data: "); ok {
			return strings.TrimSpace(data)
		}
	}
}

func TestHandlerShutdown(t *testing.T) {
	httpServer := &http.Server{}
	handler := NewHandler(newTestMCPServer(t, "token"), httpServer, "")

	if err := handler.Shutdown(t.Context()); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

func TestHandlerShutdownCancelsActiveStreamableHTTPGet(t *testing.T) {
	testServer := httptest.NewUnstartedServer(nil)
	handler := NewHandler(newTestMCPServer(t, "token"), testServer.Config, "")
	testServer.Config.Handler = handler
	testServer.Start()
	defer testServer.Close()

	resp, err := testServer.Client().Get(testServer.URL + MCPEndpoint)
	if err != nil {
		t.Fatalf("GET /mcp: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}

	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		done <- handler.Shutdown(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Shutdown() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Shutdown() timed out with active streamable HTTP GET")
	}
}

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
	lrw.Write([]byte("body"))
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
	listener.Close()

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
	defer listener.Close()
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
