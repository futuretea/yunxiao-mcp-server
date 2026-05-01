package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func appstackResourceProxyTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("get_pod_container_log",
				mcp.WithDescription("Get logs from a pod container through AppStack resource proxy."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource path segment.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Pod name.")),
				mcp.WithString("container", mcp.Required(), mcp.Description("Container name.")),
				mcp.WithNumber("tailingLines", mcp.Description("Number of log tail lines. Defaults to 1000.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPodContainerLog,
		},
		{
			Tool: mcp.NewTool("get_pod_info",
				mcp.WithDescription("Get pod information through AppStack resource proxy."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource path segment.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Pod name.")),
				mcp.WithString("taskSn", mcp.Description("Optional task serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetPodInfo,
		},
		{
			Tool: mcp.NewTool("get_kubernetes_object_info",
				mcp.WithDescription("Get Kubernetes object information through AppStack resource proxy."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("resourcePath", mcp.Required(), mcp.Description("Resource path segment.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace.")),
				mcp.WithString("kind", mcp.Required(), mcp.Description("Kubernetes object kind.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Kubernetes object name.")),
				mcp.WithString("taskSn", mcp.Description("Optional task serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetKubernetesObjectInfo,
		},
		{
			Tool: mcp.NewTool("get_deployment_revision_info",
				mcp.WithDescription("Get AppStack deployment workload revision information."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("appName", mcp.Required(), mcp.Description("Application name.")),
				mcp.WithString("envName", mcp.Required(), mcp.Description("Environment name.")),
				mcp.WithString("namespace", mcp.Required(), mcp.Description("Kubernetes namespace.")),
				mcp.WithString("name", mcp.Required(), mcp.Description("Deployment name.")),
				mcp.WithString("revision", mcp.Required(), mcp.Description("Deployment revision.")),
				mcp.WithString("taskSn", mcp.Description("Optional task serial number.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetDeploymentRevisionInfo,
		},
	}
}
