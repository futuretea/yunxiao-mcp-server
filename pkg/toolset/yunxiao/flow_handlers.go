package yunxiao

import (
	"context"
	"net/url"
)

func handleListPipelines(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "createStartTime")
	setOptionalInt(query, params, "createEndTime")
	setOptionalInt(query, params, "executeStartTime")
	setOptionalInt(query, params, "executeEndTime")
	setOptionalString(query, params, "pipelineName")
	setOptionalString(query, params, "statusList")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, flowOrganizationPath(organizationID)+"/pipelines", query)
}

func handleGetPipeline(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, flowPipelinePath(organizationID, pipelineID), nil)
}

func handleListPipelineRuns(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalInt(query, params, "startTime")
	setOptionalIntAs(query, params, "endTime", "endTme")
	setOptionalString(query, params, "status")
	setOptionalInt(query, params, "triggerMode")

	return c.GetJSONWithMetadata(ctx, flowPipelinePath(organizationID, pipelineID)+"/runs", query)
}

func handleGetPipelineRun(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	pipelineRunID, err := requiredString(params, "pipelineRunId")
	if err != nil {
		return "", err
	}

	path := flowPipelinePath(organizationID, pipelineID) + "/runs/" + url.PathEscape(pipelineRunID)
	return c.GetJSON(ctx, path, nil)
}

func handleGetLatestPipelineRun(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, flowPipelinePath(organizationID, pipelineID)+"/runs/latestPipelineRun", nil)
}

func handleListPipelineJobsByCategory(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}

	path := flowPipelinePath(organizationID, pipelineID) + "/listTasksByCategory/" + url.PathEscape(category)
	return c.GetJSON(ctx, path, nil)
}

func handleListPipelineJobHistorys(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	identifier, err := requiredString(params, "identifier")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("pipelineId", pipelineID)
	query.Set("category", category)
	query.Set("identifier", identifier)
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, flowOrganizationPath(organizationID)+"/pipelines/getComponentsWithoutButtons", query)
}

func handleGetPipelineJobRunLog(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	pipelineRunID, err := requiredString(params, "pipelineRunId")
	if err != nil {
		return "", err
	}
	jobID, err := requiredString(params, "jobId")
	if err != nil {
		return "", err
	}

	path := flowPipelinePath(organizationID, pipelineID) + "/runs/" + url.PathEscape(pipelineRunID) + "/job/" + url.PathEscape(jobID) + "/log"
	return c.GetJSON(ctx, path, nil)
}

func flowOrganizationPath(organizationID string) string {
	return "/flow/organizations/" + url.PathEscape(organizationID)
}

func flowPipelinePath(organizationID, pipelineID string) string {
	return flowOrganizationPath(organizationID) + "/pipelines/" + url.PathEscape(pipelineID)
}
