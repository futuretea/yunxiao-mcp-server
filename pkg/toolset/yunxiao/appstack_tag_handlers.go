package yunxiao

import (
	"context"
	"net/http"
	"net/url"
)

func handleSearchAppTags(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := appstackDefaultPageQuery(params)
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")

	body := map[string]any{}
	setOptionalStringBody(body, params, "search")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/appTags:search"
	resp, err := c.Request(ctx, http.MethodPost, path, query, body)
	if err != nil {
		return "", err
	}
	return prettyResponseJSON(resp), nil
}
