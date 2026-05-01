package yunxiao

import (
	"context"
	"net/url"
)

func handleListMergeRequests(ctx context.Context, client any, params map[string]any) (string, error) {
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
	setOptionalStringArrayQuery(query, params, "repositoryIds")
	setOptionalStringArrayQuery(query, params, "authorUserIds")
	setOptionalStringArrayQuery(query, params, "assigneeUserIds")
	setOptionalStringArrayQuery(query, params, "subscriberUserIds")
	setOptionalString(query, params, "state")
	setOptionalString(query, params, "search")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "createdAfter")
	setOptionalString(query, params, "createdBefore")
	setOptionalString(query, params, "targetBranch")

	return c.GetJSONWithMetadata(ctx, codeupOrganizationPath(organizationID)+"/mergeRequests", query)
}

func handleGetMergeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	iid, err := requiredNumberPathString(params, "iid")
	if err != nil {
		return "", err
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/mergeRequests/" + url.PathEscape(iid)
	return c.GetJSON(ctx, path, nil)
}
