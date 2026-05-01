package yunxiao

import (
	"context"
	"net/url"
)

func handleListSystemReleaseWorkflows(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, err := requiredOrganizationAndSystem(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackSystemPath(organizationID, systemName) + "/releaseWorkflows"
	return c.GetJSON(ctx, path, nil)
}

func handleGetRelease(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, sn, err := requiredSystemRelease(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, appstackSystemReleasePath(organizationID, systemName, sn), nil)
}

func handleListReleaseMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, sn, err := requiredSystemRelease(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackSystemReleasePath(organizationID, systemName, sn) + "/members"
	return c.GetJSON(ctx, path, nil)
}

func handleListReleaseProducts(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, sn, err := requiredSystemRelease(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackSystemReleasePath(organizationID, systemName, sn) + "/products"
	return c.GetJSON(ctx, path, nil)
}

func requiredOrganizationAndSystem(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	systemName, err := requiredString(params, "systemName")
	if err != nil {
		return "", "", err
	}
	return organizationID, systemName, nil
}

func requiredSystemRelease(params map[string]any) (string, string, string, error) {
	organizationID, systemName, err := requiredOrganizationAndSystem(params)
	if err != nil {
		return "", "", "", err
	}
	sn, err := requiredString(params, "sn")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, systemName, sn, nil
}

func appstackSystemPath(organizationID, systemName string) string {
	return "/appstack/organizations/" + url.PathEscape(organizationID) + "/systems/" + url.PathEscape(systemName)
}

func appstackSystemReleasePath(organizationID, systemName, sn string) string {
	return appstackSystemPath(organizationID, systemName) + "/releases/" + url.PathEscape(sn)
}
