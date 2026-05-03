package yunxiao

import (
	"context"
	"net/url"
)

func handleListAllWorkItemTypes(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "categories")

	return c.GetJSON(ctx, projexOrganizationPath(organizationID)+"/workitemTypes", query)
}

func handleListWorkItemTypes(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, projectID, err := requiredOrganizationAndProject(params)
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("category", category)

	return c.GetJSON(ctx, projexProjectPath(organizationID, projectID)+"/workitemTypes", query)
}

func handleGetWorkItemType(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, workItemTypeID, err := requiredOrganizationAndNamedID(params, "workItemTypeId")
	if err != nil {
		return "", err
	}

	path := projexOrganizationPath(organizationID) + "/workitemTypes/" + url.PathEscape(workItemTypeID)
	return c.GetJSON(ctx, path, nil)
}

func handleListWorkItemRelationWorkItemTypes(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, workItemTypeID, err := requiredOrganizationAndNamedID(params, "workItemTypeId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "relationType")

	path := projexOrganizationPath(organizationID) + "/workitemTypes/" + url.PathEscape(workItemTypeID) + "/relationWorkitemTypes"
	return c.GetJSON(ctx, path, query)
}

func handleGetWorkItemTypeFieldConfig(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, workItemTypeID, err := requiredOrganizationProjectAndWorkItemType(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := workItemTypeProjectPath(organizationID, projectID, workItemTypeID) + "/fields"
	return c.GetJSON(ctx, path, nil)
}

func handleGetWorkItemWorkflow(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, workItemTypeID, err := requiredOrganizationProjectAndWorkItemType(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := workItemTypeProjectPath(organizationID, projectID, workItemTypeID) + "/workflows"
	return c.GetJSON(ctx, path, nil)
}

func requiredOrganizationAndProject(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	projectID, err := requiredString(params, "projectId")
	if err != nil {
		return "", "", err
	}
	return organizationID, projectID, nil
}

func requiredOrganizationProjectAndWorkItemType(params map[string]any) (string, string, string, error) {
	organizationID, projectID, err := requiredOrganizationAndProject(params)
	if err != nil {
		return "", "", "", err
	}
	workItemTypeID, err := requiredString(params, "workItemTypeId")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, projectID, workItemTypeID, nil
}

func workItemTypeProjectPath(organizationID, projectID, workItemTypeID string) string {
	return projexProjectPath(organizationID, projectID) + "/workitemTypes/" + url.PathEscape(workItemTypeID)
}
