package yunxiao

import (
	"context"
	"net/url"
)

func handleListPackageRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "repoTypes")
	setOptionalString(query, params, "repoCategories")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := "/packages/organizations/" + url.PathEscape(organizationID) + "/repositories"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleListArtifacts(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repoID, err := requiredOrganizationAndPackageRepo(params)
	if err != nil {
		return "", err
	}
	repoType, err := requiredString(params, "repoType")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("repoType", repoType)
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "search")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")

	path := "/packages/organizations/" + url.PathEscape(organizationID) + "/repositories/" + url.PathEscape(repoID) + "/artifacts"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetArtifact(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repoID, err := requiredOrganizationAndPackageRepo(params)
	if err != nil {
		return "", err
	}
	artifactID, err := requiredNumberPathString(params, "artifactId")
	if err != nil {
		return "", err
	}
	repoType, err := requiredString(params, "repoType")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("repoType", repoType)

	path := "/packages/organizations/" + url.PathEscape(organizationID) + "/repositories/" + url.PathEscape(repoID) + "/artifacts/" + url.PathEscape(artifactID)
	return c.GetJSON(ctx, path, query)
}

func requiredOrganizationAndPackageRepo(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	repoID, err := requiredString(params, "repoId")
	if err != nil {
		return "", "", err
	}
	return organizationID, repoID, nil
}
