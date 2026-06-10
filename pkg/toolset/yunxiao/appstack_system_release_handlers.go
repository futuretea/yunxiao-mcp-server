package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func handleExecuteSystemReleaseStage(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, systemName, err := requiredOrganizationAndSystem(params)
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
	executionJSON, err := requiredString(params, "execution")
	if err != nil {
		return "", err
	}
	var execution any
	if err := json.Unmarshal([]byte(executionJSON), &execution); err != nil {
		return "", fmt.Errorf("invalid execution JSON: %w", err)
	}

	path := appstackSystemPath(organizationID, systemName) + "/releaseWorkflows/" + url.PathEscape(releaseWorkflowSn) + "/releaseStages/" + url.PathEscape(releaseStageSn) + ":execute"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, execution)
	if err != nil {
		return "", err
	}
	return PrettyResponseJSON(resp), nil
}

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

func handleListAttachedChangeRequests(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, releaseSn, err := requiredSystemReleaseWithKey(params, "releaseSn")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackSystemReleasePath(organizationID, systemName, releaseSn) + "/changeRequests"
	return c.GetJSON(ctx, path, appstackDefaultPageQuery(params))
}

func handleListReleaseExecutions(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, sn, err := requiredSystemRelease(params)
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

	path := appstackSystemReleasePath(organizationID, systemName, sn) + "/executions"
	return c.GetJSON(ctx, path, query)
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
	return requiredSystemReleaseWithKey(params, "sn")
}

func requiredSystemReleaseWithKey(params map[string]any, releaseKey string) (string, string, string, error) {
	organizationID, systemName, err := requiredOrganizationAndSystem(params)
	if err != nil {
		return "", "", "", err
	}
	sn, err := requiredString(params, releaseKey)
	if err != nil {
		return "", "", "", err
	}
	return organizationID, systemName, sn, nil
}

func appstackSystemPath(organizationID, systemName string) string {
	return "/appstack/organizations/" + encodePathValue(organizationID) + "/systems/" + encodePathValue(systemName)
}

func appstackSystemReleasePath(organizationID, systemName, sn string) string {
	return appstackSystemPath(organizationID, systemName) + "/releases/" + encodePathValue(sn)
}
