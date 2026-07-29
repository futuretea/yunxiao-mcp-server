package yunxiao

import (
	"context"
	"net/http"
	"net/url"

	sdk "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"
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
	// API field name is intentionally "endTme" (not "endTime").
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

func handleGetPipelineJobSteps(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, pipelineID, pipelineRunID, jobID, err := requiredPipelineRunJob(params)
	if err != nil {
		return "", err
	}

	path := flowPipelineRunJobPath(organizationID, pipelineID, pipelineRunID, jobID) + "/steps"
	return getJSONPreservingNumbers(ctx, c, path, nil)
}

func handleGetPipelineJobStepLog(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, pipelineID, pipelineRunID, jobID, err := requiredPipelineRunJob(params)
	if err != nil {
		return "", err
	}
	query, err := requiredPipelineJobStepLogQuery(params, true)
	if err != nil {
		return "", err
	}

	path := flowPipelineRunJobPath(organizationID, pipelineID, pipelineRunID, jobID) + "/step/log"
	return getJSONPreservingNumbers(ctx, c, path, query)
}

func handleGetPipelineJobStepLogURL(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, pipelineID, pipelineRunID, jobID, err := requiredPipelineRunJob(params)
	if err != nil {
		return "", err
	}
	query, err := requiredPipelineJobStepLogQuery(params, false)
	if err != nil {
		return "", err
	}

	path := flowPipelineRunJobPath(organizationID, pipelineID, pipelineRunID, jobID) + "/step/log/url"
	return getJSONPreservingNumbers(ctx, c, path, query)
}

func requiredPipelineRunJob(params map[string]any) (string, string, string, string, error) {
	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", "", "", "", err
	}
	pipelineRunID, err := requiredString(params, "pipelineRunId")
	if err != nil {
		return "", "", "", "", err
	}
	jobID, err := requiredString(params, "jobId")
	if err != nil {
		return "", "", "", "", err
	}
	return organizationID, pipelineID, pipelineRunID, jobID, nil
}

func requiredPipelineJobStepLogQuery(params map[string]any, withPagination bool) (url.Values, error) {
	stepIndex, err := requiredStrictIntegerString(params, "stepIndex")
	if err != nil {
		return nil, err
	}
	buildID, err := requiredDecimalString(params, "buildId")
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("stepIndex", stepIndex)
	if withPagination {
		offset, err := requiredStrictIntegerString(params, "offset")
		if err != nil {
			return nil, err
		}
		limit, err := requiredStrictIntegerString(params, "limit")
		if err != nil {
			return nil, err
		}
		query.Set("offset", offset)
		query.Set("limit", limit)
	}
	query.Set("buildId", buildID)
	return query, nil
}

func getJSONPreservingNumbers(ctx context.Context, c *Client, path string, query url.Values) (string, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return "", sdk.WrapError(sdk.FriendlyAPIError(err))
	}
	return prettyJSONPreservingNumbers(resp.Body), nil
}

func flowOrganizationPath(organizationID string) string {
	return "/flow/organizations/" + url.PathEscape(organizationID)
}

func flowPipelinePath(organizationID, pipelineID string) string {
	return flowOrganizationPath(organizationID) + "/pipelines/" + url.PathEscape(pipelineID)
}

func flowPipelineRunJobPath(organizationID, pipelineID, pipelineRunID, jobID string) string {
	return flowPipelinePath(organizationID, pipelineID) + "/pipelineRuns/" + url.PathEscape(pipelineRunID) + "/jobs/" + url.PathEscape(jobID)
}
