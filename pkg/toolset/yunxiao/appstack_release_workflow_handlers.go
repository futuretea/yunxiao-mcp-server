package yunxiao

import (
	"context"
	"net/url"
)

func handleListAppReleaseWorkflows(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/releaseWorkflows"
	return c.GetJSON(ctx, path, nil)
}

func handleListAppReleaseWorkflowBriefs(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackAppPath(organizationID, appName) + "/releaseWorkflowBriefs"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppReleaseWorkflowStage(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, err := requiredAppReleaseStage(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn) + "/releaseStage/" + url.PathEscape(releaseStageSn)
	return c.GetJSON(ctx, path, nil)
}

func handleListAppReleaseStageBriefs(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, err := requiredAppReleaseWorkflow(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn) + "/releaseStageBriefs"
	return c.GetJSON(ctx, path, nil)
}

func handleListAppReleaseStageRuns(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, err := requiredAppReleaseStage(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "pagination")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "nextToken")
	setOptionalInt(query, params, "page")

	path := appstackReleaseStageResourcePath(organizationID, appName, releaseWorkflowSn, releaseStageSn) + "/executions"
	return c.GetJSON(ctx, path, query)
}

func handleListAppReleaseStageExecMetadata(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber, err := requiredAppReleaseStageExecution(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseStageExecutionPath(organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber) + "/integratedMetadata"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppReleaseStagePipelineRun(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber, err := requiredAppReleaseStageExecution(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackReleaseStageExecutionPath(organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber) + ":getPipelineRun"
	return c.GetJSON(ctx, path, nil)
}

func handleGetAppReleaseStagePipelineJobLog(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber, err := requiredAppReleaseStageExecution(params)
	if err != nil {
		return "", err
	}
	jobID, err := requiredString(params, "jobId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("jobId", jobID)

	path := appstackReleaseStageExecutionPath(organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber) + ":pipelineJobLog"
	return c.GetJSON(ctx, path, query)
}

func requiredAppReleaseWorkflow(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	releaseWorkflowSn, err := requiredString(params, "releaseWorkflowSn")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, releaseWorkflowSn, nil
}

func requiredAppReleaseStage(params map[string]any) (string, string, string, string, error) {
	organizationID, appName, releaseWorkflowSn, err := requiredAppReleaseWorkflow(params)
	if err != nil {
		return "", "", "", "", err
	}
	releaseStageSn, err := requiredString(params, "releaseStageSn")
	if err != nil {
		return "", "", "", "", err
	}
	return organizationID, appName, releaseWorkflowSn, releaseStageSn, nil
}

func requiredAppReleaseStageExecution(params map[string]any) (string, string, string, string, string, error) {
	organizationID, appName, releaseWorkflowSn, releaseStageSn, err := requiredAppReleaseStage(params)
	if err != nil {
		return "", "", "", "", "", err
	}
	executionNumber, err := requiredString(params, "executionNumber")
	if err != nil {
		return "", "", "", "", "", err
	}
	return organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber, nil
}

func appstackReleaseWorkflowPath(organizationID, appName, releaseWorkflowSn string) string {
	return appstackAppPath(organizationID, appName) + "/releaseWorkflow/" + url.PathEscape(releaseWorkflowSn)
}

func appstackReleaseWorkflowResourcePath(organizationID, appName, releaseWorkflowSn string) string {
	return appstackAppPath(organizationID, appName) + "/releaseWorkflows/" + url.PathEscape(releaseWorkflowSn)
}

func appstackReleaseStageResourcePath(organizationID, appName, releaseWorkflowSn, releaseStageSn string) string {
	return appstackReleaseWorkflowResourcePath(organizationID, appName, releaseWorkflowSn) + "/releaseStages/" + url.PathEscape(releaseStageSn)
}

func appstackReleaseStageExecutionPath(organizationID, appName, releaseWorkflowSn, releaseStageSn, executionNumber string) string {
	return appstackReleaseStageResourcePath(organizationID, appName, releaseWorkflowSn, releaseStageSn) + "/executions/" + url.PathEscape(executionNumber)
}
