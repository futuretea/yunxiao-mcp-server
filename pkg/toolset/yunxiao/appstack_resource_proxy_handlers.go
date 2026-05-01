package yunxiao

import (
	"context"
	"net/url"
)

func handleGetPodContainerLog(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, resourcePath, namespace, name, err := requiredResourceProxyObject(params)
	if err != nil {
		return "", err
	}
	container, err := requiredString(params, "container")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("tailingLines", "1000")
	setOptionalInt(query, params, "tailingLines")

	path := appstackResourceProxyPath(organizationID, resourcePath, namespace) + "/pods/" + url.PathEscape(name) + "/containers/" + url.PathEscape(container) + ":logs"
	return c.GetJSON(ctx, path, query)
}

func handleGetPodInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, resourcePath, namespace, name, err := requiredResourceProxyObject(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "taskSn")

	path := appstackResourceProxyPath(organizationID, resourcePath, namespace) + "/pods/" + url.PathEscape(name) + "/info"
	return c.GetJSON(ctx, path, query)
}

func handleGetKubernetesObjectInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, resourcePath, namespace, name, err := requiredResourceProxyObject(params)
	if err != nil {
		return "", err
	}
	kind, err := requiredString(params, "kind")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "taskSn")

	path := appstackResourceProxyPath(organizationID, resourcePath, namespace) + "/" + url.PathEscape(kind) + "/" + url.PathEscape(name) + "/info"
	return c.GetJSON(ctx, path, query)
}

func handleGetDeploymentRevisionInfo(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, appName, envName, namespace, name, revision, err := requiredDeploymentRevision(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "taskSn")

	path := appstackAppPath(organizationID, appName) + "/envs/" + url.PathEscape(envName) + "/ns/" + url.PathEscape(namespace) + "/deployments/" + url.PathEscape(name) + "/revisions/" + url.PathEscape(revision)
	return c.GetJSON(ctx, path, query)
}

func requiredResourceProxyObject(params map[string]any) (string, string, string, string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", "", "", "", err
	}
	resourcePath, err := requiredString(params, "resourcePath")
	if err != nil {
		return "", "", "", "", err
	}
	namespace, err := requiredString(params, "namespace")
	if err != nil {
		return "", "", "", "", err
	}
	name, err := requiredString(params, "name")
	if err != nil {
		return "", "", "", "", err
	}
	return organizationID, resourcePath, namespace, name, nil
}

func requiredDeploymentRevision(params map[string]any) (string, string, string, string, string, string, error) {
	organizationID, appName, envName, err := requiredAppEnvironment(params)
	if err != nil {
		return "", "", "", "", "", "", err
	}
	namespace, err := requiredString(params, "namespace")
	if err != nil {
		return "", "", "", "", "", "", err
	}
	name, err := requiredString(params, "name")
	if err != nil {
		return "", "", "", "", "", "", err
	}
	revision, err := requiredString(params, "revision")
	if err != nil {
		return "", "", "", "", "", "", err
	}
	return organizationID, appName, envName, namespace, name, revision, nil
}

func appstackResourceProxyPath(organizationID, resourcePath, namespace string) string {
	return "/appstack/organizations/" + url.PathEscape(organizationID) + "/" + url.PathEscape(resourcePath) + "/" + url.PathEscape(namespace)
}
