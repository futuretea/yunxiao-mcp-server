package yunxiao

import (
	"context"
	"net/url"
)

func handleGetPipelineOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	pipelinePath := flowPipelinePath(organizationID, pipelineID)

	pipeline, err := getProjectOverviewSection(ctx, c, "pipeline", pipelinePath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"pipeline": pipeline,
		"filters":  pipelineOverviewFilters(params),
	}

	latestRun, err := getProjectOverviewSection(ctx, c, "latestRun", pipelinePath+"/runs/latestPipelineRun", nil)
	if err != nil {
		return "", err
	}
	overview["latestRun"] = latestRun

	if optionalBoolDefault(params, "includeRuns", true) {
		runs, err := getProjectOverviewSection(ctx, c, "runs", pipelinePath+"/runs", pageOneLimitQuery(params, "runLimit", 5))
		if err != nil {
			return "", err
		}
		overview["runs"] = runs
	}

	return marshalPretty(overview)
}

func pipelineOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeRuns": optionalBoolDefault(params, "includeRuns", true),
		"runLimit":    optionalIntDefault(params, "runLimit", 5),
	}
}

func handleGetPipelineRunOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	pipelineRunID, err := requiredString(params, "pipelineRunId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	pipelinePath := flowPipelinePath(organizationID, pipelineID)
	runPath := pipelinePath + "/runs/" + url.PathEscape(pipelineRunID)

	run, err := getProjectOverviewSection(ctx, c, "run", runPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"run":     run,
		"filters": pipelineRunOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeJobs", true) {
		category := optionalStringDefault(params, "category", "DEPLOY")
		jobs, err := getProjectOverviewSection(ctx, c, "jobs", pipelinePath+"/listTasksByCategory/"+url.PathEscape(category), nil)
		if err != nil {
			return "", err
		}
		overview["jobs"] = jobs
	}

	return marshalPretty(overview)
}

func pipelineRunOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeJobs": optionalBoolDefault(params, "includeJobs", true),
		"category":    optionalStringDefault(params, "category", "DEPLOY"),
	}
}
