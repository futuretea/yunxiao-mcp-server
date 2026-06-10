package yunxiao

import (
	"context"
)

func setChangeRequestOptionalBody(body map[string]any, params map[string]any) {
	setOptionalStringBody(body, params, "description")
	setOptionalStringBody(body, params, "sourceProjectId")
	setOptionalStringBody(body, params, "targetProjectId")
}

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
	setChangeRequestOptionalBody(body, params)

	path := codeupRepositoryPath(organizationID, repositoryID) + "/changeRequests"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleCreateMergeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
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
	setChangeRequestOptionalBody(body, params)

	if assigneeIds, ok := params["assigneeIds"].([]any); ok && len(assigneeIds) > 0 {
		body["assigneeIds"] = assigneeIds
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/mergeRequests"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleCloseChangeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/close"
	return c.PostJSONWithMetadata(ctx, path, map[string]any{})
}

func handleReopenChangeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/reopen"
	return c.PostJSONWithMetadata(ctx, path, map[string]any{})
}

func handleMergeChangeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/merge"
	return c.PostJSONWithMetadata(ctx, path, map[string]any{})
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
