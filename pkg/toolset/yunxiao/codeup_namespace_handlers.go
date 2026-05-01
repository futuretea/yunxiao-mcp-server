package yunxiao

import (
	"context"
	"net/url"
)

func handleListTemplateRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	templateType, err := requiredNumberPathString(params, "templateType")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	query.Set("templateType", templateType)

	return c.GetJSONWithMetadata(ctx, codeupOrganizationPath(organizationID)+"/repositories/templates", query)
}

func handleListNamespaces(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "parentId")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "search")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")

	return c.GetJSONWithMetadata(ctx, codeupOrganizationPath(organizationID)+"/namespaces", query)
}

func handleGetNamespace(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, namespaceID, err := requiredOrganizationAndNamedID(params, "namespaceId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupOrganizationPath(organizationID) + "/namespaces/" + encodePathValue(namespaceID)
	return c.GetJSON(ctx, path, nil)
}

func handleGetOrgNamespace(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, codeupOrganizationPath(organizationID)+"/orgNamespace", nil)
}

func codeupOrganizationPath(organizationID string) string {
	return "/codeup/organizations/" + url.PathEscape(organizationID)
}
