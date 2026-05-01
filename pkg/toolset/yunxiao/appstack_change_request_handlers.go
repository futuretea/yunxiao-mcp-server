package yunxiao

import (
	"context"
	"net/url"
)

func handleGetAppStackChangeRequestAuditItems(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, sn, err := requiredAppStackChangeRequest(params)
	if err != nil {
		return "", err
	}
	refType, err := requiredString(params, "refType")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("refType", refType)

	path := appstackChangeRequestPath(organizationID, appName, sn) + "/auditItems"
	return c.GetJSON(ctx, path, query)
}

func handleListAppStackChangeRequestExecutions(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, sn, err := requiredAppStackChangeRequest(params)
	if err != nil {
		return "", err
	}
	releaseWorkflowSn, err := requiredString(params, "releaseWorkflowSn")
	if err != nil {
		return "", err
	}
	releaseStageSn, err := requiredString(params, "releaseStageSn")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("releaseWorkflowSn", releaseWorkflowSn)
	query.Set("releaseStageSn", releaseStageSn)
	setOptionalInt(query, params, "perPage")
	setOptionalInt(query, params, "page")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")

	path := appstackChangeRequestPath(organizationID, appName, sn) + "/executions"
	return c.GetJSON(ctx, path, query)
}

func handleListAppStackChangeRequestWorkItems(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, sn, err := requiredAppStackChangeRequest(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackChangeRequestPath(organizationID, appName, sn) + "/workItems"
	return c.GetJSON(ctx, path, nil)
}

func requiredAppStackChangeRequest(params map[string]any) (string, string, string, error) {
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

func appstackChangeRequestPath(organizationID, appName, sn string) string {
	return appstackAppPath(organizationID, appName) + "/changeRequests/" + url.PathEscape(sn)
}
