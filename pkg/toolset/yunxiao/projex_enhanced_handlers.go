package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func handleGetProjectOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	projectPath := projexProjectPath(organizationID, projectID)
	overview := map[string]any{
		"project": nil,
		"filters": projectOverviewFilters(params),
	}

	project, err := getProjectOverviewSection(ctx, c, "project", projectPath, nil)
	if err != nil {
		return "", err
	}
	overview["project"] = project

	for _, section := range projectOverviewSections(projectPath, params) {
		if err := addProjectOverviewSection(ctx, c, overview, params, section); err != nil {
			return "", err
		}
	}

	return marshalPretty(overview)
}

type projectOverviewSection struct {
	flag  string
	name  string
	path  string
	query url.Values
}

func projectOverviewSections(projectPath string, params map[string]any) []projectOverviewSection {
	return []projectOverviewSection{
		{flag: "includeMembers", name: "members", path: projectPath + "/members"},
		{flag: "includeSprints", name: "sprints", path: projectPath + "/sprints", query: projectOverviewListQuery(params, true)},
		{flag: "includeMilestones", name: "milestones", path: projectPath + "/milestones", query: projectOverviewListQuery(params, true)},
		{flag: "includeVersions", name: "versions", path: projectPath + "/versions", query: projectOverviewListQuery(params, true)},
		{flag: "includeLabels", name: "labels", path: projectPath + "/labels", query: projectOverviewListQuery(params, false)},
	}
}

func addProjectOverviewSection(ctx context.Context, c *Client, overview map[string]any, params map[string]any, section projectOverviewSection) error {
	if !optionalBoolDefault(params, section.flag, true) {
		return nil
	}
	payload, err := getProjectOverviewSection(ctx, c, section.name, section.path, section.query)
	if err != nil {
		return err
	}
	overview[section.name] = payload
	return nil
}

func getProjectOverviewSection(ctx context.Context, c *Client, name, path string, query url.Values) (any, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return responsePayload(resp), nil
}

func projectOverviewListQuery(params map[string]any, withStatus bool) url.Values {
	query := url.Values{}
	query.Set("page", strconv.Itoa(optionalIntDefault(params, "page", 1)))
	query.Set("perPage", strconv.Itoa(optionalIntDefault(params, "perPage", 20)))
	if withStatus && optionalBoolDefault(params, "activeOnly", true) {
		status := optionalStringDefault(params, "status", "TODO,DOING")
		if status != "" {
			query.Set("status", status)
		}
	}
	return query
}

func projectOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"activeOnly": optionalBoolDefault(params, "activeOnly", true),
		"status":     optionalStringDefault(params, "status", "TODO,DOING"),
		"page":       optionalIntDefault(params, "page", 1),
		"perPage":    optionalIntDefault(params, "perPage", 20),
	}
}

func optionalIntDefault(params map[string]any, key string, defaultValue int) int {
	switch value := params[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	case int64:
		return int(value)
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func responsePayload(resp *Response) any {
	var data any
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		data = string(resp.Body)
	}

	if resp.Pagination == nil && resp.NextToken == "" && resp.RequestID == "" {
		return data
	}

	payload := map[string]any{"data": data}
	if resp.Pagination != nil {
		payload["pagination"] = resp.Pagination
	}
	if resp.NextToken != "" {
		payload["nextToken"] = resp.NextToken
	}
	if resp.RequestID != "" {
		payload["requestId"] = resp.RequestID
	}
	return payload
}

func marshalPretty(value any) (string, error) {
	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}
