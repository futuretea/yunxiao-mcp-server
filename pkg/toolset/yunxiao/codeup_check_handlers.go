package yunxiao

import (
	"context"
	"net/url"
)

func handleListCommitStatuses(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	sha, err := requiredString(params, "sha")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := codeupRepositoryPath(organizationID, repositoryID) + "/commits/" + url.PathEscape(sha) + "/statuses"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleListCheckRuns(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	ref, err := requiredString(params, "ref")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("ref", ref)
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/checkRuns", query)
}

func handleGetCheckRun(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	checkRunID, err := requiredNumberPathString(params, "checkRunId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/checkRuns/" + url.PathEscape(checkRunID)
	return c.GetJSON(ctx, path, nil)
}
