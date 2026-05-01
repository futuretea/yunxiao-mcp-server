package yunxiao

import "context"

func handleGetUser(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	idOrUsername, err := requiredString(params, "idOrUsername")
	if err != nil {
		return "", err
	}

	return c.GetJSON(ctx, "/platform/users/"+encodePathValue(idOrUsername), nil)
}

func handleListAppExtensionFeatures(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, appExtensionType, err := requiredOrganizationAndNamedID(params, "type")
	if err != nil {
		return "", err
	}

	path := organizationPath(organizationID) + "/appExtensions/" + encodePathValue(appExtensionType) + "/features"
	return c.GetJSON(ctx, path, nil)
}
