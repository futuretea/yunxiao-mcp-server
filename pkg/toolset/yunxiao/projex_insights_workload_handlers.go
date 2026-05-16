package yunxiao

import (
	"context"
	"fmt"
	"time"
)

func handleGetMemberWorkloadTrend(ctx context.Context, client any, params map[string]any) (string, error) {
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

	daysBack := optionalIntDefault(params, "daysBack", 30)
	memberLimit := optionalIntDefault(params, "memberLimit", 20)

	members, assigneeIDs, err := projectTaskStatusMembers(ctx, c, organizationID, projectID, params)
	if err != nil {
		return "", fmt.Errorf("members: %w", err)
	}
	if len(assigneeIDs) == 0 {
		return "", fmt.Errorf("assigneeIds is empty and no project members with userId were returned")
	}

	cutoffDate := time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")

	memberWorkloads := make(map[string]any, len(assigneeIDs))
	totalActiveTasks := 0
	totalOverdue := 0

	for _, assigneeID := range assigneeIDs {
		workload, active, overdue := buildMemberWorkload(ctx, c, organizationID, projectID, assigneeID, members[assigneeID], categories, params, cutoffDate)
		memberWorkloads[assigneeID] = workload
		totalActiveTasks += active
		totalOverdue += overdue
	}

	denominator := len(assigneeIDs)
	if denominator == 0 {
		denominator = 1
	}

	result := map[string]any{
		"filters": memberWorkloadTrendFilters(params, categories, daysBack, memberLimit),
		"members": memberWorkloads,
		"summary": map[string]any{
			"memberCount":       len(assigneeIDs),
			"totalActiveTasks":  totalActiveTasks,
			"totalOverdue":      totalOverdue,
			"avgTasksPerMember": float64(totalActiveTasks) / float64(denominator),
		},
	}

	return marshalPretty(result)
}

func buildMemberWorkload(ctx context.Context, c *Client, organizationID, projectID, assigneeID string, member any, categories []string, params map[string]any, cutoffDate string) (workload map[string]any, totalActive, totalOverdue int) {
	workload = map[string]any{
		"member":        member,
		"tasksByStatus": map[string]int{},
		"totalAssigned": 0,
		"overdueCount":  0,
	}

	for _, category := range categories {
		searchParams := copyParams(params)
		searchParams["assignedTo"] = assigneeID
		searchParams["sampleLimit"] = 1000

		payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, category, searchParams)
		if err != nil {
			continue
		}

		data, _, _ := extractWorkitemData(payload)
		active, overdue := countWorkitemStatuses(data, workload)
		totalActive += active
		totalOverdue += overdue

		recentParams := copyParams(searchParams)
		recentParams["updatedAfter"] = cutoffDate
		recentPayload, _ := searchProjectWorkitems(ctx, c, organizationID, projectID, category, recentParams)
		recentData, _, _ := extractWorkitemData(recentPayload)
		workload["recentActivity"] = map[string]any{
			"updatedInPeriod": len(recentData),
		}
	}

	return workload, totalActive, totalOverdue
}

func countWorkitemStatuses(data []any, workload map[string]any) (totalActive, totalOverdue int) {
	tasksByStatus := workload["tasksByStatus"].(map[string]int)

	for _, item := range data {
		if itemMap, ok := item.(map[string]any); ok {
			workload["totalAssigned"] = workload["totalAssigned"].(int) + 1
			totalActive++

			statusName := extractWorkitemStatusName(itemMap)
			tasksByStatus[statusName]++

			if finishTime, ok := itemMap["finishTime"].(string); ok && finishTime != "" {
				if finishTime < todayDate() {
					workload["overdueCount"] = workload["overdueCount"].(int) + 1
					totalOverdue++
				}
			}
		}
	}
	return totalActive, totalOverdue
}

func memberWorkloadTrendFilters(params map[string]any, categories []string, daysBack, memberLimit int) map[string]any {
	return map[string]any{
		"categories":  categories,
		"daysBack":    daysBack,
		"memberLimit": memberLimit,
	}
}

func handleGetTeamWorkloadBreakdown(ctx context.Context, client any, params map[string]any) (string, error) {
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

	memberLimit := optionalIntDefault(params, "memberLimit", 20)
	taskLimit := optionalIntDefault(params, "taskLimit", 10)
	if taskLimit < 1 {
		taskLimit = 1
	}
	if taskLimit > 50 {
		taskLimit = 50
	}

	members, assigneeIDs, err := projectTaskStatusMembers(ctx, c, organizationID, projectID, params)
	if err != nil {
		return "", fmt.Errorf("members: %w", err)
	}
	if len(assigneeIDs) == 0 {
		return "", fmt.Errorf("assigneeIds is empty and no project members with userId were returned")
	}

	result := map[string]any{
		"filters": teamWorkloadBreakdownFilters(params, categories, memberLimit, taskLimit),
		"members": map[string]any{},
	}

	memberBreakdowns := result["members"].(map[string]any)

	for _, assigneeID := range assigneeIDs {
		member := members[assigneeID]
		breakdown := map[string]any{
			"member": member,
			"tasks":  []map[string]any{},
		}

		for _, category := range categories {
			searchParams := copyParams(params)
			searchParams["assignedTo"] = assigneeID
			searchParams["sampleLimit"] = taskLimit

			payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, category, searchParams)
			if err != nil {
				continue
			}

			data, _, _ := extractWorkitemData(payload)

			for _, item := range data {
				if itemMap, ok := item.(map[string]any); ok {
					task := map[string]any{
						"id":           itemMap["id"],
						"serialNumber": itemMap["serialNumber"],
						"subject":      itemMap["subject"],
						"status":       extractWorkitemStatusName(itemMap),
						"labels":       extractLabelNames(itemMap),
						"gmtCreate":    itemMap["gmtCreate"],
						"category":     category,
					}
					breakdown["tasks"] = append(breakdown["tasks"].([]map[string]any), task)
				}
			}
		}

		memberBreakdowns[assigneeID] = breakdown
	}

	return marshalPretty(result)
}

func teamWorkloadBreakdownFilters(params map[string]any, categories []string, memberLimit, taskLimit int) map[string]any {
	return map[string]any{
		"categories":  categories,
		"memberLimit": memberLimit,
		"taskLimit":   taskLimit,
		"status":      optionalStringDefault(params, "status", ""),
	}
}

func extractWorkitemStatusName(itemMap map[string]any) string {
	if status, ok := itemMap["status"].(map[string]any); ok {
		if name, ok := status["name"].(string); ok {
			return name
		}
	}
	return "Unknown"
}

func extractLabelNames(itemMap map[string]any) []string {
	labels, _ := itemMap["labels"].([]any)
	names := make([]string, 0, len(labels))
	for _, l := range labels {
		if labelMap, ok := l.(map[string]any); ok {
			if name, ok := labelMap["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return names
}
