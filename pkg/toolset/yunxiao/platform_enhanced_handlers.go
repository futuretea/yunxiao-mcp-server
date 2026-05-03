package yunxiao

import (
	"context"
	"net/url"
	"strconv"
)

func handleGetOrganizationOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	orgPath := organizationPath(organizationID)

	org, err := c.GetJSON(ctx, orgPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"organization": org,
		"filters":      organizationOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeDepartments", true) {
		deptQuery := url.Values{}
		deptQuery.Set("page", "1")
		deptQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "departmentLimit", 5)))
		depts, err := c.GetJSONWithMetadata(ctx, orgPath+"/departments", deptQuery)
		if err != nil {
			return "", err
		}
		overview["departments"] = depts
	}

	if optionalBoolDefault(params, "includeMembers", true) {
		memberQuery := url.Values{}
		memberQuery.Set("page", "1")
		memberQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "memberLimit", 5)))
		members, err := c.GetJSONWithMetadata(ctx, orgPath+"/members", memberQuery)
		if err != nil {
			return "", err
		}
		overview["members"] = members
	}

	if optionalBoolDefault(params, "includeGroups", true) {
		groupQuery := url.Values{}
		groupQuery.Set("page", "1")
		groupQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "groupLimit", 5)))
		groups, err := c.GetJSONWithMetadata(ctx, orgPath+"/groups", groupQuery)
		if err != nil {
			return "", err
		}
		overview["groups"] = groups
	}

	if optionalBoolDefault(params, "includeRoles", true) {
		roles, err := c.GetJSON(ctx, orgPath+"/roles", nil)
		if err != nil {
			return "", err
		}
		overview["roles"] = roles
	}

	return marshalPretty(overview)
}

func organizationOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeDepartments": optionalBoolDefault(params, "includeDepartments", true),
		"includeMembers":     optionalBoolDefault(params, "includeMembers", true),
		"includeGroups":      optionalBoolDefault(params, "includeGroups", true),
		"includeRoles":       optionalBoolDefault(params, "includeRoles", true),
		"departmentLimit":    optionalIntDefault(params, "departmentLimit", 5),
		"memberLimit":        optionalIntDefault(params, "memberLimit", 5),
		"groupLimit":         optionalIntDefault(params, "groupLimit", 5),
	}
}

func handleGetOrganizationDepartmentOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, departmentID, err := requiredOrganizationAndNamedID(params, "departmentId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	deptPath := organizationPath(organizationID) + "/departments/" + url.PathEscape(departmentID)

	dept, err := c.GetJSON(ctx, deptPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"department": dept,
		"filters":    organizationDepartmentOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeAncestors", true) {
		ancestors, err := c.GetJSON(ctx, deptPath+"/ancestors", nil)
		if err != nil {
			return "", err
		}
		overview["ancestors"] = ancestors
	}

	return marshalPretty(overview)
}

func organizationDepartmentOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeAncestors": optionalBoolDefault(params, "includeAncestors", true),
	}
}

func handleGetOrganizationGroupOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, groupID, err := requiredOrganizationAndNamedID(params, "groupId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	groupPath := organizationPath(organizationID) + "/groups/" + encodePathValue(groupID)

	group, err := c.GetJSON(ctx, groupPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"group":   group,
		"filters": organizationGroupOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeMembers", true) {
		memberQuery := url.Values{}
		memberQuery.Set("page", "1")
		memberQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "memberLimit", 5)))
		members, err := c.GetJSONWithMetadata(ctx, groupPath+"/members", memberQuery)
		if err != nil {
			return "", err
		}
		overview["members"] = members
	}

	return marshalPretty(overview)
}

func organizationGroupOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeMembers": optionalBoolDefault(params, "includeMembers", true),
		"memberLimit":    optionalIntDefault(params, "memberLimit", 5),
	}
}
