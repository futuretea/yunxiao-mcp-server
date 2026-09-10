package yunxiao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func handleGetSystemOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, err := requiredOrganizationAndSystem(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	systemPath := appstackSystemPath(organizationID, systemName)

	system, err := findAppStackSystem(ctx, c, organizationID, systemName)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"system":  system,
		"filters": systemOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeApps", true) {
		appQuery := appstackDefaultPageQuery(params)
		appQuery.Set("pageSize", strconv.Itoa(optionalIntDefault(params, "appLimit", 10)))
		apps, err := getProjectOverviewSection(ctx, c, "apps", systemPath+"/apps", appQuery)
		if err != nil {
			return "", err
		}
		overview["apps"] = apps
	}

	if optionalBoolDefault(params, "includeMembers", true) {
		memberQuery := appstackDefaultPageQuery(params)
		memberQuery.Set("pageSize", strconv.Itoa(optionalIntDefault(params, "memberLimit", 10)))
		members, err := getProjectOverviewSection(ctx, c, "members", systemPath+"/members", memberQuery)
		if err != nil {
			return "", err
		}
		overview["members"] = members
	}

	return marshalPretty(overview)
}

func findAppStackSystem(ctx context.Context, c *Client, organizationID, systemName string) (any, error) {
	const maxPages = 20
	page := 1
	for {
		query := url.Values{}
		query.Set("current", strconv.Itoa(page))
		query.Set("pageSize", "100")
		payload, err := getProjectOverviewSection(ctx, c, "systems", appstackSystemsPath(organizationID), query)
		if err != nil {
			return nil, err
		}

		pagePayload, ok := payload.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unexpected systems payload")
		}
		records, _ := pagePayload["records"].([]any)
		for _, record := range records {
			system, ok := record.(map[string]any)
			if !ok {
				continue
			}
			if name, _ := system["name"].(string); name == systemName {
				return system, nil
			}
		}

		pages, _ := pagePayload["pages"].(float64)
		if page >= int(pages) || page >= maxPages || len(records) == 0 {
			break
		}
		page++
	}
	return nil, fmt.Errorf("system %q not found in organization %q", systemName, organizationID)
}

func handleGetChangeOrderOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, changeOrderSn, err := requiredAppChangeOrder(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	coPath := appstackChangeOrderPath(organizationID, appName, changeOrderSn)

	changeOrder, err := getProjectOverviewSection(ctx, c, "changeOrder", coPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"changeOrder": changeOrder,
		"filters":     changeOrderOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeJobLogs", true) {
		if detail, ok := changeOrder.(map[string]any); ok {
			if jobs, ok := detail["jobs"]; ok {
				overview["jobs"] = jobs
			}
		}
	}

	return marshalPretty(overview)
}

func systemOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeApps":    optionalBoolDefault(params, "includeApps", true),
		"includeMembers": optionalBoolDefault(params, "includeMembers", true),
		"appLimit":       optionalIntDefault(params, "appLimit", 10),
		"memberLimit":    optionalIntDefault(params, "memberLimit", 10),
	}
}

func handleGetAppReleaseStageOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber, err := requiredAppReleaseStageExecution(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	execPath := appstackReleaseStageExecutionPath(organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber)

	overview := map[string]any{
		"filters": stageOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeStageInfo", true) {
		stagePath := appstackReleaseStagePath(organizationID, appName, releaseWorkflowSn, releaseStageSn)
		stage, err := c.GetJSON(ctx, stagePath, nil)
		if err != nil {
			return "", err
		}
		overview["stage"] = stage
	}

	if optionalBoolDefault(params, "includePipelineRun", true) {
		pipelineRun, err := c.GetJSON(ctx, execPath+":getPipelineRun", nil)
		if err != nil {
			overview["pipelineRunError"] = err.Error()
		} else {
			overview["pipelineRun"] = pipelineRun
		}
	}

	if optionalBoolDefault(params, "includeMetadata", true) {
		metadata, err := c.GetJSON(ctx, execPath+"/integratedMetadata", nil)
		if err != nil {
			overview["metadataError"] = err.Error()
		} else {
			overview["metadata"] = metadata
		}
	}

	return marshalPretty(overview)
}

func handleGetAppReleaseWorkflowOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, err := requiredAppReleaseWorkflow(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	wfPath := appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn)

	workflow, err := findAppReleaseWorkflow(ctx, c, organizationID, appName, releaseWorkflowSn)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"workflow": workflow,
		"filters":  workflowOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeStageBriefs", true) {
		briefs, err := getProjectOverviewSection(ctx, c, "stageBriefs", wfPath+"/releaseStageBriefs", nil)
		if err != nil {
			return "", err
		}
		overview["stageBriefs"] = briefs
	}

	return marshalPretty(overview)
}

func findAppReleaseWorkflow(ctx context.Context, c *Client, organizationID, appName, releaseWorkflowSn string) (any, error) {
	payload, err := getProjectOverviewSection(ctx, c, "workflows", appstackReleaseWorkflowsPath(organizationID, appName), nil)
	if err != nil {
		return nil, err
	}

	items, ok := payload.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected release workflows payload")
	}
	for _, item := range items {
		workflow, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if sn, _ := workflow["sn"].(string); sn == releaseWorkflowSn {
			return workflow, nil
		}
	}
	return nil, fmt.Errorf("release workflow %q not found in app %q", releaseWorkflowSn, appName)
}

func workflowOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeStageBriefs": optionalBoolDefault(params, "includeStageBriefs", true),
	}
}

func stageOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeStageInfo":   optionalBoolDefault(params, "includeStageInfo", true),
		"includePipelineRun": optionalBoolDefault(params, "includePipelineRun", true),
		"includeMetadata":    optionalBoolDefault(params, "includeMetadata", true),
	}
}

func changeOrderOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeJobLogs": optionalBoolDefault(params, "includeJobLogs", true),
	}
}
