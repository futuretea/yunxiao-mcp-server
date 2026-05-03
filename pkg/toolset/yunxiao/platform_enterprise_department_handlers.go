package yunxiao

import (
	"context"
	"net/url"
)

func handleListEnterpriseDepartments(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "parentId")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")

	return c.GetJSONWithMetadata(ctx, "/platform/departments", query)
}

func handleGetEnterpriseDepartment(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	departmentID, err := requiredString(params, "departmentId")
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, "/platform/departments/"+encodePathValue(departmentID), nil)
}
