package mcp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	yunxiaoToolset "github.com/futuretea/yunxiao-mcp-server/pkg/toolset/yunxiao"
)

func TestServeSSEReturnsWorkingSSEServer(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.AccessToken = "token"
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	sseServer := s.ServeSSE("", &http.Server{})
	if sseServer == nil {
		t.Fatal("ServeSSE() returned nil")
	}

	ts := httptest.NewServer(sseServer.SSEHandler())
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatalf("GET /sse: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
}

func TestServeSSEUsesBaseURL(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.AccessToken = "token"
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	sseServer := s.ServeSSE("http://localhost:8080", &http.Server{})
	if sseServer == nil {
		t.Fatal("ServeSSE() returned nil")
	}

	ep, err := sseServer.CompleteSseEndpoint()
	if err != nil {
		t.Fatalf("CompleteSseEndpoint() error = %v", err)
	}
	if ep != "http://localhost:8080/sse" {
		t.Fatalf("SSE endpoint = %q, want http://localhost:8080/sse", ep)
	}
}

func TestServeStreamableHTTPReturnsWorkingHandler(t *testing.T) {
	s, err := NewServer(Configuration{StaticConfig: newTestConfig(func(c *config.StaticConfig) {
		c.AccessToken = "token"
	})})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	streamable := s.ServeStreamableHTTP(&http.Server{})
	if streamable == nil {
		t.Fatal("ServeStreamableHTTP() returned nil")
	}

	ts := httptest.NewServer(streamable)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatalf("GET /mcp: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestCloseIsNoOp(t *testing.T) {
	s := newTestServer(nil, nil)
	s.Close() // should not panic
}

func TestNewTestServerCreatesUnhealthyServer(t *testing.T) {
	s := NewTestServer(nil, nil)
	if s.IsHealthy() {
		t.Fatal("NewTestServer(nil, nil) should not be healthy")
	}
}

func TestNewTestServerCreatesHealthyServer(t *testing.T) {
	s := NewTestServer(&yunxiaoToolset.Client{}, []string{"tool-1"})
	if !s.IsHealthy() {
		t.Fatal("NewTestServer with client and tools should be healthy")
	}
}
