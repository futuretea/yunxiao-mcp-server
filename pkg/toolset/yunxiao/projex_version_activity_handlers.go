package yunxiao

import (
	"context"
	"net/url"
)

func handleListVersions(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, projectID, err := requiredOrganizationAndNamedID(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "status")
	setOptionalString(query, params, "name")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, projexProjectPath(organizationID, projectID)+"/versions", query)
}

func handleListWorkitemActivities(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "id")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexOrganizationPath(organizationID) + "/workitems/" + url.PathEscape(workItemID) + "/activities"
	return c.GetJSON(ctx, path, nil)
}
