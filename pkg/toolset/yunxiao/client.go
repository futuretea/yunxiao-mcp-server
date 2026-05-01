package yunxiao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
)

const apiBasePath = "/oapi/v1"

// Client is a minimal Yunxiao OpenAPI client.
type Client struct {
	baseURL     *url.URL
	accessToken string
	httpClient  *http.Client
	userAgent   string
}

// APIError includes response context from a failed Yunxiao API call.
type APIError struct {
	StatusCode int
	Method     string
	URL        string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("Yunxiao API error: %s %s returned status %d", e.Method, e.URL, e.StatusCode)
	}
	return fmt.Sprintf("Yunxiao API error: %s %s returned status %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}

// NewClient creates a Yunxiao OpenAPI client.
func NewClient(baseURL, accessToken string, timeout time.Duration) (*Client, error) {
	parsed, err := normalizeAPIBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:     parsed,
		accessToken: strings.TrimSpace(accessToken),
		httpClient:  &http.Client{Timeout: timeout},
		userAgent:   fmt.Sprintf("modelcontextprotocol/%s/%s", version.BinaryName, version.Version),
	}, nil
}

func normalizeAPIBaseURL(raw string) (*url.URL, error) {
	raw = strings.TrimRight(strings.TrimSpace(raw), "/")
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid base URL %q", raw)
	}
	if !strings.HasSuffix(parsed.Path, apiBasePath) {
		parsed.Path = strings.TrimRight(parsed.Path, "/") + apiBasePath
	}
	return parsed, nil
}

// GetJSON sends a GET request and returns a pretty-printed JSON response.
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values) (string, error) {
	raw, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", err
	}
	return prettyJSON(raw), nil
}

// Request sends an authenticated Yunxiao OpenAPI request.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body any) (json.RawMessage, error) {
	if c.accessToken == "" {
		return nil, fmt.Errorf("access token is required; set --access-token or YUNXIAO_MCP_ACCESS_TOKEN")
	}

	requestURL := c.resolveURL(path, query)

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("x-yunxiao-token", c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Method:     method,
			URL:        requestURL,
			Body:       string(responseBody),
		}
	}
	return json.RawMessage(responseBody), nil
}

func (c *Client) resolveURL(path string, query url.Values) string {
	u := *c.baseURL
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	u.RawQuery = query.Encode()
	return u.String()
}

func prettyJSON(raw json.RawMessage) string {
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return string(raw)
	}
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(raw)
	}
	return string(formatted)
}
