package yunxiao

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/version"
)

const apiBasePath = "/oapi/v1"

const (
	// AccessTokenHeader is the request header used by Yunxiao OpenAPI.
	AccessTokenHeader = "x-yunxiao-token"
	// AccessTokenQueryParam is the SSE-compatible query parameter used by the Node reference server.
	AccessTokenQueryParam = "yunxiao_access_token"
)

type accessTokenContextKey struct{}

type clientOptions struct {
	insecureSkipTLSVerify bool
}

// ClientOption customizes the Yunxiao OpenAPI client.
type ClientOption func(*clientOptions)

// WithInsecureSkipTLSVerify configures whether the client skips server TLS verification.
func WithInsecureSkipTLSVerify(skip bool) ClientOption {
	return func(options *clientOptions) {
		options.insecureSkipTLSVerify = skip
	}
}

// Client is a minimal Yunxiao OpenAPI client.
type Client struct {
	baseURL      *url.URL
	accessToken  string
	httpClient   *http.Client
	userAgent    string
	DefaultOrgID string
}

// APIError includes response context from a failed Yunxiao API call.
type APIError struct {
	StatusCode int
	Method     string
	URL        string
	Body       string
}

// Response contains a Yunxiao response body and selected response metadata.
type Response struct {
	Body       json.RawMessage `json:"body"`
	Pagination *Pagination     `json:"pagination,omitempty"`
	NextToken  string          `json:"nextToken,omitempty"`
	RequestID  string          `json:"requestId,omitempty"`
}

// Pagination contains standard Yunxiao pagination headers when present.
type Pagination struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"perPage,omitempty"`
	PrevPage   int `json:"prevPage,omitempty"`
	NextPage   int `json:"nextPage,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

// WithAccessToken returns a context carrying a request-scoped Yunxiao access token.
func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return ctx
	}
	return context.WithValue(ctx, accessTokenContextKey{}, accessToken)
}

// AccessTokenFromContext returns the request-scoped Yunxiao access token, if present.
func AccessTokenFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if accessToken, ok := ctx.Value(accessTokenContextKey{}).(string); ok {
		return strings.TrimSpace(accessToken)
	}
	return ""
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("Yunxiao API error: %s %s returned status %d", e.Method, e.URL, e.StatusCode)
	}
	return fmt.Sprintf("Yunxiao API error: %s %s returned status %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}

// friendlyAPIError wraps an APIError with actionable guidance for LLM consumers.
// Non-API errors are returned unchanged.
func friendlyAPIError(err error) error {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return err
	}

	var suggestion string
	switch apiErr.StatusCode {
	case http.StatusUnauthorized:
		suggestion = "Authentication failed. Verify that your access token is valid and not expired."
	case http.StatusForbidden:
		suggestion = "Access denied. Your token may not have permission for this resource."
	case http.StatusNotFound:
		suggestion = "Resource not found. Verify that the project ID, work item ID, pipeline ID, or other identifiers are correct. Use search_projects, search_workitems, or list_pipelines to find valid IDs."
	case http.StatusBadRequest:
		suggestion = "Invalid request parameters. Check that required fields are present, IDs are correct, and enum values are valid. Use the corresponding get_*_context or list_* tools to discover valid values."
	case http.StatusTooManyRequests:
		suggestion = "Rate limit exceeded. Wait a moment before retrying."
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		suggestion = "Yunxiao service temporarily unavailable. Retry the request later."
	default:
		return err
	}

	return fmt.Errorf("%w\n\nSuggestion: %s", apiErr, suggestion)
}

// NewClient creates a Yunxiao OpenAPI client.
func NewClient(baseURL, accessToken string, timeout time.Duration, opts ...ClientOption) (*Client, error) {
	parsed, err := normalizeAPIBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	options := clientOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return &Client{
		baseURL:     parsed,
		accessToken: strings.TrimSpace(accessToken),
		httpClient:  newHTTPClient(timeout, options),
		userAgent:   fmt.Sprintf("modelcontextprotocol/%s/%s", version.BinaryName, version.Version),
	}, nil
}

func newHTTPClient(timeout time.Duration, options clientOptions) *http.Client {
	client := &http.Client{Timeout: timeout}
	if !options.insecureSkipTLSVerify {
		return client
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // Explicit opt-in for private/self-signed Yunxiao endpoints.
	client.Transport = transport
	return client
}

// ResolveDefaultOrgID fetches the user's organizations and, if exactly one exists,
// caches its ID as the default organization for automatic parameter filling.
func (c *Client) ResolveDefaultOrgID(ctx context.Context) error {
	resp, err := c.Request(ctx, http.MethodGet, "/platform/organizations", nil, nil)
	if err != nil {
		return fmt.Errorf("list organizations: %w", err)
	}

	// Try array format first
	var orgList []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(resp.Body, &orgList); err == nil && len(orgList) == 1 {
		c.DefaultOrgID = orgList[0].ID
		return nil
	}

	// Try { data: [...] } format
	var wrapped struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &wrapped); err == nil && len(wrapped.Data) == 1 {
		c.DefaultOrgID = wrapped.Data[0].ID
		return nil
	}

	return nil
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
