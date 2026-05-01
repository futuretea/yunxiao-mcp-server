package yunxiao

import (
	"context"
	"net/url"
)

func handleGetLatestOrchestration(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	envName, err := requiredString(params, "envName")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/envs/" + url.PathEscape(envName) + "/orchestration:latestAvailable"
	return c.GetJSON(ctx, path, nil)
}

func handleListAppOrchestration(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/orchestrations"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppOrchestration(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, sn, err := requiredAppOrchestration(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "tagName")
	setOptionalString(query, params, "sha")

	path := appstackOrchestrationPath(organizationID, appName, sn)
	return c.GetJSON(ctx, path, query)
}

func requiredAppOrchestration(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	sn, err := requiredString(params, "sn")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, sn, nil
}

func appstackOrchestrationPath(organizationID, appName, sn string) string {
	return appstackAppPath(organizationID, appName) + "/orchestrations/" + url.PathEscape(sn)
}
