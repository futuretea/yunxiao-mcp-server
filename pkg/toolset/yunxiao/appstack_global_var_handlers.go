package yunxiao

import (
	"context"
	"net/http"
	"net/url"
)

func handleGetGlobalVar(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, name, err := requiredGlobalVar(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "revisionSha")

	return c.GetJSON(ctx, appstackGlobalVarPath(organizationID, name), query)
}

func handleListGlobalVars(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := appstackDefaultPageQuery(params)
	body := map[string]any{}
	setOptionalStringBody(body, params, "search")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/globalVars:search"
	resp, err := c.Request(ctx, http.MethodPost, path, query, body)
	if err != nil {
		return "", err
	}
	return prettyResponseJSON(resp), nil
}

func requiredGlobalVar(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	name, err := requiredString(params, "name")
	if err != nil {
		return "", "", err
	}
	return organizationID, name, nil
}

func appstackGlobalVarPath(organizationID, name string) string {
	return "/appstack/organizations/" + url.PathEscape(organizationID) + "/globalVars/" + url.PathEscape(name)
}
