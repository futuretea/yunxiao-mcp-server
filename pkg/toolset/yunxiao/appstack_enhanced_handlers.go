package yunxiao

import (
	"context"
	"net/url"
	"strconv"
)

func handleGetApplicationOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	appPath := appstackAppPath(organizationID, appName)

	app, err := getProjectOverviewSection(ctx, c, "application", appPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"application": app,
		"filters":     applicationOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeEnvironments", true) {
		envQuery := url.Values{}
		envQuery.Set("page", "1")
		envQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "envLimit", 5)))
		envs, err := getProjectOverviewSection(ctx, c, "environments", appPath+"/envs", envQuery)
		if err != nil {
			return "", err
		}
		overview["environments"] = envs
	}

	if optionalBoolDefault(params, "includeOrchestrations", true) {
		orchQuery := url.Values{}
		orchQuery.Set("page", "1")
		orchQuery.Set("perPage", strconv.Itoa(optionalIntDefault(params, "orchestrationLimit", 5)))
		orchs, err := getProjectOverviewSection(ctx, c, "orchestrations", appPath+"/orchestrations", orchQuery)
		if err != nil {
			return "", err
		}
		overview["orchestrations"] = orchs
	}

	return marshalPretty(overview)
}

func handleGetEnvironmentOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, envName, err := requiredAppEnvironment(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	envPath := appstackAppPath(organizationID, appName) + "/envs/" + url.PathEscape(envName)

	env, err := c.GetJSON(ctx, envPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"environment": env,
		"filters":     environmentOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeVariableGroups", true) {
		vgs, err := c.GetJSON(ctx, envPath+"/variableGroups", nil)
		if err != nil {
			return "", err
		}
		overview["variableGroups"] = vgs
	}

	if optionalBoolDefault(params, "includeLatestOrchestration", true) {
		orch, err := c.GetJSON(ctx, envPath+"/orchestration:latestAvailable", nil)
		if err != nil {
			return "", err
		}
		overview["latestOrchestration"] = orch
	}

	return marshalPretty(overview)
}

func environmentOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeVariableGroups":      optionalBoolDefault(params, "includeVariableGroups", true),
		"includeLatestOrchestration": optionalBoolDefault(params, "includeLatestOrchestration", true),
	}
}

func applicationOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeEnvironments":   optionalBoolDefault(params, "includeEnvironments", true),
		"includeOrchestrations": optionalBoolDefault(params, "includeOrchestrations", true),
		"envLimit":              optionalIntDefault(params, "envLimit", 5),
		"orchestrationLimit":    optionalIntDefault(params, "orchestrationLimit", 5),
	}
}

func handleGetReleaseOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, systemName, sn, err := requiredSystemRelease(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	releasePath := appstackSystemReleasePath(organizationID, systemName, sn)

	release, err := c.GetJSON(ctx, releasePath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"release": release,
		"filters": releaseOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeMembers", true) {
		members, err := c.GetJSON(ctx, releasePath+"/members", nil)
		if err != nil {
			return "", err
		}
		overview["members"] = members
	}

	if optionalBoolDefault(params, "includeProducts", true) {
		products, err := c.GetJSON(ctx, releasePath+"/products", nil)
		if err != nil {
			return "", err
		}
		overview["products"] = products
	}

	if optionalBoolDefault(params, "includeChangeRequests", true) {
		crQuery := url.Values{}
		crQuery.Set("current", "1")
		crQuery.Set("pageSize", strconv.Itoa(optionalIntDefault(params, "changeRequestLimit", 5)))
		crs, err := c.GetJSON(ctx, releasePath+"/changeRequests", crQuery)
		if err != nil {
			return "", err
		}
		overview["changeRequests"] = crs
	}

	return marshalPretty(overview)
}

func releaseOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeMembers":        optionalBoolDefault(params, "includeMembers", true),
		"includeProducts":       optionalBoolDefault(params, "includeProducts", true),
		"includeChangeRequests": optionalBoolDefault(params, "includeChangeRequests", true),
		"changeRequestLimit":    optionalIntDefault(params, "changeRequestLimit", 5),
	}
}
