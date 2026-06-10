package yunxiao

import (
	"context"
	"fmt"
	"net/url"
)

func handlePassPipelineValidate(ctx context.Context, client any, params map[string]any) (string, error) {
	return handlePipelineValidateAction(ctx, client, params, "pass")
}

func handleRefusePipelineValidate(ctx context.Context, client any, params map[string]any) (string, error) {
	return handlePipelineValidateAction(ctx, client, params, "refuse")
}

func handlePipelineValidateAction(ctx context.Context, client any, params map[string]any, action string) (string, error) {
	if action != "pass" && action != "refuse" {
		return "", &ValidationError{Msg: fmt.Sprintf("invalid action %q, must be pass or refuse", action)}
	}

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
		"/jobs/" + url.PathEscape(jobID) + "/" + action
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
