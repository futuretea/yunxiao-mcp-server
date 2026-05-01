package yunxiao

import (
	"context"
	"net/url"
)

func handleListGroupMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, groupID, err := requiredOrganizationAndNamedID(params, "groupId")
	if err != nil {
		return "", err
	}
	query := url.Values{}
	setOptionalInt(query, params, "accessLevel")

	path := codeupOrganizationPath(organizationID) + "/groups/" + encodePathValue(groupID) + "/members"
	return c.GetJSON(ctx, path, query)
}

func handleGetMemberHTTPSCloneUsername(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, userID, err := requiredOrganizationAndNamedID(params, "userId")
	if err != nil {
		return "", err
	}

	path := codeupOrganizationPath(organizationID) + "/users/" + encodePathValue(userID) + "/httpsCloneUsername"
	return c.GetJSON(ctx, path, nil)
}
