package yunxiao

import (
	"context"
	"net/url"
)

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
