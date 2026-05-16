package yunxiao

import (
	"context"
	"net/url"
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

	if err := addOrgSection(ctx, c, params, overview, orgPath, "includeDepartments", "departments", "/departments", "departmentLimit"); err != nil {
		return "", err
	}
	if err := addOrgSection(ctx, c, params, overview, orgPath, "includeMembers", "members", "/members", "memberLimit"); err != nil {
		return "", err
	}
	if err := addOrgSection(ctx, c, params, overview, orgPath, "includeGroups", "groups", "/groups", "groupLimit"); err != nil {
		return "", err
	}
	if err := addOrgSectionDirect(ctx, c, params, overview, orgPath, "includeRoles", "roles", "/roles"); err != nil {
		return "", err
	}

	return marshalPretty(overview)
}

func addOrgSection(ctx context.Context, c *Client, params map[string]any, overview map[string]any, orgPath, flagKey, key, pathSuffix, limitKey string) error {
	if !optionalBoolDefault(params, flagKey, true) {
		return nil
	}
	val, err := c.GetJSONWithMetadata(ctx, orgPath+pathSuffix, pageOneLimitQuery(params, limitKey, 5))
	if err != nil {
		return err
	}
	overview[key] = val
	return nil
}

func addOrgSectionDirect(ctx context.Context, c *Client, params map[string]any, overview map[string]any, orgPath, flagKey, key, pathSuffix string) error {
	if !optionalBoolDefault(params, flagKey, true) {
		return nil
	}
	val, err := c.GetJSON(ctx, orgPath+pathSuffix, nil)
	if err != nil {
		return err
	}
	overview[key] = val
	return nil
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
		members, err := c.GetJSONWithMetadata(ctx, groupPath+"/members", pageOneLimitQuery(params, "memberLimit", 5))
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
