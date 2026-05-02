package yunxiao

import "github.com/futuretea/yunxiao-mcp-server/pkg/toolset"

// projectFocusedDomains are the tool domains enabled in project-focused mode.
var projectFocusedDomains = map[string]struct{}{
	"platform": {},
	"projex":   {},
}

// projectFocusedHiddenTools are raw projex tools superseded by enhanced alternatives.
var projectFocusedHiddenTools = map[string]struct{}{
	"get_project":             {},
	"get_sprint":              {},
	"list_work_item_comments": {},
}

// Toolset exposes Yunxiao OpenAPI tools.
type Toolset struct {
	ReadOnly bool
}

func (t *Toolset) GetName() string {
	return "yunxiao"
}

func (t *Toolset) GetDescription() string {
	return "Yunxiao organization and DevOps OpenAPI tools"
}

func (t *Toolset) GetTools(_ any) []toolset.ServerTool {
	tools := make([]toolset.ServerTool, 0, 174)
	tools = append(tools, withDomain(platformTools(), "platform")...)
	tools = append(tools, withDomain(codeupTools(), "codeup")...)
	tools = append(tools, withDomain(flowTools(), "flow")...)
	tools = append(tools, withDomain(projexTools(), "projex")...)
	tools = append(tools, withDomain(packageTools(), "packages")...)
	tools = append(tools, withDomain(appstackTools(), "appstack")...)
	tools = append(tools, withDomain(lingmaTools(), "lingma")...)
	return tools
}

// GetProjectFocusedTools returns only platform + projex tools, hiding
// low-value raw tools that have enhanced alternatives.
func (t *Toolset) GetProjectFocusedTools(client any) []toolset.ServerTool {
	all := t.GetTools(client)
	filtered := make([]toolset.ServerTool, 0, len(all))
	for _, tool := range all {
		if _, ok := projectFocusedDomains[tool.Domain]; !ok {
			continue
		}
		if _, ok := projectFocusedHiddenTools[tool.Tool.Name]; ok {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

func withDomain(tools []toolset.ServerTool, domain string) []toolset.ServerTool {
	for i := range tools {
		tools[i].Domain = domain
	}
	return tools
}
