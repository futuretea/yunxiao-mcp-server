package yunxiao

import (
	"context"
	"net/url"
)

func handleCreateChangeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	title, err := requiredString(params, "title")
	if err != nil {
		return "", err
	}
	sourceBranch, err := requiredString(params, "sourceBranch")
	if err != nil {
		return "", err
	}
	targetBranch, err := requiredString(params, "targetBranch")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"title":        title,
		"sourceBranch": sourceBranch,
		"targetBranch": targetBranch,
	}
	setOptionalStringBody(body, params, "description")
	setOptionalStringBody(body, params, "sourceProjectId")
	setOptionalStringBody(body, params, "targetProjectId")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/changeRequests"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleAddChangeRequestComment(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	localID, err := requiredString(params, "localId")
	if err != nil {
		return "", err
	}
	content, err := requiredString(params, "content")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"content": content,
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/comments"
	return c.PostJSONWithMetadata(ctx, path, body)
}
