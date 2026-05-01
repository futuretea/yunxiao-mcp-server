package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	mcpserver "github.com/futuretea/yunxiao-mcp-server/pkg/server/mcp"
)

func newTestMCPServer(t *testing.T, accessToken string) *mcpserver.Server {
	t.Helper()

	server, err := mcpserver.NewServer(mcpserver.Configuration{StaticConfig: &config.StaticConfig{
		BaseURL:               config.DefaultBaseURL,
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

func TestHealthEndpointRequiresAccessToken(t *testing.T) {
	handler := NewHandler(newTestMCPServer(t, ""), &http.Server{}, "")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, HealthEndpoint, nil)

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.String() != "unhealthy" {
		t.Fatalf("body = %q", rec.Body.String())
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
