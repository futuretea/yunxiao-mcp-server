package toolset

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// Toolset defines a group of MCP tools.
type Toolset interface {
	GetName() string
	GetDescription() string
	GetTools(client any) []ServerTool
}

// ServerTool combines an MCP tool definition with its handler.
type ServerTool struct {
	Tool    mcp.Tool
	Handler ToolHandler
	Domain  string
}

// ToolHandler handles a tool call.
type ToolHandler func(ctx context.Context, client any, params map[string]any) (string, error)
