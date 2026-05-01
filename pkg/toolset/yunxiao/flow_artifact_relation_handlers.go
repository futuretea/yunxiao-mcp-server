package yunxiao

import (
	"context"
	"net/url"
)

func handleGetPipelineScanReportURL(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	reportPath, err := requiredString(params, "reportPath")
	if err != nil {
		return "", err
	}

	query := url.Values{"reportPath": []string{reportPath}}
	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/getPipelineScanReportUrl"
	return c.GetJSON(ctx, path, query)
}

func handleGetPipelineArtifactURL(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	filePath, err := requiredString(params, "filePath")
	if err != nil {
		return "", err
	}
	fileName, err := requiredString(params, "fileName")
	if err != nil {
		return "", err
	}

	query := url.Values{"filePath": []string{filePath}, "fileName": []string{fileName}}
	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/getArtifactDownloadUrl"
	return c.GetJSON(ctx, path, query)
}

func handleGetPipelineEmasArtifactURL(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, query, err := requiredPipelineEmasArtifactURLQuery(params)
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/getEmasArtifactDownloadUrl"
	return c.GetJSON(ctx, path, query)
}

func requiredPipelineEmasArtifactURLQuery(params map[string]any) (string, url.Values, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", nil, err
	}
	emasJobInstanceID, err := requiredString(params, "emasJobInstanceId")
	if err != nil {
		return "", nil, err
	}
	md5, err := requiredString(params, "md5")
	if err != nil {
		return "", nil, err
	}
	pipelineID, err := requiredNumberPathString(params, "pipelineId")
	if err != nil {
		return "", nil, err
	}
	pipelineRunID, err := requiredNumberPathString(params, "pipelineRunId")
	if err != nil {
		return "", nil, err
	}
	serviceConnectionID, err := requiredNumberPathString(params, "serviceConnectionId")
	if err != nil {
		return "", nil, err
	}

	query := url.Values{
		"emasJobInstanceId":   []string{emasJobInstanceID},
		"md5":                 []string{md5},
		"pipelineId":          []string{pipelineID},
		"pipelineRunId":       []string{pipelineRunID},
		"serviceConnectionId": []string{serviceConnectionID},
	}
	return organizationID, query, nil
}

func handleListPipelineRelations(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}
	relObjectType, err := requiredString(params, "relObjectType")
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/pipelines/" + url.PathEscape(pipelineID) + "/pipelineObjRel/" + url.PathEscape(relObjectType) + "/list"
	return c.GetJSON(ctx, path, nil)
}

func handleGetLastInstance(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, pipelineID, err := requiredOrganizationAndPipeline(params)
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/createServiceConnection/" + url.PathEscape(pipelineID) + "/getLastInstance"
	return c.GetJSON(ctx, path, nil)
}
