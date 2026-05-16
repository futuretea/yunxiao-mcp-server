package yunxiao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetJSON sends a GET request and returns a pretty-printed JSON response.
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values) (string, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", WrapError(friendlyAPIError(err))
	}
	return prettyJSON(resp.Body), nil
}

// GetJSONWithMetadata sends a GET request and includes pagination metadata in the response.
func (c *Client) GetJSONWithMetadata(ctx context.Context, path string, query url.Values) (string, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", WrapError(friendlyAPIError(err))
	}
	return prettyResponseJSON(resp), nil
}

// PostJSONWithMetadata sends a POST request and includes pagination metadata in the response.
func (c *Client) PostJSONWithMetadata(ctx context.Context, path string, body any) (string, error) {
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return "", WrapError(friendlyAPIError(err))
	}
	return prettyResponseJSON(resp), nil
}

// PutJSONWithMetadata sends a PUT request and includes pagination metadata in the response.
func (c *Client) PutJSONWithMetadata(ctx context.Context, path string, body any) (string, error) {
	resp, err := c.Request(ctx, http.MethodPut, path, nil, body)
	if err != nil {
		return "", WrapError(friendlyAPIError(err))
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
