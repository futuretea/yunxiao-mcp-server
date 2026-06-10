package yunxiao

import (
	"context"
	"fmt"
	"net/http"
)

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

	activityList := parseListData(activities)
	if activityList == nil {
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
