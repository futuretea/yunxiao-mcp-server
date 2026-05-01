package yunxiao

import (
	"context"
	"net/url"
)

func handleSearchReleases(ctx context.Context, client any, params map[string]any) (string, error) {
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
	setOptionalString(query, params, "nameKeyword")
	setOptionalString(query, params, "systemName")
	setOptionalStringArrayQuery(query, params, "states")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/releases:search"
	return c.GetJSON(ctx, path, query)
}
