package yunxiao

import (
	"context"
	"fmt"
)

func handleCreateWorkitem(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	projectID, err := requiredString(params, "projectId")
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	workitemTypeID, err := requiredString(params, "workitemTypeId")
	if err != nil {
		return "", err
	}
	subject, err := requiredString(params, "subject")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"spaceId":        projectID,
		"category":       category,
		"workitemTypeId": workitemTypeID,
		"subject":        subject,
	}
	setOptionalStringBody(body, params, "description")
	setOptionalStringBody(body, params, "assignedTo")
	setOptionalStringBody(body, params, "priority")
	setOptionalStringBody(body, params, "parentId")
	setOptionalStringBody(body, params, "sprint")
	setOptionalStringArrayBody(body, params, "labels")

	path := projexOrganizationPath(organizationID) + "/workitems"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleUpdateWorkitem(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	workitemID, err := requiredString(params, "workitemId")
	if err != nil {
		return "", err
	}

	body := map[string]any{}
	setOptionalStringBody(body, params, "subject")
	setOptionalStringBody(body, params, "description")
	setOptionalStringBody(body, params, "assignedTo")
	setOptionalStringBody(body, params, "priority")
	setOptionalStringBody(body, params, "sprint")
	setOptionalStringArrayBody(body, params, "labels")

	if len(body) == 0 {
		return "", fmt.Errorf("at least one field to update must be provided")
	}

	return c.PutJSONWithMetadata(ctx, projexWorkitemPath(organizationID, workitemID), body)
}

func handleUpdateWorkitemStatus(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	workitemID, err := requiredString(params, "workitemId")
	if err != nil {
		return "", err
	}
	statusID, err := requiredString(params, "statusId")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"statusId": statusID,
	}

	result, err := c.PostJSONWithMetadata(ctx, projexWorkitemPath(organizationID, workitemID)+"/status", body)
	if err != nil {
		return "", err
	}

	if comment, _ := params["comment"].(string); comment != "" {
		commentBody := map[string]any{"content": comment}
		_, _ = c.PostJSONWithMetadata(ctx, projexWorkitemPath(organizationID, workitemID)+"/comments", commentBody)
	}

	return result, nil
}

func handleAddWorkitemComment(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	workitemID, err := requiredString(params, "workitemId")
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

	return c.PostJSONWithMetadata(ctx, projexWorkitemPath(organizationID, workitemID)+"/comments", body)
}
