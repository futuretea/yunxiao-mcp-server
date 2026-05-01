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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines"
	return c.GetJSONWithMetadata(ctx, path, query)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID)
	return c.GetJSON(ctx, path, nil)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs"
	return c.GetJSONWithMetadata(ctx, path, query)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs/" + url.PathEscape(pipelineRunID)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs/latestPipelineRun"
	return c.GetJSON(ctx, path, nil)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/listTasksByCategory/" + url.PathEscape(category)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/getComponentsWithoutButtons"
	return c.GetJSONWithMetadata(ctx, path, query)
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

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/runs/" + url.PathEscape(pipelineRunID) + "/job/" + url.PathEscape(jobID) + "/log"
	return c.GetJSON(ctx, path, nil)
}
