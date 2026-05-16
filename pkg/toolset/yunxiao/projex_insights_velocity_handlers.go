package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func handleGetSprintVelocity(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Task,Bug"))
	if len(categories) == 0 {
		return "", errNoCategories
	}

	sprintCount := clampSprintCount(optionalIntDefault(params, "sprintCount", 5))

	query := url.Values{}
	query.Set("page", "1")
	query.Set("perPage", strconv.Itoa(sprintCount))
	sprintStatus := optionalStringDefault(params, "sprintStatus", "ARCHIVED,DONE")
	if sprintStatus != "" {
		query.Set("status", sprintStatus)
	}

	sprintPath := projexProjectPath(organizationID, projectID) + "/sprints"
	sprintResp, err := c.Request(ctx, http.MethodGet, sprintPath, query, nil)
	if err != nil {
		return "", fmt.Errorf("list sprints: %w", err)
	}

	result := map[string]any{
		"filters": sprintVelocityFilters(params, categories, sprintCount, sprintStatus),
		"sprints": []any{},
	}

	sprints := parseSprintList(responsePayload(sprintResp))
	velocityData := make([]map[string]any, 0, len(sprints))

	for _, sprint := range sprints {
		sprintStats := buildSprintVelocityStats(ctx, c, organizationID, projectID, sprint, categories, params)
		if sprintStats != nil {
			velocityData = append(velocityData, sprintStats)
		}
	}
	result["sprints"] = velocityData

	return marshalPretty(result)
}

func clampSprintCount(n int) int {
	if n > 20 {
		return 20
	}
	if n <= 0 {
		return 1
	}
	return n
}

func parseSprintList(data any) []any {
	switch d := data.(type) {
	case []any:
		return d
	case map[string]any:
		if data, ok := d["data"].([]any); ok {
			return data
		}
	}
	return nil
}

func buildSprintVelocityStats(ctx context.Context, c *Client, organizationID, projectID string, sprint any, categories []string, params map[string]any) map[string]any {
	sprintMap, ok := sprint.(map[string]any)
	if !ok {
		return nil
	}
	sprintID, _ := sprintMap["id"].(string)
	sprintName, _ := sprintMap["name"].(string)

	sprintStats := map[string]any{
		"id":   sprintID,
		"name": sprintName,
	}

	categoryStats := map[string]any{}
	for _, category := range categories {
		searchParams := copyParams(params)
		searchParams["sprint"] = sprintID
		searchParams["sampleLimit"] = 1000

		payload, err := searchSprintWorkitems(ctx, c, organizationID, projectID, sprintID, category, searchParams)
		if err != nil {
			continue
		}

		data, total, _ := extractWorkitemData(payload)
		completed := countCompletedWorkitems(data)

		denominator := total
		if denominator == 0 {
			denominator = 1
		}
		categoryStats[category] = map[string]any{
			"total":     total,
			"completed": completed,
			"rate":      float64(completed) / float64(denominator),
		}
	}
	sprintStats["stats"] = categoryStats
	return sprintStats
}

func countCompletedWorkitems(data []any) int {
	completed := 0
	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			if status, ok := m["status"].(map[string]any); ok {
				if stage, ok := status["stage"].(string); ok && stage == "DONE" {
					completed++
				}
			}
		}
	}
	return completed
}

func sprintVelocityFilters(params map[string]any, categories []string, sprintCount int, sprintStatus string) map[string]any {
	return map[string]any{
		"categories":   categories,
		"sprintCount":  sprintCount,
		"sprintStatus": sprintStatus,
	}
}
