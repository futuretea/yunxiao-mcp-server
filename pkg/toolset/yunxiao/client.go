package yunxiao

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// GetJSON sends a GET request and returns a pretty-printed JSON response.
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values) (string, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", friendlyAPIError(err)
	}
	return prettyJSON(resp.Body), nil
}

// GetJSONWithMetadata sends a GET request and includes pagination metadata in the response.
func (c *Client) GetJSONWithMetadata(ctx context.Context, path string, query url.Values) (string, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", friendlyAPIError(err)
	}
	return prettyResponseJSON(resp), nil
}

// PostJSONWithMetadata sends a POST request and includes pagination metadata in the response.
func (c *Client) PostJSONWithMetadata(ctx context.Context, path string, body any) (string, error) {
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return "", friendlyAPIError(err)
	}
	return prettyResponseJSON(resp), nil
}

// PutJSONWithMetadata sends a PUT request and includes pagination metadata in the response.
func (c *Client) PutJSONWithMetadata(ctx context.Context, path string, body any) (string, error) {
	resp, err := c.Request(ctx, http.MethodPut, path, nil, body)
	if err != nil {
		return "", friendlyAPIError(err)
	}
	return prettyResponseJSON(resp), nil
}

// Request sends an authenticated Yunxiao OpenAPI request.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body any) (*Response, error) {
	accessToken := c.accessToken
	if scopedToken := AccessTokenFromContext(ctx); scopedToken != "" {
		accessToken = scopedToken
	}
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required; set --access-token, YUNXIAO_MCP_ACCESS_TOKEN, or pass x-yunxiao-token/yunxiao_access_token on the HTTP request")
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
	req.Header.Set(AccessTokenHeader, accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

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
	return &Response{
		Body:       json.RawMessage(responseBody),
		Pagination: parsePagination(resp.Header),
		NextToken:  resp.Header.Get("x-next-token"),
		RequestID:  resp.Header.Get("x-request-id"),
	}, nil
}

func (c *Client) resolveURL(path string, query url.Values) string {
	u := *c.baseURL
	escapedPath := joinEscapedPath(c.baseURL.EscapedPath(), path)
	decodedPath, err := url.PathUnescape(escapedPath)
	if err == nil {
		u.Path = decodedPath
		u.RawPath = escapedPath
	} else {
		u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
		u.RawPath = ""
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func joinEscapedPath(basePath, path string) string {
	basePath = strings.TrimRight(basePath, "/")
	path = strings.TrimLeft(path, "/")
	if basePath == "" {
		return "/" + path
	}
	if path == "" {
		return basePath
	}
	return basePath + "/" + path
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

func prettyResponseJSON(resp *Response) string {
	var data any
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		data = string(resp.Body)
	}

	payload := map[string]any{"data": data}
	if resp.Pagination != nil {
		payload["pagination"] = resp.Pagination
	}
	if resp.NextToken != "" {
		payload["nextToken"] = resp.NextToken
	}
	if resp.RequestID != "" {
		payload["requestId"] = resp.RequestID
	}
	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return prettyJSON(resp.Body)
	}
	return string(formatted)
}

func parsePagination(header http.Header) *Pagination {
	pagination := &Pagination{
		Page:       parseHeaderInt(header, "x-page"),
		PerPage:    parseHeaderInt(header, "x-per-page"),
		PrevPage:   parseHeaderInt(header, "x-prev-page"),
		NextPage:   parseHeaderInt(header, "x-next-page"),
		Total:      parseHeaderInt(header, "x-total"),
		TotalPages: parseHeaderInt(header, "x-total-pages"),
	}
	if pagination.Page == 0 &&
		pagination.PerPage == 0 &&
		pagination.PrevPage == 0 &&
		pagination.NextPage == 0 &&
		pagination.Total == 0 &&
		pagination.TotalPages == 0 {
		return nil
	}
	return pagination
}

func parseHeaderInt(header http.Header, key string) int {
	value := header.Get(key)
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func EncodeRepositoryID(repositoryID string) string {
	repositoryID = strings.TrimSpace(repositoryID)
	if repositoryID == "" {
		return ""
	}
	if strings.Contains(repositoryID, "%2F") || strings.Contains(repositoryID, "%2f") {
		return repositoryID
	}
	if !strings.Contains(repositoryID, "/") {
		return url.PathEscape(repositoryID)
	}

	parts := strings.SplitN(repositoryID, "/", 2)
	return url.PathEscape(parts[0]) + "%2F" + url.PathEscape(parts[1])
}

func encodePathValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "%2F") || strings.Contains(value, "%2f") {
		return value
	}
	return url.PathEscape(value)
}

func encodeFilePath(filePath string) string {
	return encodePathValue(strings.TrimPrefix(strings.TrimSpace(filePath), "/"))
}
