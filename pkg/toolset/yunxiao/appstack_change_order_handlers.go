package yunxiao

import (
	"context"
	"net/url"
)

func handleListChangeOrderVersions(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "envNames")
	setOptionalString(query, params, "creators")
	query.Set("current", "1")
	query.Set("pageSize", "10")
	setOptionalInt(query, params, "current")
	setOptionalInt(query, params, "pageSize")

	path := appstackAppPath(organizationID, appName) + "/changeOrders/versions"
	return c.GetJSON(ctx, path, query)
}

func handleGetChangeOrder(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, changeOrderSn, err := requiredAppChangeOrder(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackChangeOrderPath(organizationID, appName, changeOrderSn)
	return c.GetJSON(ctx, path, nil)
}

func handleListChangeOrderJobLogs(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, changeOrderSn, jobSn, err := requiredAppChangeOrderJob(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("current", "1")
	query.Set("pageSize", "10")
	setOptionalInt(query, params, "current")
	setOptionalInt(query, params, "pageSize")

	path := appstackChangeOrderPath(organizationID, appName, changeOrderSn) + "/jobs/" + url.PathEscape(jobSn) + "/logs"
	return c.GetJSON(ctx, path, query)
}

func handleFindTaskOperationLog(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, changeOrderSn, jobSn, err := requiredAppChangeOrderJob(params)
	if err != nil {
		return "", err
	}
	stageSn, err := requiredString(params, "stageSn")
	if err != nil {
		return "", err
	}
	taskSn, err := requiredString(params, "taskSn")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackChangeOrderPath(organizationID, appName, changeOrderSn) + "/jobs/" + url.PathEscape(jobSn) + "/stages/" + url.PathEscape(stageSn) + "/tasks/" + url.PathEscape(taskSn) + "/operationLog"
	return c.GetJSON(ctx, path, nil)
}

func handleListChangeOrdersByOrigin(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	originType, err := requiredString(params, "originType")
	if err != nil {
		return "", err
	}
	originID, err := requiredString(params, "originId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("originType", originType)
	query.Set("originId", originID)
	setOptionalString(query, params, "appName")
	setOptionalString(query, params, "envName")

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/changeOrders:byOrigin"
	return c.GetJSON(ctx, path, query)
}

func requiredOrganizationAndApp(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	appName, err := requiredString(params, "appName")
	if err != nil {
		return "", "", err
	}
	return organizationID, appName, nil
}

func requiredAppChangeOrder(params map[string]any) (string, string, string, error) {
	organizationID, appName, err := requiredOrganizationAndApp(params)
	if err != nil {
		return "", "", "", err
	}
	changeOrderSn, err := requiredString(params, "changeOrderSn")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, appName, changeOrderSn, nil
}

func requiredAppChangeOrderJob(params map[string]any) (string, string, string, string, error) {
	organizationID, appName, changeOrderSn, err := requiredAppChangeOrder(params)
	if err != nil {
		return "", "", "", "", err
	}
	jobSn, err := requiredString(params, "jobSn")
	if err != nil {
		return "", "", "", "", err
	}
	return organizationID, appName, changeOrderSn, jobSn, nil
}

func appstackAppPath(organizationID, appName string) string {
	return "/appstack/organizations/" + url.PathEscape(organizationID) + "/apps/" + url.PathEscape(appName)
}

func appstackChangeOrderPath(organizationID, appName, changeOrderSn string) string {
	return appstackAppPath(organizationID, appName) + "/changeOrders/" + url.PathEscape(changeOrderSn)
}
