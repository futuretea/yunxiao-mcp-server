package yunxiao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func getClient(client any) (*Client, error) {
	c, ok := client.(*Client)
	if !ok || c == nil {
		return nil, fmt.Errorf("Yunxiao client is not configured")
	}
	return c, nil
}

func handleGetCurrentUser(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	return c.GetJSON(ctx, "/platform/users:me", nil)
}

func handleGetCurrentOrganizationInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	return handleGetCurrentUser(ctx, client, params)
}

func handleListOrganizations(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	return c.GetJSON(ctx, "/platform/organizations", query)
}

func handleGetUserOrganizations(ctx context.Context, client any, params map[string]any) (string, error) {
	return handleListOrganizations(ctx, client, params)
}

func handleGetOrganization(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	id, _ := params["id"].(string)
	if id == "" {
		return "", fmt.Errorf("id is required")
	}
	return c.GetJSON(ctx, "/platform/organizations/"+url.PathEscape(id), nil)
}

func handleListRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")
	setOptionalBool(query, params, "archived")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetRepository(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	repositoryID, err := requiredString(params, "repositoryId")
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID)
	return c.GetJSON(ctx, path, nil)
}

func handleListBranches(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	repositoryID, err := requiredString(params, "repositoryId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/branches"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleListPipelines(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "createStartTime")
	setOptionalInt(query, params, "createEndTime")
	setOptionalInt(query, params, "executeStartTime")
	setOptionalInt(query, params, "executeEndTime")
	setOptionalString(query, params, "pipelineName")
	setOptionalString(query, params, "statusList")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetPipeline(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID)
	return c.GetJSON(ctx, path, nil)
}

func handleListPipelineRuns(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalInt(query, params, "startTime")
	setOptionalIntAs(query, params, "endTime", "endTme")
	setOptionalString(query, params, "status")
	setOptionalInt(query, params, "triggerMode")

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetPipelineRun(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	pipelineRunID, err := requiredString(params, "pipelineRunId")
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs/" + url.PathEscape(pipelineRunID)
	return c.GetJSON(ctx, path, nil)
}

func handleGetLatestPipelineRun(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs/latestPipelineRun"
	return c.GetJSON(ctx, path, nil)
}

func requiredString(params map[string]any, key string) (string, error) {
	value, _ := params[key].(string)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func requiredOrganizationAndPipeline(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	pipelineID, err := requiredString(params, "pipelineId")
	if err != nil {
		return "", "", err
	}
	return organizationID, pipelineID, nil
}

func setOptionalInt(query url.Values, params map[string]any, key string) {
	setOptionalIntAs(query, params, key, key)
}

func setOptionalIntAs(query url.Values, params map[string]any, fromKey, toKey string) {
	switch value := params[fromKey].(type) {
	case float64:
		query.Set(toKey, strconv.Itoa(int(value)))
	case int:
		query.Set(toKey, strconv.Itoa(value))
	case int64:
		query.Set(toKey, strconv.FormatInt(value, 10))
	case string:
		if value != "" {
			query.Set(toKey, value)
		}
	}
}

func setOptionalString(query url.Values, params map[string]any, key string) {
	value, _ := params[key].(string)
	if value != "" {
		query.Set(key, value)
	}
}

func setOptionalBool(query url.Values, params map[string]any, key string) {
	value, ok := params[key].(bool)
	if ok {
		query.Set(key, strconv.FormatBool(value))
	}
}
