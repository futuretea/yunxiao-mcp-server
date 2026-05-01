package yunxiao

import (
	"context"
	"net/url"
)

func handleListSSHKeys(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := codeupKeyListQuery(params)
	return c.GetJSONWithMetadata(ctx, codeupOrganizationPath(organizationID)+"/keys", query)
}

func handleGetSSHKey(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	keyID, err := requiredNumberPathString(params, "keyId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupOrganizationPath(organizationID) + "/keys/" + url.PathEscape(keyID)
	return c.GetJSON(ctx, path, nil)
}

func handleListUserSSHKeys(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, userID, err := requiredOrganizationAndNamedID(params, "userId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupOrganizationPath(organizationID) + "/users/" + url.PathEscape(userID) + "/keys"
	return c.GetJSONWithMetadata(ctx, path, codeupKeyListQuery(params))
}

func handleListWebHooks(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
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

	return c.GetJSONWithMetadata(ctx, codeupRepositoryPath(organizationID, repositoryID)+"/webhooks", query)
}

func handleGetWebHook(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	hookID, err := requiredNumberPathString(params, "hookId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := codeupRepositoryPath(organizationID, repositoryID) + "/webhooks/" + url.PathEscape(hookID)
	return c.GetJSON(ctx, path, nil)
}

func codeupKeyListQuery(params map[string]any) url.Values {
	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	return query
}
