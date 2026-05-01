package yunxiao

import (
	"context"
	"net/url"
)

func handleListAppReleaseWorkflows(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/releaseWorkflows"
	return c.GetJSON(ctx, path, nil)
}

func handleListAppReleaseWorkflowBriefs(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/releaseWorkflowBriefs"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppReleaseWorkflowStage(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, err := requiredAppReleaseStage(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn) + "/releaseStage/" + url.PathEscape(releaseStageSn)
	return c.GetJSON(ctx, path, nil)
}

func handleListAppReleaseStageBriefs(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, err := requiredAppReleaseWorkflow(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn) + "/releaseStageBriefs"
	return c.GetJSON(ctx, path, nil)
}

func requiredAppReleaseWorkflow(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	releaseWorkflowSn, err := requiredString(params, "releaseWorkflowSn")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, releaseWorkflowSn, nil
}

func requiredAppReleaseStage(params map[string]any) (string, string, string, string, error) {
	organizationID, appName, releaseWorkflowSn, err := requiredAppReleaseWorkflow(params)
	if err != nil {
		return "", "", "", "", err
	}
	releaseStageSn, err := requiredString(params, "releaseStageSn")
	if err != nil {
		return "", "", "", "", err
	}
	return organizationID, appName, releaseWorkflowSn, releaseStageSn, nil
}

func appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn string) string {
	return appstackAppPath(organizationID, appName) + "/releaseWorkflow/" + url.PathEscape(releaseWorkflowSn)
}
