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

	path := packagesOrganizationPath(organizationID) + "/repositories"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleListArtifacts(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repoID, err := requiredOrganizationAndNamedID(params, "repoId")
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

	path := packagesRepoPath(organizationID, repoID) + "/artifacts"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetArtifact(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repoID, err := requiredOrganizationAndNamedID(params, "repoId")
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

	path := packagesRepoPath(organizationID, repoID) + "/artifacts/" + encodePathValue(artifactID)
	return c.GetJSON(ctx, path, query)
}

func packagesOrganizationPath(organizationID string) string {
	return "/packages/organizations/" + encodePathValue(organizationID)
}

func packagesRepoPath(organizationID, repoID string) string {
	return packagesOrganizationPath(organizationID) + "/repositories/" + encodePathValue(repoID)
}
