package yunxiao

import (
	"errors"
	"fmt"
)

var errNoCategories = errors.New("categories must include at least one category")

func getClient(client any) (*Client, error) {
	c, ok := client.(*Client)
	if !ok || c == nil {
		return nil, fmt.Errorf("yunxiao client is not configured")
	}
	return c, nil
}

func requiredString(params map[string]any, key string) (string, error) {
	value, _ := params[key].(string)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func requiredOrganizationAndRepository(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	repositoryID, err := requiredString(params, "repositoryId")
	if err != nil {
		return "", "", err
	}
	return organizationID, repositoryID, nil
}

func requiredOrganizationAndNamedID(params map[string]any, key string) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	id, err := requiredString(params, key)
	if err != nil {
		return "", "", err
	}
	return organizationID, id, nil
}

func requiredOrganizationRepositoryAndLocalID(params map[string]any) (string, string, string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", "", "", err
	}
	localID, err := requiredString(params, "localId")
	if err != nil {
		return "", "", "", err
	}
	return organizationID, repositoryID, localID, nil
}

func requiredOrganizationAndPipeline(params map[string]any) (string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", err
	}
	pipelineID, err := requiredString(params, "pipelineId")
	if err != nil {
		return "", "", err
	}
	return organizationID, pipelineID, nil
}
