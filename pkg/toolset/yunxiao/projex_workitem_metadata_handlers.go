package yunxiao

import (
	"context"
	"net/url"
)

func handleListWorkitemAttachments(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexWorkitemPath(organizationID, workItemID)+"/attachments", nil)
}

func handleGetWorkitemFile(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, fileID, err := requiredOrganizationWorkitemAndFile(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexWorkitemPath(organizationID, workItemID) + "/files/" + url.PathEscape(fileID)
	return c.GetJSON(ctx, path, nil)
}

func handleListWorkitemRelationRecords(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", err
	}
	relationType, err := requiredString(params, "relationType")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("relationType", relationType)

	return c.GetJSON(ctx, projexWorkitemPath(organizationID, workItemID)+"/relationRecords", query)
}

func handleListLabels(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
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

	return c.GetJSONWithMetadata(ctx, projexProjectPath(organizationID, projectID)+"/labels", query)
}

func requiredOrganizationWorkitemAndFile(params map[string]any) (string, string, string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", "", "", err
	}
	fileID, err := requiredString(params, "fileId")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, workItemID, fileID, nil
}
