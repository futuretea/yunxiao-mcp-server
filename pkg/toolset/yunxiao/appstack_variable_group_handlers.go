package yunxiao

import (
	"context"
	"net/url"
)

func handleGetEnvVariableGroups(ctx context.Context, client any, params map[string]any) (string, error) {
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

	path := appstackAppPath(organizationID, appName) + "/envs/" + url.PathEscape(envName) + "/variableGroups"
	return c.GetJSON(ctx, path, nil)
}

func handleGetVariableGroup(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, variableGroupName, err := requiredAppVariableGroup(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackVariableGroupPath(organizationID, appName, variableGroupName)
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppVariableGroups(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/variableGroups"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppVariableGroupsRevision(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/variableGroups:revision"
	return c.GetJSON(ctx, path, nil)
}

func requiredAppVariableGroup(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	variableGroupName, err := requiredString(params, "variableGroupName")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, variableGroupName, nil
}

func appstackVariableGroupPath(organizationID, appName, variableGroupName string) string {
	return appstackAppPath(organizationID, appName) + "/variableGroup/" + url.PathEscape(variableGroupName)
}
