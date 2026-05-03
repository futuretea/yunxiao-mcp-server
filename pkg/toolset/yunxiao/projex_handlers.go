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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects:search"
	return c.PostJSONWithMetadata(ctx, path, body)
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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects/" + url.PathEscape(projectID)
	return c.GetJSON(ctx, path, nil)
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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects/" + url.PathEscape(projectID) + "/sprints"
	return c.GetJSONWithMetadata(ctx, path, query)
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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects/" + url.PathEscape(projectID) + "/sprints/" + url.PathEscape(sprintID)
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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems:search"
	return c.PostJSONWithMetadata(ctx, path, body)
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

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems/" + url.PathEscape(workitemID)
	return c.GetJSON(ctx, path, nil)
}

func handleListWorkItemComments(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workItemId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems/" + url.PathEscape(workItemID) + "/comments"
	return c.GetJSONWithMetadata(ctx, path, query)
}
