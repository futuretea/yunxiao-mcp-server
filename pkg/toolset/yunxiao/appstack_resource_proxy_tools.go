package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackResourceProxyTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_pod_container_log",
				mcp.WithDescription("Get container logs from a pod in an AppStack resource proxy. Use this to retrieve recent logs for debugging deployment issues."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource proxy path. Typically discovered from application environment resource configurations.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace. Use get_pod_info or get_kubernetes_object_info to discover valid namespaces.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Pod name. Use get_kubernetes_object_info to discover valid pod names.")),
				mcp.WithString("container", mcp.Required(), mcp.Description("Container name within the pod. Typically discovered from pod info.")),
				mcp.WithNumber("tailingLines", mcp.Description("Number of recent log lines to return. Defaults to 1000.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPodContainerLog,
		},
		{
			Tool: mcp.NewTool("get_pod_info",
				mcp.WithDescription("Get detailed information about a pod in an AppStack resource proxy. Use this to inspect pod status, containers, and events."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource proxy path. Typically discovered from application environment resource configurations.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace. Use get_kubernetes_object_info to discover valid namespaces.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Pod name. Use get_kubernetes_object_info to discover valid pod names.")),
				mcp.WithString("taskSn", mcp.Description("Optional deployment task serial number for filtering pod info by task.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPodInfo,
		},
		{
			Tool: mcp.NewTool("get_kubernetes_object_info",
				mcp.WithDescription("Get detailed information about a Kubernetes object in an AppStack resource proxy. Supports pods, deployments, services, and other Kubernetes resources."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource proxy path. Typically discovered from application environment resource configurations.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace. Use get_application_overview or get_environment_overview to discover valid namespaces.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Kubernetes object name. Use Kubernetes conventions to identify objects by name.")),
				mcp.WithString("kind", mcp.Required(), mcp.Description("Kubernetes object kind. Valid values: Pod, Deployment, Service, ConfigMap, Secret, Ingress, StatefulSet, DaemonSet, etc.")),
				mcp.WithString("taskSn", mcp.Description("Optional deployment task serial number for filtering object info by task.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetKubernetesObjectInfo,
		},
		{
			Tool: mcp.NewTool("get_deployment_revision_info",
				mcp.WithDescription("Get revision information for an AppStack deployment. Use this to inspect the rollout history and revision details of a deployment."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name. Use list_applications to discover valid app names.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name. Use list_environments to discover valid environment names.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace. Use get_application_overview or get_environment_overview to discover valid namespaces.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Deployment name. Use get_kubernetes_object_info with kind=Deployment to discover valid deployment names.")),
				mcp.WithString("revision", mcp.Required(), mcp.Description("Deployment revision number. Use Kubernetes rollout history to discover valid revision numbers.")),
				mcp.WithString("taskSn", mcp.Description("Optional deployment task serial number for filtering revision info by task.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeploymentRevisionInfo,
		},
	}
}
