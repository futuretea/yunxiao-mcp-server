package yunxiao

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var errNoCategories = errors.New("categories must include at least one category")

func getClient(client any) (*Client, error) {
	c, ok := client.(*Client)
	if !ok || c == nil {
		return nil, &ValidationError{Msg: "yunxiao client is not configured"}
	}
	return c, nil
}

func requiredString(params map[string]any, key string) (string, error) {
	value, _ := params[key].(string)
	if strings.TrimSpace(value) == "" {
		return "", &ValidationError{Msg: fmt.Sprintf("%s is required", key)}
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

func requiredNumberPathString(params map[string]any, key string) (string, error) {
	switch value := params[key].(type) {
	case float64:
		if value != math.Trunc(value) {
			return "", &ValidationError{Msg: fmt.Sprintf("%s must be an integer", key)}
		}
		return strconv.FormatInt(int64(value), 10), nil
	case int:
		return strconv.Itoa(value), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case string:
		if value = strings.TrimSpace(value); value != "" {
			return value, nil
		}
	}
	return "", &ValidationError{Msg: fmt.Sprintf("%s is required", key)}
}
