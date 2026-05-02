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

func handleGetProjectWorkitemSummary(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Req,Task,Bug,Risk"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	result, err := buildCategoryResult(ctx, categories, projectWorkitemSummaryFilters(params, categories),
		func(cat string) (any, error) {
			return searchProjectWorkitemSummaryCategory(ctx, c, organizationID, projectID, cat, params)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func handleGetProjectWorkitemContext(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	projectPath := projexProjectPath(organizationID, projectID)
	payload := map[string]any{"filters": projectWorkitemContextFilters(params, category)}
	if err := addProjectWorkitemContextBaseSections(ctx, c, payload, params, projectPath, category); err != nil {
		return "", err
	}
	if err := addProjectWorkitemTypeContext(ctx, c, payload, params, projectPath); err != nil {
		return "", err
	}
	return marshalPretty(payload)
}

func addProjectWorkitemContextBaseSections(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath, category string) error {
	typeQuery := url.Values{}
	typeQuery.Set("category", category)
	types, err := getProjectOverviewSection(ctx, c, "workItemTypes", projectPath+"/workitemTypes", typeQuery)
	if err != nil {
		return err
	}
	payload["workItemTypes"] = types

	if optionalBoolDefault(params, "includeMembers", true) {
		members, err := getProjectOverviewSection(ctx, c, "members", projectPath+"/members", nil)
		if err != nil {
			return err
		}
		payload["members"] = members
	}
	if optionalBoolDefault(params, "includeLabels", true) {
		labels, err := getProjectOverviewSection(ctx, c, "labels", projectPath+"/labels", projectOverviewListQuery(params, false))
		if err != nil {
			return err
		}
		payload["labels"] = labels
	}
	return nil
}

func addProjectWorkitemTypeContext(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath string) error {
	workItemTypeID, _ := params["workItemTypeId"].(string)
	if strings.TrimSpace(workItemTypeID) == "" {
		return nil
	}
	typePath := projectPath + "/workitemTypes/" + url.PathEscape(strings.TrimSpace(workItemTypeID))
	if optionalBoolDefault(params, "includeFields", true) {
		fields, err := getProjectOverviewSection(ctx, c, "fields", typePath+"/fields", nil)
		if err != nil {
			return err
		}
		payload["fields"] = fields
	}
	if optionalBoolDefault(params, "includeWorkflow", true) {
		workflow, err := getProjectOverviewSection(ctx, c, "workflow", typePath+"/workflows", nil)
		if err != nil {
			return err
		}
		payload["workflow"] = workflow
	}
	return nil
}

func projectWorkitemContextFilters(params map[string]any, category string) map[string]any {
	return map[string]any{
		"category":       category,
		"workItemTypeId": optionalStringDefault(params, "workItemTypeId", ""),
		"page":           optionalIntDefault(params, "page", 1),
		"perPage":        optionalIntDefault(params, "perPage", 20),
	}
}

func searchProjectWorkitemSummaryCategory(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func projectWorkitemSummaryBody(projectID, category string, params map[string]any) map[string]any {
	body := map[string]any{
		"category": category,
		"spaceId":  projectID,
		"page":     1,
		"perPage":  normalizedSampleLimit(params),
	}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	if body["conditions"] == nil {
		if conditions := buildWorkitemConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}
	return body
}

func projectWorkitemSummaryFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"subject":     optionalStringDefault(params, "subject", ""),
		"status":      optionalStringDefault(params, "status", ""),
		"assignedTo":  optionalStringDefault(params, "assignedTo", ""),
		"creator":     optionalStringDefault(params, "creator", ""),
		"tag":         optionalStringDefault(params, "tag", ""),
		"orderBy":     optionalStringDefault(params, "orderBy", ""),
		"sort":        optionalStringDefault(params, "sort", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func normalizedSampleLimit(params map[string]any) int {
	limit := optionalIntDefault(params, "sampleLimit", 5)
	if limit < 0 {
		return 0
	}
	if limit > 200 {
		return 200
	}
	return limit
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

func handleGetMyProjectWorkitems(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	userID, err := requiredString(params, "userId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	relation := optionalStringDefault(params, "relation", "assigned")
	searchParams := copyParams(params)
	switch relation {
	case "assigned":
		searchParams["assignedTo"] = userID
	case "created":
		searchParams["creator"] = userID
	default:
		return "", fmt.Errorf("relation must be assigned or created")
	}

	result, err := buildCategoryResult(ctx, categories, myProjectWorkitemFilters(params, userID, relation, categories),
		func(cat string) (any, error) {
			return searchProjectWorkitemSummaryCategory(ctx, c, organizationID, projectID, cat, searchParams)
		})
	if err != nil {
		return "", err
	}
	return marshalPretty(result)
}

func myProjectWorkitemFilters(params map[string]any, userID, relation string, categories []string) map[string]any {
	return map[string]any{
		"userId":      userID,
		"relation":    relation,
		"categories":  categories,
		"status":      optionalStringDefault(params, "status", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func handleGetSprintOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	sprintID, err := requiredString(params, "sprintId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", fmt.Errorf("categories must include at least one category")
	}

	sprintPath := projexProjectPath(organizationID, projectID) + "/sprints/" + url.PathEscape(strings.TrimSpace(sprintID))
	sprintResp, err := c.Request(ctx, http.MethodGet, sprintPath, nil, nil)
	if err != nil {
		return "", fmt.Errorf("sprint: %w", err)
	}

	result, err := buildCategoryResult(ctx, categories, sprintOverviewFilters(params, categories),
		func(cat string) (any, error) {
			return searchSprintWorkitems(ctx, c, organizationID, projectID, sprintID, cat, params)
		})
	if err != nil {
		return "", err
	}
	result["sprint"] = responsePayload(sprintResp)
	return marshalPretty(result)
}

func buildCategoryResult(ctx context.Context, categories []string, filters map[string]any, searchFn func(string) (any, error)) (map[string]any, error) {
	result := map[string]any{
		"filters":    filters,
		"categories": map[string]any{},
	}
	payloads := result["categories"].(map[string]any)
	for _, category := range categories {
		payload, err := searchFn(category)
		if err != nil {
			return nil, err
		}
		payloads[category] = payload
	}
	return result, nil
}

func searchSprintWorkitems(ctx context.Context, c *Client, organizationID, projectID, sprintID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	sprintConditions := buildWorkitemConditions(map[string]any{"sprint": sprintID})
	if sprintConditions != "" {
		existing, _ := body["conditions"].(string)
		body["conditions"] = mergeConditions(existing, sprintConditions)
	}
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func mergeConditions(existing, extra string) string {
	if existing == "" {
		return extra
	}
	if extra == "" {
		return existing
	}
	var existingArr, extraArr []map[string]any
	_ = json.Unmarshal([]byte(existing), &existingArr)
	_ = json.Unmarshal([]byte(extra), &extraArr)
	merged := append(existingArr, extraArr...)
	mergedBytes, _ := json.Marshal(merged)
	return string(mergedBytes)
}

func sprintOverviewFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"status":      optionalStringDefault(params, "status", ""),
		"sampleLimit": normalizedSampleLimit(params),
	}
}

func handleGetProjectWorkitemBoard(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	body := projectWorkitemSummaryBody(projectID, category, params)
	// sprint filter is already handled by buildWorkitemConditions inside projectWorkitemSummaryBody
	// when params contains a "sprint" key.
	sprintID := optionalStringDefault(params, "sprint", "")

	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return "", fmt.Errorf("search: %w", err)
	}

	payload := responsePayload(resp)
	board := map[string]any{
		"filters": map[string]any{
			"category":    category,
			"sprint":      sprintID,
			"sampleLimit": normalizedSampleLimit(params),
		},
		"columns": map[string]any{},
	}

	data, total, err := extractWorkitemData(payload)
	if err != nil {
		return "", err
	}
	board["total"] = total

	columns := board["columns"].(map[string]any)
	for _, item := range data {
		statusName := extractStatusName(item)
		if statusName == "" {
			statusName = "Unknown"
		}
		if columns[statusName] == nil {
			columns[statusName] = []any{}
		}
		columns[statusName] = append(columns[statusName].([]any), item)
	}

	return marshalPretty(board)
}

func extractWorkitemData(payload any) ([]any, int, error) {
	switch p := payload.(type) {
	case []any:
		return p, len(p), nil
	case map[string]any:
		data, _ := p["data"].([]any)
		total := 0
		if pagination, ok := p["pagination"].(map[string]any); ok {
			if t, ok := pagination["total"].(float64); ok {
				total = int(t)
			}
		}
		return data, total, nil
	default:
		return nil, 0, fmt.Errorf("unexpected payload type %T", payload)
	}
}

func extractStatusName(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}
	status, ok := m["status"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := status["name"].(string)
	return name
}

func marshalPretty(value any) (string, error) {
	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}
