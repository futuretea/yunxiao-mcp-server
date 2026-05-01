package yunxiao

import (
	"context"
	"net/url"
)

func handleListSystems(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/systems"
	return c.GetJSON(ctx, path, appstackDefaultPageQuery(params))
}

func handleListAttachedApps(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, err := requiredOrganizationAndSystem(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackSystemPath(organizationID, systemName) + "/apps"
	return c.GetJSON(ctx, path, appstackDefaultPageQuery(params))
}

func appstackDefaultPageQuery(params map[string]any) url.Values {
	query := url.Values{}
	query.Set("current", "1")
	query.Set("pageSize", "10")
	setOptionalInt(query, params, "current")
	setOptionalInt(query, params, "pageSize")
	return query
}
