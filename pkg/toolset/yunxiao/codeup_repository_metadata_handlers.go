package yunxiao

import (
	"context"
	"net/url"
)

func handleListTags(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
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
	setOptionalString(query, params, "search")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "orderBy")

	return c.GetJSONWithMetadata(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/tags", query)
}

func handleListRepositoryMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "accessLevel")

	return c.GetJSON(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/members", query)
}

func handleListProtectedBranches(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/protectedBranches", nil)
}

func handleGetProtectedBranch(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	id, err := requiredNumberPathString(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/protectedBranches/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}

func handleListPushRules(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/pushRules", nil)
}

func handleGetPushRule(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	pushRuleID, err := requiredNumberPathString(params, "pushRuleId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/pushRules/" + url.PathEscape(pushRuleID)
	return c.GetJSON(ctx, path, nil)
}

func codeupRepositoryPath(organizationID, repositoryID string) string {
	return "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID)
}
