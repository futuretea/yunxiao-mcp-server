package yunxiao

import (
	"context"
	"net/http"
)

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

	blockedWorkitems := make([]any, 0)
	blockingWorkitems := make([]any, 0)
	blockedCounts := make(map[string]int)
	blockingCounts := make(map[string]int)

	for _, category := range categories {
		data := fetchCategoryWorkitems(ctx, c, organizationID, projectID, category, params)
		for _, item := range data {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			itemID, _ := itemMap["id"].(string)
			if itemID == "" {
				continue
			}

			hasBlockingDeps, hasBlockedItems := checkWorkitemBlockers(ctx, c, organizationID, itemID)

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

	result := map[string]any{
		"filters": blockerAnalysisFilters(params, categories),
		"blocked": map[string]any{
			"workitems":  blockedWorkitems,
			"byCategory": blockedCounts,
		},
		"blocking": map[string]any{
			"workitems":  blockingWorkitems,
			"byCategory": blockingCounts,
		},
		"summary": map[string]any{
			"totalBlocked":  len(blockedWorkitems),
			"totalBlocking": len(blockingWorkitems),
		},
	}

	return marshalPretty(result)
}

func fetchCategoryWorkitems(ctx context.Context, c *Client, organizationID, projectID, category string, params map[string]any) []any {
	searchParams := copyParams(params)
	searchParams["perPage"] = normalizedSampleLimit(params)

	body := projectWorkitemSummaryBody(projectID, category, searchParams)
	path := projexOrganizationPath(organizationID) + "/workitems:search"
	resp, err := c.Request(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return nil
	}

	payload := responsePayload(resp)
	data, _, _ := extractWorkitemData(payload)
	return data
}

func checkWorkitemBlockers(ctx context.Context, c *Client, organizationID, itemID string) (hasBlockingDeps, hasBlockedItems bool) {
	relPath := projexWorkitemPath(organizationID, itemID) + "/relationRecords"
	relResp, err := c.Request(ctx, http.MethodGet, relPath, nil, nil)
	if err != nil {
		return false, false
	}

	relations := responsePayload(relResp)
	relList := parseAnyList(relations)

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
	return hasBlockingDeps, hasBlockedItems
}

func parseAnyList(data any) []any {
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

func blockerAnalysisFilters(params map[string]any, categories []string) map[string]any {
	return map[string]any{
		"categories":  categories,
		"sampleLimit": normalizedSampleLimit(params),
	}
}
