package yunxiao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func getClient(client any) (*Client, error) {
	c, ok := client.(*Client)
	if !ok || c == nil {
		return nil, fmt.Errorf("Yunxiao client is not configured")
	}
	return c, nil
}

func handleGetCurrentUser(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	return c.GetJSON(ctx, "/platform/users:me", nil)
}

func handleGetCurrentOrganizationInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	return handleGetCurrentUser(ctx, client, params)
}

func handleListOrganizations(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	return c.GetJSON(ctx, "/platform/organizations", query)
}

func handleGetUserOrganizations(ctx context.Context, client any, params map[string]any) (string, error) {
	return handleListOrganizations(ctx, client, params)
}

func handleGetOrganization(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	id, _ := params["id"].(string)
	if id == "" {
		return "", fmt.Errorf("id is required")
	}
	return c.GetJSON(ctx, "/platform/organizations/"+url.PathEscape(id), nil)
}

func setOptionalInt(query url.Values, params map[string]any, key string) {
	switch value := params[key].(type) {
	case float64:
		query.Set(key, strconv.Itoa(int(value)))
	case int:
		query.Set(key, strconv.Itoa(value))
	case int64:
		query.Set(key, strconv.FormatInt(value, 10))
	case string:
		if value != "" {
			query.Set(key, value)
		}
	}
}
