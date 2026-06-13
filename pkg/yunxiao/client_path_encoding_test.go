package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEncodeRepositoryID(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "numeric", in: "2813489", want: "2813489"},
		{name: "path", in: "org/Demo Repo", want: "org%2FDemo%20Repo"},
		{name: "already encoded", in: "org%2FDemoRepo", want: "org%2FDemoRepo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeRepositoryID(tt.in); got != tt.want {
				t.Fatalf("EncodeRepositoryID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestJoinEscapedPath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{"both non-empty", "/oapi/v1", "platform/users", "/oapi/v1/platform/users"},
		{"empty base", "", "platform/users", "/platform/users"},
		{"empty path", "/oapi/v1", "", "/oapi/v1"},
		{"both empty", "", "", "/"},
		{"trailing slash", "/oapi/v1/", "/platform/users", "/oapi/v1/platform/users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinEscapedPath(tt.basePath, tt.path); got != tt.want {
				t.Fatalf("joinEscapedPath(%q, %q) = %q, want %q", tt.basePath, tt.path, got, tt.want)
			}
		})
	}
}

func TestEncodePathValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"already encoded", "org%2Frepo", "org%2Frepo"},
		{"normal", "org/repo", "org%2Frepo"},
		{"whitespace", "  org/repo  ", "org%2Frepo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodePathValue(tt.in); got != tt.want {
				t.Fatalf("EncodePathValue(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEncodeFilePath(t *testing.T) {
	if got := EncodeFilePath(" /src/main.go "); got != "src%2Fmain.go" {
		t.Fatalf("EncodeFilePath() = %q, want src%%2Fmain.go", got)
	}
}

func TestParseHeaderInt(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{"empty", "", 0},
		{"valid", "42", 42},
		{"invalid", "abc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.value != "" {
				h.Set("x-test", tt.value)
			}
			if got := parseHeaderInt(h, "x-test"); got != tt.want {
				t.Fatalf("parseHeaderInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEncodeRepositoryIDEmpty(t *testing.T) {
	if got := EncodeRepositoryID(""); got != "" {
		t.Fatalf("EncodeRepositoryID(\"\") = %q, want empty", got)
	}
}

func TestResolveURLWithInvalidEscapeSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetJSON(context.Background(), "/bad%ZZ", nil)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if result != "{}" {
		t.Fatalf("result = %q, want {}", result)
	}
}

func TestClientPostJSONWithMetadataReturnsMarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "token-1", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.PostJSONWithMetadata(context.Background(), "/platform/users:me", map[string]any{"key": make(chan int)})
	if err == nil {
		t.Fatal("PostJSONWithMetadata() expected marshal error")
	}
}

type errorReader struct{}

func (e errorReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}
func (e errorReader) Close() error { return nil }

type errorBodyTransport struct{}

func (t errorBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       errorReader{},
		Header:     http.Header{},
		Request:    req,
	}, nil
}

type failingTransport struct {
	err error
}

func (t *failingTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, t.err
}

func TestClientRequestReturnsReadError(t *testing.T) {
	client, err := NewClient("https://example.com", "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: errorBodyTransport{}}

	_, err = client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("Request() expected read error")
	}
	if !strings.Contains(err.Error(), "read response body") {
		t.Fatalf("error = %v", err)
	}
}

type blockingTransport struct {
	blockUntil context.Context
}

func (t *blockingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	select {
	case <-t.blockUntil.Done():
		return nil, t.blockUntil.Err()
	case <-req.Context().Done():
		return nil, req.Context().Err()
	}
}

func TestClientRequestRespectsContextCancellation(t *testing.T) {
	client, err := NewClient("https://example.com", "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	blockCtx, blockCancel := context.WithCancel(context.Background())
	defer blockCancel()
	client.httpClient = &http.Client{Transport: &blockingTransport{blockUntil: blockCtx}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.Request(ctx, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("Request() expected context cancellation error")
	}
	if !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("error = %v", err)
	}
}
