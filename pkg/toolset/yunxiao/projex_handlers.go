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
	id, err := requiredString(params, "id")
	if err != nil {
		return "", err
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/projects/" + url.PathEscape(id)
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
	spaceID, err := requiredString(params, "spaceId")
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"category": category,
		"spaceId":  spaceID,
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
	id, err := requiredString(params, "id")
	if err != nil {
		return "", err
	}

	path := "/projex/organizations/" + url.PathEscape(organizationID) + "/workitems/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}
