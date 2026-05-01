package yunxiao

import (
	"context"
	"net/url"
)

func handleListMilestones(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "status")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, projexProjectPath(organizationID, projectID)+"/milestones", query)
}

func handleListDirectories(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexTestRepoPath(organizationID, testRepoID)+"/directories", nil)
}

func handleGetTestcaseFieldConfig(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexTestRepoPath(organizationID, testRepoID)+"/testcases/fields", nil)
}

func handleGetTestcase(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, testRepoID, testcaseID, err := requiredOrganizationTestRepoAndTestcase(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexTestRepoPath(organizationID, testRepoID) + "/testcases/" + url.PathEscape(testcaseID)
	return c.GetJSON(ctx, path, nil)
}

func requiredOrganizationTestRepoAndTestcase(params map[string]any) (string, string, string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "testRepoId")
	if err != nil {
		return "", "", "", err
	}
	testcaseID, err := requiredString(params, "id")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, testRepoID, testcaseID, nil
}

func projexTestRepoPath(organizationID, testRepoID string) string {
	return projexOrganizationPath(organizationID) + "/testRepos/" + url.PathEscape(testRepoID)
}
