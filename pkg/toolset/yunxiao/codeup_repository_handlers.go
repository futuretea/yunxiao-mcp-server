package yunxiao

import (
	"context"
	"net/url"
)

func handleListRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")
	setOptionalBool(query, params, "archived")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetRepository(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID)
	return c.GetJSON(ctx, path, nil)
}

func handleListBranches(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/branches"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetBranch(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	branchName, err := requiredString(params, "branchName")
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/branches/" + encodePathValue(branchName)
	return c.GetJSON(ctx, path, nil)
}

func handleListFiles(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "path")
	setOptionalString(query, params, "ref")
	setOptionalString(query, params, "type")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/files/tree"
	return c.GetJSON(ctx, path, query)
}

func handleGetFileBlobs(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	filePath, err := requiredString(params, "filePath")
	if err != nil {
		return "", err
	}
	ref, err := requiredString(params, "ref")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("ref", ref)

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/files/" + encodeFilePath(filePath)
	return c.GetJSON(ctx, path, query)
}
