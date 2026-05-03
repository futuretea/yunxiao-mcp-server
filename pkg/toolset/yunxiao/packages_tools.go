package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func packageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_package_repositories",
				mcp.WithDescription("List artifact repositories (Packages) in a Yunxiao organization. Use this to discover repository IDs for listing artifacts."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repoTypes", mcp.Description("Comma-separated repository types: GENERIC, DOCKER, MAVEN, NPM, or NUGET.")),
				mcp.WithString("repoCategories", mcp.Description("Comma-separated repository modes: Hybrid, Local, Proxy, ProxyCache, or Group.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPackageRepositories,
		},
		{
			Tool: mcp.NewTool("list_artifacts",
				mcp.WithDescription("List artifacts in a Packages repository. Requires a repository ID from list_package_repositories."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("repoId", mcp.Required(), mcp.Description("Packages repository ID. Use list_package_repositories to discover valid IDs.")),
				mcp.WithString("repoType", mcp.Required(), mcp.Description("Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 100 when omitted.")),
				mcp.WithString("search", mcp.Description("Package name search text.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: latestUpdate or gmtDownload.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListArtifacts,
		},
	}
}
