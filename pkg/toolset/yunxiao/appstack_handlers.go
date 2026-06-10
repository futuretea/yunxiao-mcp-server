package yunxiao

import (
	"context"
	"net/url"
)

func handleListApplications(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "pagination")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "nextToken")
	setOptionalInt(query, params, "page")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/apps:search"
	return c.GetJSON(ctx, path, query)
}

func handleGetApplication(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	appName, err := requiredString(params, "appName")
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, appstackAppPath(organizationID, appName), nil)
}
