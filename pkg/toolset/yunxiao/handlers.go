package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

func handleSearchProjects(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	body := map[string]any{}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "extraConditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	if body["conditions"] == nil {
		if conditions := buildProjectConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects:search"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleGetProject(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	id, err := requiredString(params, "id")
	if err != nil {
		return "", err
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}

func handleSearchWorkitems(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	spaceID, err := requiredString(params, "spaceId")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"category": category,
		"spaceId":  spaceID,
	}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	if body["conditions"] == nil {
		if conditions := buildWorkitemConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems:search"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleGetWorkitem(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	id, err := requiredString(params, "id")
	if err != nil {
		return "", err
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}

func requiredString(params map[string]any, key string) (string, error) {
	value, _ := params[key].(string)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func buildProjectConditions(params map[string]any) string {
	filterConditions := make([]map[string]any, 0)
	if name, _ := params["name"].(string); name != "" {
		filterConditions = append(filterConditions, stringContainsCondition("name", name))
	}
	if status, _ := params["status"].(string); status != "" {
		filterConditions = append(filterConditions, listContainsCondition("status", "status", splitCSV(status)))
	}
	if creator, _ := params["creator"].(string); creator != "" {
		filterConditions = append(filterConditions, listContainsCondition("creator", "user", splitCSV(creator)))
	}
	if len(filterConditions) == 0 {
		return ""
	}
	return marshalConditions(filterConditions)
}

func buildWorkitemConditions(params map[string]any) string {
	filterConditions := make([]map[string]any, 0)
	if subject, _ := params["subject"].(string); subject != "" {
		filterConditions = append(filterConditions, stringContainsCondition("subject", subject))
	}
	if status, _ := params["status"].(string); status != "" {
		filterConditions = append(filterConditions, listContainsCondition("status", "status", splitCSV(status)))
	}
	if assignedTo, _ := params["assignedTo"].(string); assignedTo != "" {
		filterConditions = append(filterConditions, listContainsCondition("assignedTo", "user", splitCSV(assignedTo)))
	}
	if creator, _ := params["creator"].(string); creator != "" {
		filterConditions = append(filterConditions, listContainsCondition("creator", "user", splitCSV(creator)))
	}
	if tag, _ := params["tag"].(string); tag != "" {
		filterConditions = append(filterConditions, containsCondition("tag", "tag", "multiList", splitCSV(tag)))
	}
	if len(filterConditions) == 0 {
		return ""
	}
	return marshalConditions(filterConditions)
}

func stringContainsCondition(fieldIdentifier, value string) map[string]any {
	return map[string]any{
		"className":       "string",
		"fieldIdentifier": fieldIdentifier,
		"format":          "input",
		"operator":        "CONTAINS",
		"toValue":         nil,
		"value":           []string{value},
	}
}

func listContainsCondition(fieldIdentifier, className string, values []string) map[string]any {
	return containsCondition(fieldIdentifier, className, "list", values)
}

func containsCondition(fieldIdentifier, className, format string, values []string) map[string]any {
	return map[string]any{
		"className":       className,
		"fieldIdentifier": fieldIdentifier,
		"format":          format,
		"operator":        "CONTAINS",
		"toValue":         nil,
		"value":           values,
	}
}

func marshalConditions(filterConditions []map[string]any) string {
	conditions := map[string]any{
		"conditionGroups": []any{filterConditions},
	}
	data, err := json.Marshal(conditions)
	if err != nil {
		return ""
	}
	return string(data)
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
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

func setOptionalStringBody(body map[string]any, params map[string]any, key string) {
	value, _ := params[key].(string)
	if value != "" {
		body[key] = value
	}
}

func setOptionalIntBody(body map[string]any, params map[string]any, key string) {
	switch value := params[key].(type) {
	case float64:
		body[key] = int(value)
	case int:
		body[key] = value
	case int64:
		body[key] = value
	case string:
		if value != "" {
			body[key] = value
		}
	}
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
