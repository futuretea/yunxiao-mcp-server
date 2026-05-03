package yunxiao

import (
	"context"
	"net/url"
)

func handleSearchProjects(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	body := map[string]any{}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "extraConditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	if body["conditions"] == nil {
		if conditions := buildProjectConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}

	return c.PostJSONWithMetadata(ctx, projexOrganizationPath(organizationID)+"/projects:search", body)
}

func handleGetProject(ctx context.Context, client any, params map[string]any) (string, error) {
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

	return c.GetJSON(ctx, projexProjectPath(organizationID, projectID), nil)
}

func handleListSprints(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "projectId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "status")
	setOptionalString(query, params, "name")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, projexProjectPath(organizationID, projectID)+"/sprints", query)
}

func handleGetSprint(ctx context.Context, client any, params map[string]any) (string, error) {
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
	sprintID, err := requiredString(params, "sprintId")
	if err != nil {
		return "", err
	}

	path := projexProjectPath(organizationID, projectID) + "/sprints/" + url.PathEscape(sprintID)
	return c.GetJSON(ctx, path, nil)
}

func handleSearchWorkitems(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	category, err := requiredString(params, "category")
	if err != nil {
		return "", err
	}
	projectID, err := requiredString(params, "projectId")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"category": category,
		"spaceId":  projectID,
	}
	setOptionalStringBody(body, params, "conditions")
	setOptionalStringBody(body, params, "orderBy")
	setOptionalStringBody(body, params, "sort")
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	if body["conditions"] == nil {
		if conditions := buildWorkitemConditions(params); conditions != "" {
			body["conditions"] = conditions
		}
	}

	return c.PostJSONWithMetadata(ctx, projexOrganizationPath(organizationID)+"/workitems:search", body)
}

func handleGetWorkitem(ctx context.Context, client any, params map[string]any) (string, error) {
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

	return c.GetJSON(ctx, projexWorkitemPath(organizationID, workitemID), nil)
}

func handleListWorkItemComments(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, workitemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, projexWorkitemPath(organizationID, workitemID)+"/comments", query)
}
