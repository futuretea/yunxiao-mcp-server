package yunxiao

import (
	"context"
	"net/url"
)

func handleListCommits(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	refName, err := requiredString(params, "refName")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("refName", refName)
	setOptionalString(query, params, "since")
	setOptionalString(query, params, "until")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "path")
	setOptionalString(query, params, "search")
	setOptionalBool(query, params, "showSignature")
	setOptionalString(query, params, "committerIds")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/commits"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetCommit(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	sha, err := requiredString(params, "sha")
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/commits/" + url.PathEscape(sha)
	return c.GetJSON(ctx, path, nil)
}

func handleCompare(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	from, err := requiredString(params, "from")
	if err != nil {
		return "", err
	}
	to, err := requiredString(params, "to")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("from", from)
	query.Set("to", to)
	setOptionalString(query, params, "sourceType")
	setOptionalString(query, params, "targetType")
	setOptionalString(query, params, "straight")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/compares"
	return c.GetJSON(ctx, path, query)
}
