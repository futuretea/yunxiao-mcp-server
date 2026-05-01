package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func packageTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_package_repositories",
				mcp.WithDescription("List Packages repositories in a Yunxiao organization."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repoTypes", mcp.Description("Comma-separated repository types: GENERIC, DOCKER, MAVEN, NPM, or NUGET.")),
				mcp.WithString("repoCategories", mcp.Description("Comma-separated repository modes: Hybrid, Local, Proxy, ProxyCache, or Group.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListPackageRepositories,
		},
		{
			Tool: mcp.NewTool("list_artifacts",
				mcp.WithDescription("List artifacts in a Packages repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repoId", mcp.Required(), mcp.Description("Packages repository ID.")),
				mcp.WithString("repoType", mcp.Required(), mcp.Description("Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size.")),
				mcp.WithString("search", mcp.Description("Package name search text.")),
				mcp.WithString("orderBy", mcp.Description("Sort field: latestUpdate or gmtDownload.")),
				mcp.WithString("sort", mcp.Description("Sort direction: asc or desc.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListArtifacts,
		},
		{
			Tool: mcp.NewTool("get_artifact",
				mcp.WithDescription("Get one artifact from a Packages repository."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("repoId", mcp.Required(), mcp.Description("Packages repository ID.")),
				mcp.WithNumber("id", mcp.Required(), mcp.Description("Artifact ID.")),
				mcp.WithString("repoType", mcp.Required(), mcp.Description("Repository type: GENERIC, DOCKER, MAVEN, NPM, NUGET, or PYPI.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleGetArtifact,
		},
	}
}
