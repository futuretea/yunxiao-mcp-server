package yunxiao

import (
	"context"
)

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
