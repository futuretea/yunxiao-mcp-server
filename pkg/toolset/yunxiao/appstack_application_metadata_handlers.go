package yunxiao

import (
	"context"
	"net/url"
)

func handleSearchAppTemplates(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := appstackOptionalPageQuery(params)
	setOptionalString(query, params, "displayNameKeyword")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/appTemplates:search"
	return c.GetJSON(ctx, path, query)
}

func handleListEnvironments(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/envs"
	return c.GetJSON(ctx, path, appstackOptionalPageQuery(params))
}

func handleGetEnvironment(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, envName, err := requiredAppEnvironment(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/envs/" + url.PathEscape(envName)
	return c.GetJSON(ctx, path, nil)
}

func handleListApplicationMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/members"
	return c.GetJSON(ctx, path, appstackDefaultPageQuery(params))
}

func handleListApplicationSources(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/sources"
	return c.GetJSON(ctx, path, appstackOptionalPageQuery(params))
}

func requiredAppEnvironment(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	envName, err := requiredString(params, "envName")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, envName, nil
}

func appstackOptionalPageQuery(params map[string]any) url.Values {
	query := url.Values{}
	setOptionalString(query, params, "pagination")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "nextToken")
	setOptionalInt(query, params, "page")
	return query
}
