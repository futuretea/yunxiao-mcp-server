package yunxiao

import (
	"context"
	"net/url"
)

func handlePassPipelineValidate(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, pipelineID, pipelineRunID, jobID, err := requiredOrganizationPipelineRunAndJob(params)
	if err != nil {
		return "", err
	}
	path := "/flow/organizations/" + url.PathEscape(organizationID) +
		"/pipelines/" + url.PathEscape(pipelineID) +
		"/pipelineRuns/" + url.PathEscape(pipelineRunID) +
		"/jobs/" + url.PathEscape(jobID) + "/pass"
	return c.PostJSONWithMetadata(ctx, path, map[string]any{})
}

func handleRefusePipelineValidate(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, pipelineID, pipelineRunID, jobID, err := requiredOrganizationPipelineRunAndJob(params)
	if err != nil {
		return "", err
	}
	path := "/flow/organizations/" + url.PathEscape(organizationID) +
		"/pipelines/" + url.PathEscape(pipelineID) +
		"/pipelineRuns/" + url.PathEscape(pipelineRunID) +
		"/jobs/" + url.PathEscape(jobID) + "/refuse"
	return c.PostJSONWithMetadata(ctx, path, map[string]any{})
}

func requiredOrganizationPipelineRunAndJob(params map[string]any) (orgID, pipelineID, pipelineRunID, jobID string, err error) {
	orgID, err = requiredString(params, "organizationId")
	if err != nil {
		return "", "", "", "", err
	}
	pipelineID, err = requiredString(params, "pipelineId")
	if err != nil {
		return "", "", "", "", err
	}
	pipelineRunID, err = requiredString(params, "pipelineRunId")
	if err != nil {
		return "", "", "", "", err
	}
	jobID, err = requiredString(params, "jobId")
	if err != nil {
		return "", "", "", "", err
	}
	return orgID, pipelineID, pipelineRunID, jobID, nil
}
