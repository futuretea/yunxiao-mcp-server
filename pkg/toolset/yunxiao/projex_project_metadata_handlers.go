package yunxiao

import (
	"context"
	"net/url"
)

func handleListProjectMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "name")
	setOptionalString(query, params, "roleId")

	return c.GetJSON(ctx, projexProjectPath(organizationID, projectID)+"/members", query)
}

func handleListProjectTemplates(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexOrganizationPath(organizationID)+"/projectTemplates", nil)
}

func handleGetProjectTemplateFieldConfig(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, templateID, err := requiredOrganizationAndNamedID(params, "projectTemplateId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexOrganizationPath(organizationID) + "/projectTemplates/" + url.PathEscape(templateID) + "/fields"
	return c.GetJSON(ctx, path, nil)
}

func handleListProjectProgram(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, programIdentifier, err := requiredOrganizationAndNamedID(params, "programIdentifier")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexOrganizationPath(organizationID) + "/" + url.PathEscape(programIdentifier) + "/binding/project/list"
	return c.GetJSON(ctx, path, nil)
}

func handleListProjectRoles(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexProjectPath(organizationID, projectID)+"/roles", nil)
}

func handleListAllProjectRoles(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, projexOrganizationPath(organizationID)+"/roles", nil)
}

func projexProjectPath(organizationID, projectID string) string {
	return projexOrganizationPath(organizationID) + "/projects/" + url.PathEscape(projectID)
}

func projexOrganizationPath(organizationID string) string {
	return "/projex/organizations/" + url.PathEscape(organizationID)
}
