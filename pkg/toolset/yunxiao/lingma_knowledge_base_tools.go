package yunxiao

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func lingmaKnowledgeBaseTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("list_knowledge_bases",
				mcp.WithDescription("List Tongyi Lingma knowledge bases in a Yunxiao organization. Knowledge bases contain curated documents for AI-assisted code completion and chat."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("query", mcp.Description("Knowledge base name fuzzy query.")),
				mcp.WithString("sceneType", mcp.Description("Scene type, such as chat or completion.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithString("userId", mcp.Description("User ID permission filter.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 20 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKnowledgeBases,
		},
		{
			Tool: mcp.NewTool("list_kb_files",
				mcp.WithDescription("List files in a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("kbId", mcp.Required(), mcp.Description("Knowledge base ID. Use list_knowledge_bases to discover valid IDs.")),
				mcp.WithString("query", mcp.Description("File name fuzzy query.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 20 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKbFiles,
		},
		{
			Tool: mcp.NewTool("list_kb_members",
				mcp.WithDescription("List members with access to a Tongyi Lingma knowledge base. Use list_knowledge_bases to discover valid kbId values."),
				mcp.WithString("organizationId", mcp.Description("Yunxiao organization ID. When omitted, the server uses the user's default organization.")),
				mcp.WithString("kbId", mcp.Required(), mcp.Description("Knowledge base ID. Use list_knowledge_bases to discover valid IDs.")),
				mcp.WithString("query", mcp.Description("Member name fuzzy query.")),
				mcp.WithString("orderBy", mcp.Description("Sort field.")),
				mcp.WithString("sort", mcp.Description("Sort order: desc or asc.")),
				mcp.WithNumber("page", mcp.Description("Page number for pagination. Starts at 1.")),
				mcp.WithNumber("perPage", mcp.Description("Page size for pagination. Supports 1-100. Defaults to 20 when omitted.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleListKbMembers,
		},
	}
}
