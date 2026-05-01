package yunxiao

import (
	"context"
	"net/url"
)

func handleGetMachineDeployLog(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	tunnelID, err := requiredNumberPathString(params, "tunnelId")
	if err != nil {
		return "", err
	}
	machineSn, err := requiredString(params, "machineSn")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("tunnelId", tunnelID)
	query.Set("machineSn", machineSn)

	path := "/appstack/organizations/" + url.PathEscape(organizationID) + "/host/deployLog"
	return c.GetJSON(ctx, path, query)
}

func handleGetDeployGroup(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, poolName, err := requiredOrganizationAndPool(params)
	if err != nil {
		return "", err
	}
	deployGroupName, err := requiredString(params, "deployGroupName")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackPoolPath(organizationID, poolName) + "/deployGroups/" + url.PathEscape(deployGroupName)
	return c.GetJSON(ctx, path, nil)
}

func handleListResourceInstances(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, poolName, err := requiredOrganizationAndPool(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "pagination")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "nextToken")
	setOptionalInt(query, params, "page")

	path := appstackPoolPath(organizationID, poolName) + "/instances"
	return c.GetJSON(ctx, path, query)
}

func handleGetResourceInstance(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, poolName, err := requiredOrganizationAndPool(params)
	if err != nil {
		return "", err
	}
	instanceName, err := requiredString(params, "instanceName")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := appstackPoolPath(organizationID, poolName) + "/instances/" + url.PathEscape(instanceName)
	return c.GetJSON(ctx, path, nil)
}

func requiredOrganizationAndPool(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	poolName, err := requiredString(params, "poolName")
	if err != nil {
		return "", "", err
	}
	return organizationID, poolName, nil
}

func appstackPoolPath(organizationID, poolName string) string {
	return "/appstack/organizations/" + url.PathEscape(organizationID) + "/pools/" + url.PathEscape(poolName)
}
