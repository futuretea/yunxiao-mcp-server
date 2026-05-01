package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaKnowledgeBaseTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_knowledge_bases",
				mcp.WithDescription("List Tongyi Lingma knowledge bases."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("query", mcp.Description("Knowledge base name fuzzy query.")),
				mcp.WithString("sceneType", mcp.Description("Scene type, such as chat or completion.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithString("userId", mcp.Description("User ID permission filter.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKnowledgeBases,
		},
		{
			Tool: mcp.NewTool("list_kb_files",
				mcp.WithDescription("List Tongyi Lingma knowledge base files."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("kbId", mcp.Required(), mcp.Description("Knowledge base ID.")),
				mcp.WithString("query", mcp.Description("File name fuzzy query.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKbFiles,
		},
		{
			Tool: mcp.NewTool("list_kb_members",
				mcp.WithDescription("List Tongyi Lingma knowledge base members."),
				mcp.WithString("organizationId", mcp.Required(), mcp.Description("Yunxiao organization ID.")),
				mcp.WithString("kbId", mcp.Required(), mcp.Description("Knowledge base ID.")),
				mcp.WithString("query", mcp.Description("Member name fuzzy query.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithNumber("page", mcp.Description("Page number.")),
				mcp.WithNumber("perPage", mcp.Description("Page size. Default is 20.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKbMembers,
		},
	}
}
