package yunxiao

import (
	"context"
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

	system, err := c.GetJSON(ctx, systemPath, nil)
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

	changeOrder, err := c.GetJSON(ctx, coPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"changeOrder": changeOrder,
		"filters":     changeOrderOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeJobLogs", true) {
		jobs, err := getProjectOverviewSection(ctx, c, "jobs", coPath+"/jobs", nil)
		if err != nil {
			overview["jobLogsError"] = err.Error()
		} else {
			overview["jobs"] = jobs
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
		stagePath := appstackReleaseStageResourcePath(organizationID, appName, releaseWorkflowSn, releaseStageSn)
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
