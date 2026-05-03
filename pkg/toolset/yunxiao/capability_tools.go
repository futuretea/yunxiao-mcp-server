package yunxiao

import (
	"context"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func capabilityTools() []toolset.ServerTool {
	return []toolset.ServerTool{
		{
			Tool: mcp.NewTool("describe_toolset",
				mcp.WithDescription("Get a catalog of all available Yunxiao tools organized by domain. Use this when you are unsure which tool to use or want to discover what operations are supported. The returned catalog includes each tool's name, domain, description, and read-only status."),
				mcp.WithString("domain", mcp.Description("Optional domain filter. Valid values: platform, projex, codeup, flow, packages, appstack, lingma, api. Omit to list all domains.")),
				mcp.WithReadOnlyHintAnnotation(true),
			),
			Handler: handleDescribeToolset,
		},
	}
}

func handleDescribeToolset(ctx context.Context, client any, params map[string]any) (string, error) {
	domainFilter := optionalStringDefault(params, "domain", "")

	// Use a fresh Toolset with ReadOnly=false to always show the full catalog.
	ts := &Toolset{ReadOnly: false}
	allTools := ts.GetTools(nil)

	domainTools := make(map[string][]map[string]any)
	var domains []string

	for _, t := range allTools {
		if domainFilter != "" && t.Domain != domainFilter {
			continue
		}
		info := map[string]any{
			"name":        t.Tool.Name,
			"description": t.Tool.Description,
			"readOnly":    isReadOnlyTool(t.Tool),
		}
		if _, ok := domainTools[t.Domain]; !ok {
			domains = append(domains, t.Domain)
		}
		domainTools[t.Domain] = append(domainTools[t.Domain], info)
	}

	sort.Strings(domains)

	catalog := make([]map[string]any, 0, len(domains))
	for _, d := range domains {
		catalog = append(catalog, map[string]any{
			"domain":      d,
			"toolCount":   len(domainTools[d]),
			"tools":       domainTools[d],
		})
	}

	result := map[string]any{
		"totalTools": len(allTools),
		"catalog":    catalog,
	}
	if domainFilter != "" {
		result["domainFilter"] = domainFilter
	}

	return marshalPretty(result)
}

func isReadOnlyTool(tool mcp.Tool) bool {
	if tool.Annotations.ReadOnlyHint != nil {
		return *tool.Annotations.ReadOnlyHint
	}
	return false
}
