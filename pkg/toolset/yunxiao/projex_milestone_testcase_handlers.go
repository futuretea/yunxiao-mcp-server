package yunxiao

import (
	"context"
	"net/http"
	"net/url"
)

func handleListMilestones(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
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
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "testRepoId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexTestRepoPath(organizationID, testRepoID)+"/directories", nil)
}

func handleListTestcaseRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
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

	resp, err := c.Request(ctx, http.MethodPost, projexOrganizationPath(organizationID)+"/projects/repo/list", query, nil)
	if err != nil {
		return "", err
	}
	return PrettyResponseJSON(resp), nil
}

func handleGetTestcaseFieldConfig(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "testRepoId")
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

func handleSearchTestcases(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "testRepoId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	body := map[string]any{}
	setOptionalStringBody(body, params, "directoryId")
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	if body["conditions"] == nil {
		if conditions := buildTestcaseConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}

	path := projexTestRepoPath(organizationID, testRepoID) + "/testcases:search"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleListTestPlans(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.PostJSONWithMetadata(ctx, projexOrganizationPath(organizationID)+"/testPlan/list", nil)
}

func handleGetTestResultList(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	testPlanIdentifier, err := requiredString(params, "testPlanIdentifier")
	if err != nil {
		return "", err
	}
	directoryIdentifier, err := requiredString(params, "directoryIdentifier")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexOrganizationPath(organizationID) + "/" + url.PathEscape(testPlanIdentifier) + "/result/list/" + url.PathEscape(directoryIdentifier)
	return c.PostJSONWithMetadata(ctx, path, nil)
}

func buildTestcaseConditions(params map[string]any) string {
	if subject, _ := params["subject"].(string); subject != "" {
		return marshalConditions([]map[string]any{stringContainsCondition("subject", subject)})
	}
	return ""
}

func requiredOrganizationTestRepoAndTestcase(params map[string]any) (string, string, string, error) {
	organizationID, testRepoID, err := requiredOrganizationAndNamedID(params, "testRepoId")
	if err != nil {
		return "", "", "", err
	}
	testcaseID, err := requiredString(params, "testcaseId")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, testRepoID, testcaseID, nil
}

func projexTestRepoPath(organizationID, testRepoID string) string {
	return projexOrganizationPath(organizationID) + "/testRepos/" + url.PathEscape(testRepoID)
}
