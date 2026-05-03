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

	appPath := "/appstack/organizations/" + url.PathEscape(organizationID) + "/apps/" + url.PathEscape(appName)

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

func applicationOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeEnvironments":   optionalBoolDefault(params, "includeEnvironments", true),
		"includeOrchestrations": optionalBoolDefault(params, "includeOrchestrations", true),
		"envLimit":              optionalIntDefault(params, "envLimit", 5),
		"orchestrationLimit":    optionalIntDefault(params, "orchestrationLimit", 5),
	}
}
