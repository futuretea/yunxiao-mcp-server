package yunxiao

import (
	"context"
	"net/url"
)

func handleListResourceMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	resourceType, err := requiredString(params, "resourceType")
	if err != nil {
		return "", err
	}
	resourceID, err := requiredString(params, "resourceId")
	if err != nil {
		return "", err
	}

	path := "/flow/organizations/" + url.PathEscape(organizationID) + "/resourceMembers/resourceTypes/" + url.PathEscape(resourceType) + "/resourceIds/" + url.PathEscape(resourceID)
	return c.GetJSON(ctx, path, nil)
}
