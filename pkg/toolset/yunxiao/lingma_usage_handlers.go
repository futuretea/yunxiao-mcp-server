package yunxiao

import (
	"context"
	"fmt"
	"net/url"
)

func handleGetDepartmentUsage(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	query, err := requiredLingmaUsageQuery(params, "departmentId")
	if err != nil {
		return "", err
	}

	return c.GetJSONWithMetadata(ctx, lingmaOrganizationPath(organizationID)+"/departmentUsage", query)
}

func handleListDeveloperMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	query := url.Values{}
	setOptionalString(query, params, "departmentId")
	setOptionalString(query, params, "userId")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, lingmaOrganizationPath(organizationID)+"/developer/members", query)
}

func handleGetDeveloperUsage(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	query, err := requiredLingmaUsageQuery(params)
	if err != nil {
		return "", err
	}
	userID, _ := params["userId"].(string)
	departmentID, _ := params["departmentId"].(string)
	if userID == "" && departmentID == "" {
		return "", fmt.Errorf("userId or departmentId is required")
	}
	setOptionalString(query, params, "userId")
	setOptionalString(query, params, "departmentId")

	return c.GetJSONWithMetadata(ctx, lingmaOrganizationPath(organizationID)+"/developerUsage", query)
}

func requiredLingmaUsageQuery(params map[string]any, requiredKeys ...string) (url.Values, error) {
	startTime, err := requiredString(params, "startTime")
	if err != nil {
		return nil, err
	}
	endTime, err := requiredString(params, "endTime")
	if err != nil {
		return nil, err
	}

	query := url.Values{"startTime": []string{startTime}, "endTime": []string{endTime}}
	for _, key := range requiredKeys {
		value, err := requiredString(params, key)
		if err != nil {
			return nil, err
		}
		query.Set(key, value)
	}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	return query, nil
}

func lingmaOrganizationPath(organizationID string) string {
	return "/lingma/organizations/" + url.PathEscape(organizationID)
}
