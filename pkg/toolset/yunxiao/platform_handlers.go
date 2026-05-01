package yunxiao

import (
	"context"
	"net/url"
)

func handleListOrganizationDepartments(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "parentId")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	path := organizationPath(organizationID) + "/departments"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetOrganizationDepartmentInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	return getOrganizationDepartment(ctx, client, params, "")
}

func handleGetOrganizationDepartmentAncestors(ctx context.Context, client any, params map[string]any) (string, error) {
	return getOrganizationDepartment(ctx, client, params, "/ancestors")
}

func getOrganizationDepartment(ctx context.Context, client any, params map[string]any, suffix string) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, id, err := requiredOrganizationAndID(params)
	if err != nil {
		return "", err
	}

	path := organizationPath(organizationID) + "/departments/" + url.PathEscape(id) + suffix
	return c.GetJSON(ctx, path, nil)
}

func handleListOrganizationMembers(ctx context.Context, client any, params map[string]any) (string, error) {
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

	path := organizationPath(organizationID) + "/members"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetOrganizationMemberInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, id, err := requiredOrganizationAndID(params)
	if err != nil {
		var memberID string
		organizationID, memberID, err = requiredOrganizationAndNamedID(params, "memberId")
		if err != nil {
			return "", err
		}
		id = memberID
	}

	path := organizationPath(organizationID) + "/members/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}

func handleGetOrganizationMemberInfoByUserID(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	userID, err := requiredString(params, "userId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("userId", userID)

	path := organizationPath(organizationID) + "/members:readByUser"
	return c.GetJSON(ctx, path, query)
}

func handleSearchOrganizationMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	body := map[string]any{}
	setOptionalIntBody(body, params, "page")
	setOptionalIntBody(body, params, "perPage")
	setOptionalBoolBody(body, params, "includeChildren")
	setOptionalStringBody(body, params, "nextToken")
	setOptionalStringBody(body, params, "query")
	setOptionalStringArrayBody(body, params, "deptIds")
	setOptionalStringArrayBody(body, params, "roleIds")
	setOptionalStringArrayBody(body, params, "statuses")

	path := organizationPath(organizationID) + "/members:search"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleListOrganizationRoles(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	path := organizationPath(organizationID) + "/roles"
	return c.GetJSON(ctx, path, nil)
}

func handleGetOrganizationRole(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, id, err := requiredOrganizationAndID(params)
	if err != nil {
		var roleID string
		organizationID, roleID, err = requiredOrganizationAndNamedID(params, "roleId")
		if err != nil {
			return "", err
		}
		id = roleID
	}

	path := organizationPath(organizationID) + "/roles/" + url.PathEscape(id)
	return c.GetJSON(ctx, path, nil)
}

func handleListUsers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "filter")
	setOptionalString(query, params, "status")
	setOptionalString(query, params, "deptId")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, "/platform/users", query)
}

func organizationPath(organizationID string) string {
	return "/platform/organizations/" + url.PathEscape(organizationID)
}
