package yunxiao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func handleGetProjectRiskDashboard(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	categories := splitCSV(optionalStringDefault(params, "categories", "Risk,Bug,Task"))
	if len(categories) == 0 {
		return "", errNoCategories
	}

	dashboard := map[string]any{
		"filters":      riskDashboardFilters(params, categories),
		"byCategory":   map[string]any{},
		"overdue":      nil,
		"highPriority": nil,
		"stale":        nil,
	}
	if err := addRiskCategorySections(ctx, c, dashboard["byCategory"].(map[string]any), organizationID, projectID, categories, params); err != nil {
		return "", err
	}
	if err := addRiskFocusSections(ctx, c, dashboard, organizationID, projectID, categories, params); err != nil {
		return "", err
	}
	return marshalPretty(dashboard)
}

func handleGetProjectMemberTaskStatus(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	members, assigneeIDs, err := projectTaskStatusMembers(ctx, c, organizationID, projectID, params)
	if err != nil {
		return "", err
	}
	if len(assigneeIDs) == 0 {
		return "", fmt.Errorf("assigneeIds is empty and no project members with userId were returned")
	}
	groups, err := parseStatusGroups(params)
	if err != nil {
		return "", err
	}

	status := map[string]any{
		"filters": projectTaskStatusFilters(params, assigneeIDs, groups),
		"members": map[string]any{},
	}
	memberPayloads := status["members"].(map[string]any)
	for _, assigneeID := range assigneeIDs {
		payload, err := projectMemberTaskPayload(ctx, c, organizationID, projectID, assigneeID, members[assigneeID], groups, params)
		if err != nil {
			return "", err
		}
		memberPayloads[assigneeID] = payload
	}
	return marshalPretty(status)
}

func addRiskCategorySections(ctx context.Context, c *Client, target map[string]any, organizationID, projectID string, categories []string, params map[string]any) error {
	for _, category := range categories {
		payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, category, params)
		if err != nil {
			return err
		}
		target[category] = payload
	}
	return nil
}

func addRiskFocusSections(ctx context.Context, c *Client, dashboard map[string]any, organizationID, projectID string, categories []string, params map[string]any) error {
	overdueParams := copyParams(params)
	overdueParams["finishTimeBefore"] = optionalStringDefault(params, "overdueBefore", todayDate())
	overdue, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), overdueParams)
	if err != nil {
		return err
	}
	dashboard["overdue"] = overdue

	if highPriority := optionalStringDefault(params, "highPriority", ""); highPriority != "" {
		priorityParams := copyParams(params)
		priorityParams["priority"] = highPriority
		highPriorityPayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), priorityParams)
		if err != nil {
			return err
		}
		dashboard["highPriority"] = highPriorityPayload
	}
	if staleBefore := optionalStringDefault(params, "staleBefore", ""); staleBefore != "" {
		staleParams := copyParams(params)
		staleParams["updateStatusAtBefore"] = staleBefore
		stalePayload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, strings.Join(categories, ","), staleParams)
		if err != nil {
			return err
		}
		dashboard["stale"] = stalePayload
	}
	return nil
}

func projectMemberTaskPayload(ctx context.Context, c *Client, organizationID, projectID, assigneeID string, member any, groups map[string]string, params map[string]any) (map[string]any, error) {
	memberParams := copyParams(params)
	memberParams["assignedTo"] = assigneeID
	categories := optionalStringDefault(params, "categories", "Task,Bug")

	assigned, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, memberParams)
	if err != nil {
		return nil, err
	}
	overdueParams := copyParams(memberParams)
	overdueParams["finishTimeBefore"] = optionalStringDefault(params, "overdueBefore", todayDate())
	overdue, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, overdueParams)
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"member":   member,
		"assigned": assigned,
		"overdue":  overdue,
	}
	if len(groups) > 0 {
		groupPayloads, err := projectMemberStatusGroups(ctx, c, organizationID, projectID, categories, memberParams, groups)
		if err != nil {
			return nil, err
		}
		payload["statusGroups"] = groupPayloads
	}
	return payload, nil
}

func projectMemberStatusGroups(ctx context.Context, c *Client, organizationID, projectID, categories string, baseParams map[string]any, groups map[string]string) (map[string]any, error) {
	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	sort.Strings(names)
	payloads := make(map[string]any, len(groups))
	for _, name := range names {
		groupParams := copyParams(baseParams)
		groupParams["status"] = groups[name]
		payload, err := searchProjectWorkitems(ctx, c, organizationID, projectID, categories, groupParams)
		if err != nil {
			return nil, err
		}
		payloads[name] = payload
	}
	return payloads, nil
}

func projectTaskStatusMembers(ctx context.Context, c *Client, organizationID, projectID string, params map[string]any) (map[string]any, []string, error) {
	if assigneeIDs := splitCSV(optionalStringDefault(params, "assigneeIds", "")); len(assigneeIDs) > 0 {
		members := make(map[string]any, len(assigneeIDs))
		for _, assigneeID := range assigneeIDs {
			members[assigneeID] = map[string]any{"userId": assigneeID}
		}
		return members, assigneeIDs, nil
	}

	resp, err := c.Request(ctx, http.MethodGet, projexProjectPath(organizationID, projectID)+"/members", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("members: %w", err)
	}
	return projectMembersFromResponse(resp, optionalIntDefault(params, "memberLimit", 20))
}

func projectMembersFromResponse(resp *Response, limit int) (map[string]any, []string, error) {
	var members []map[string]any
	if err := json.Unmarshal(resp.Body, &members); err != nil {
		return nil, nil, fmt.Errorf("decode members: %w", err)
	}
	if limit < 0 {
		limit = 0
	}
	result := make(map[string]any, len(members))
	ids := make([]string, 0, len(members))
	for _, member := range members {
		userID, _ := member["userId"].(string)
		if strings.TrimSpace(userID) == "" {
			continue
		}
		if limit > 0 && len(ids) >= limit {
			break
		}
		userID = strings.TrimSpace(userID)
		result[userID] = member
		ids = append(ids, userID)
	}
	return result, ids, nil
}

func searchProjectWorkitems(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) (any, error) {
	body := projectWorkitemSummaryBody(projectID, category, params)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", category, err)
	}
	return responsePayload(resp), nil
}

func parseStatusGroups(params map[string]any) (map[string]string, error) {
	raw := optionalStringDefault(params, "statusGroups", "")
	if raw == "" {
		return nil, nil
	}
	var groups map[string]string
	if err := json.Unmarshal([]byte(raw), &groups); err != nil {
		return nil, fmt.Errorf("statusGroups must be a JSON object of name to comma-separated status IDs: %w", err)
	}
	for name, value := range groups {
		if strings.TrimSpace(name) == "" || strings.TrimSpace(value) == "" {
			delete(groups, name)
			continue
		}
		groups[name] = strings.TrimSpace(value)
	}
	return groups, nil
}

func riskDashboardFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":    categories,
		"subject":       optionalStringDefault(params, "subject", ""),
		"status":        optionalStringDefault(params, "status", ""),
		"statusStage":   optionalStringDefault(params, "statusStage", ""),
		"assignedTo":    optionalStringDefault(params, "assignedTo", ""),
		"creator":       optionalStringDefault(params, "creator", ""),
		"sprint":        optionalStringDefault(params, "sprint", ""),
		"workitemType":  optionalStringDefault(params, "workitemType", ""),
		"tag":           optionalStringDefault(params, "tag", ""),
		"overdueBefore": optionalStringDefault(params, "overdueBefore", todayDate()),
		"highPriority":  optionalStringDefault(params, "highPriority", ""),
		"staleBefore":   optionalStringDefault(params, "staleBefore", ""),
		"sampleLimit":   normalizedSampleLimit(params),
	}
}

func projectTaskStatusFilters(params map[string]any, assigneeIDs []string, groups map[string]string) map[string]any {
	return map[string]any{
		"assigneeIds":   assigneeIDs,
		"categories":    splitCSV(optionalStringDefault(params, "categories", "Task,Bug")),
		"subject":       optionalStringDefault(params, "subject", ""),
		"status":        optionalStringDefault(params, "status", ""),
		"statusStage":   optionalStringDefault(params, "statusStage", ""),
		"assignedTo":    optionalStringDefault(params, "assignedTo", ""),
		"creator":       optionalStringDefault(params, "creator", ""),
		"sprint":        optionalStringDefault(params, "sprint", ""),
		"workitemType":  optionalStringDefault(params, "workitemType", ""),
		"tag":           optionalStringDefault(params, "tag", ""),
		"overdueBefore": optionalStringDefault(params, "overdueBefore", todayDate()),
		"statusGroups":  groups,
		"memberLimit":   optionalIntDefault(params, "memberLimit", 20),
		"sampleLimit":   normalizedSampleLimit(params),
	}
}

func copyParams(params map[string]any) map[string]any {
	copied := make(map[string]any, len(params))
	for key, value := range params {
		copied[key] = value
	}
	return copied
}

func todayDate() string {
	return time.Now().Format("2006-01-02")
}

// --- Predictive Risk Assessment Handlers ---

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

	sprintCount := optionalIntDefault(params, "sprintCount", 5)
	if sprintCount > 20 {
		sprintCount = 20
	}
	if sprintCount <= 0 {
		sprintCount = 1
	}

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

	sprintData := responsePayload(sprintResp)
	var sprints []any
	switch d := sprintData.(type) {
	case []any:
		sprints = d
	case map[string]any:
		if data, ok := d["data"].([]any); ok {
			sprints = data
		}
	}

	velocityData := make([]map[string]any, 0, len(sprints))
	for _, sprint := range sprints {
		sprintMap, ok := sprint.(map[string]any)
		if !ok {
			continue
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
		velocityData = append(velocityData, sprintStats)
	}
	result["sprints"] = velocityData

	return marshalPretty(result)
}

func sprintVelocityFilters(params map[string]any, categories []string, sprintCount int, sprintStatus string) map[string]any {
	return map[string]any{
		"categories":   categories,
		"sprintCount":  sprintCount,
		"sprintStatus": sprintStatus,
	}
}

func handleGetWorkitemStatusTimeline(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	workitemID, err := requiredString(params, "workitemId")
	if err != nil {
		return "", err
	}

	result := map[string]any{
		"filters": workitemStatusTimelineFilters(params),
	}

	if optionalBoolDefault(params, "includeWorkitem", true) {
		workitemPath := projexWorkitemPath(organizationID, workitemID)
		workitem, err := getProjectOverviewSection(ctx, c, "workitem", workitemPath, nil)
		if err != nil {
			return "", err
		}
		result["workitem"] = workitem
	}

	activitiesPath := projexWorkitemPath(organizationID, workitemID) + "/activities"
	activitiesResp, err := c.Request(ctx, http.MethodGet, activitiesPath, nil, nil)
	if err != nil {
		return "", fmt.Errorf("activities: %w", err)
	}
	activities := responsePayload(activitiesResp)

	timeline := parseStatusTimeline(activities)
	result["timeline"] = timeline

	if len(timeline) > 0 {
		stats := calculateTimelineStats(timeline)
		result["summary"] = stats
	}

	return marshalPretty(result)
}

func workitemStatusTimelineFilters(params map[string]any) map[string]any {
	return map[string]any{
		"workitemId":      optionalStringDefault(params, "workitemId", ""),
		"includeWorkitem": optionalBoolDefault(params, "includeWorkitem", true),
	}
}

func parseStatusTimeline(activities any) []map[string]any {
	timeline := make([]map[string]any, 0)

	var activityList []any
	switch a := activities.(type) {
	case []any:
		activityList = a
	case map[string]any:
		if data, ok := a["data"].([]any); ok {
			activityList = data
		}
	default:
		return timeline
	}

	for _, act := range activityList {
		actMap, ok := act.(map[string]any)
		if !ok {
			continue
		}

		action, _ := actMap["action"].(string)
		if action == "UPDATE" {
			field, _ := actMap["field"].(string)
			if field == "status" {
				timeline = append(timeline, map[string]any{
					"timestamp":  actMap["gmtCreate"],
					"operator":   actMap["operator"],
					"fromStatus": actMap["oldValue"],
					"toStatus":   actMap["newValue"],
					"comment":    actMap["comment"],
				})
			}
		}
	}

	return timeline
}

func calculateTimelineStats(timeline []map[string]any) map[string]any {
	if len(timeline) == 0 {
		return map[string]any{"totalChanges": 0}
	}

	statuses := make(map[string]int)
	for _, entry := range timeline {
		if to, ok := entry["toStatus"].(map[string]any); ok {
			if name, ok := to["name"].(string); ok {
				statuses[name]++
			}
		}
	}

	return map[string]any{
		"totalChanges": len(timeline),
		"statusVisits": statuses,
	}
}

func handleGetBlockerAnalysis(ctx context.Context, client any, params map[string]any) (string, error) {
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

	result := map[string]any{
		"filters":  blockerAnalysisFilters(params, categories),
		"blocked":  map[string]any{},
		"blocking": map[string]any{},
		"summary":  map[string]any{},
	}

	blockedWorkitems := make([]any, 0)
	blockingWorkitems := make([]any, 0)
	blockedCounts := make(map[string]int)
	blockingCounts := make(map[string]int)

	for _, category := range categories {
		searchParams := copyParams(params)
		searchParams["perPage"] = normalizedSampleLimit(params)

		body := projectWorkitemSummaryBody(projectID, category, searchParams)
		path := projexOrganizationPath(organizationID) + "/workitems:search"
		resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
		if err != nil {
			continue
		}

		payload := responsePayload(resp)
		data, _, _ := extractWorkitemData(payload)

		for _, item := range data {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			itemID, _ := itemMap["id"].(string)
			if itemID == "" {
				continue
			}

			relPath := projexWorkitemPath(organizationID, itemID) + "/relationRecords"
			relResp, err := c.Request(ctx, http.MethodGet, relPath, nil, nil)
			if err != nil {
				continue
			}
			relations := responsePayload(relResp)

			hasBlockingDeps := false
			hasBlockedItems := false

			var relList []any
			switch r := relations.(type) {
			case []any:
				relList = r
			case map[string]any:
				if data, ok := r["data"].([]any); ok {
					relList = data
				}
			}

			for _, rel := range relList {
				relMap, ok := rel.(map[string]any)
				if !ok {
					continue
				}
				relType, _ := relMap["relationType"].(string)

				if target, ok := relMap["target"].(map[string]any); ok {
					if status, ok := target["status"].(map[string]any); ok {
						stage, _ := status["stage"].(string)
						if relType == "DEPEND_ON" && stage != "DONE" {
							hasBlockingDeps = true
						}
						if relType == "DEPENDED_BY" && stage != "DONE" {
							hasBlockedItems = true
						}
					}
				}
			}

			if hasBlockingDeps {
				blockedWorkitems = append(blockedWorkitems, itemMap)
				blockedCounts[category]++
			}
			if hasBlockedItems {
				blockingWorkitems = append(blockingWorkitems, itemMap)
				blockingCounts[category]++
			}
		}
	}

	result["blocked"] = map[string]any{
		"workitems":  blockedWorkitems,
		"byCategory": blockedCounts,
	}
	result["blocking"] = map[string]any{
		"workitems":  blockingWorkitems,
		"byCategory": blockingCounts,
	}
	result["summary"] = map[string]any{
		"totalBlocked":  len(blockedWorkitems),
		"totalBlocking": len(blockingWorkitems),
	}

	return marshalPretty(result)
}

func blockerAnalysisFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":    categories,
		"sampleLimit":   normalizedSampleLimit(params),
	}
}

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

	result := map[string]any{
		"filters": memberWorkloadTrendFilters(params, categories, daysBack, memberLimit),
		"members": map[string]any{},
		"summary": map[string]any{},
	}

	totalActiveTasks := 0
	totalOverdue := 0
	memberWorkloads := result["members"].(map[string]any)

	for _, assigneeID := range assigneeIDs {
		workload := map[string]any{
			"member":        members[assigneeID],
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

			for _, item := range data {
				if itemMap, ok := item.(map[string]any); ok {
					workload["totalAssigned"] = workload["totalAssigned"].(int) + 1
					totalActiveTasks++

					statusName := "Unknown"
					if status, ok := itemMap["status"].(map[string]any); ok {
						if name, ok := status["name"].(string); ok {
							statusName = name
						}
					}
					workload["tasksByStatus"].(map[string]int)[statusName]++

					if finishTime, ok := itemMap["finishTime"].(string); ok && finishTime != "" {
						if finishTime < todayDate() {
							workload["overdueCount"] = workload["overdueCount"].(int) + 1
							totalOverdue++
						}
					}
				}
			}

			recentParams := copyParams(searchParams)
			recentParams["updatedAfter"] = cutoffDate
			recentPayload, _ := searchProjectWorkitems(ctx, c, organizationID, projectID, category, recentParams)
			recentData, _, _ := extractWorkitemData(recentPayload)
			workload["recentActivity"] = map[string]any{
				"updatedInPeriod": len(recentData),
			}
		}

		memberWorkloads[assigneeID] = workload
	}

	denominator := len(assigneeIDs)
	if denominator == 0 {
		denominator = 1
	}
	result["summary"] = map[string]any{
		"memberCount":       len(assigneeIDs),
		"totalActiveTasks":  totalActiveTasks,
		"totalOverdue":      totalOverdue,
		"avgTasksPerMember": float64(totalActiveTasks) / float64(denominator),
	}

	return marshalPretty(result)
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
						"id":          itemMap["id"],
						"serialNumber": itemMap["serialNumber"],
						"subject":     itemMap["subject"],
						"status":      extractWorkitemStatusName(itemMap),
						"labels":      extractLabelNames(itemMap),
						"gmtCreate":   itemMap["gmtCreate"],
						"category":    category,
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
