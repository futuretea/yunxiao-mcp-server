package yunxiao

import (
	"context"
	"net/url"
)

func handleListOrganizationGroups(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, organizationPath(organizationID)+"/groups", query)
}

func handleGetOrganizationGroup(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, groupID, err := requiredOrganizationAndNamedID(params, "groupId")
	if err != nil {
		return "", err
	}

	path := organizationPath(organizationID) + "/groups/" + encodePathValue(groupID)
	return c.GetJSON(ctx, path, nil)
}

func handleListOrganizationGroupMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, groupID, err := requiredOrganizationAndNamedID(params, "groupId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := organizationPath(organizationID) + "/groups/" + encodePathValue(groupID) + "/members"
	return c.GetJSONWithMetadata(ctx, path, query)
}
