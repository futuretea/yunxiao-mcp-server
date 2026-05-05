package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Param struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

type Tool struct {
	Name           string
	Description    string
	Params         []Param
	PaginationMode string
	Enhanced       bool
}

func main() {
	baseDir := "pkg/toolset/yunxiao"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}

	domainTools := map[string][]Tool{}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), "_tools.go") {
			continue
		}

		domain := extractDomain(f.Name())
		if domain == "" {
			continue
		}

		path := filepath.Join(baseDir, f.Name())
		enhanced := strings.Contains(f.Name(), "_enhanced_")
		tools, err := extractTools(path, enhanced)
		if err != nil {
			fmt.Fprintf(os.Stderr, "extract %s: %v\n", path, err)
			continue
		}
		domainTools[domain] = append(domainTools[domain], tools...)
	}

	if err := os.MkdirAll("docs", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir docs: %v\n", err)
		os.Exit(1)
	}

	for domain, tools := range domainTools {
		if len(tools) == 0 {
			continue
		}
		filename := fmt.Sprintf("docs/%s-tools.md", domain)
		if err := writeDomainDoc(filename, domain, tools); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", filename, err)
		}
	}

	fmt.Println("Generated docs for", len(domainTools), "domains")
}

func extractDomain(filename string) string {
	name := strings.TrimSuffix(filename, "_tools.go")
	switch {
	case strings.HasPrefix(name, "platform"):
		return "platform"
	case strings.HasPrefix(name, "codeup"):
		return "codeup"
	case strings.HasPrefix(name, "flow"):
		return "flow"
	case strings.HasPrefix(name, "appstack"):
		return "appstack"
	case strings.HasPrefix(name, "projex"):
		return "projex"
	case strings.HasPrefix(name, "package"):
		return "packages"
	case strings.HasPrefix(name, "lingma"):
		return "lingma"
	}
	return ""
}

func extractTools(filename string, enhanced bool) ([]Tool, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var tools []Tool
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "mcp" || sel.Sel.Name != "NewTool" {
			return true
		}
		if len(call.Args) < 1 {
			return true
		}

		tool := Tool{Enhanced: enhanced}
		if lit, ok := call.Args[0].(*ast.BasicLit); ok {
			tool.Name = strings.Trim(lit.Value, `"`)
		}

		for _, arg := range call.Args[1:] {
			optCall, ok := arg.(*ast.CallExpr)
			if !ok {
				continue
			}
			optSel, ok := optCall.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			optIdent, ok := optSel.X.(*ast.Ident)
			if !ok || optIdent.Name != "mcp" {
				continue
			}

			switch optSel.Sel.Name {
			case "WithDescription":
				if len(optCall.Args) > 0 {
					tool.Description = extractStringLit(optCall.Args[0])
				}
			case "WithString", "WithNumber", "WithBoolean", "WithArray":
				param := extractParam(optCall, optSel.Sel.Name)
				if param != nil {
					tool.Params = append(tool.Params, *param)
				}
			}
		}

		tools = append(tools, tool)
		return true
	})

	for i := range tools {
		tools[i].PaginationMode = detectPaginationMode(tools[i].Params)
	}
	return tools, nil
}

func detectPaginationMode(params []Param) string {
	hasNextToken := false
	hasCurrent := false
	hasPageSize := false
	hasPage := false
	hasPerPage := false
	for _, p := range params {
		switch p.Name {
		case "nextToken":
			hasNextToken = true
		case "current":
			hasCurrent = true
		case "pageSize":
			hasPageSize = true
		case "page":
			hasPage = true
		case "perPage":
			hasPerPage = true
		}
	}
	if hasNextToken {
		return "Keyset (nextToken)"
	}
	if hasCurrent || hasPageSize {
		return "Offset (current/pageSize)"
	}
	if hasPage || hasPerPage {
		return "Offset (page/perPage)"
	}
	return ""
}

func extractParam(call *ast.CallExpr, funcName string) *Param {
	if len(call.Args) < 1 {
		return nil
	}
	name := extractStringLit(call.Args[0])
	if name == "" {
		return nil
	}

	param := &Param{Name: name}
	switch funcName {
	case "WithString":
		param.Type = "string"
	case "WithNumber":
		param.Type = "number"
	case "WithBoolean":
		param.Type = "boolean"
	case "WithArray":
		param.Type = "array"
	}

	for _, arg := range call.Args[1:] {
		optCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		optSel, ok := optCall.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		optIdent, ok := optSel.X.(*ast.Ident)
		if !ok || optIdent.Name != "mcp" {
			continue
		}

		switch optSel.Sel.Name {
		case "Required":
			param.Required = true
		case "Description":
			if len(optCall.Args) > 0 {
				param.Description = extractStringLit(optCall.Args[0])
			}
		}
	}

	return param
}

func extractStringLit(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return ""
	}
	return strings.Trim(lit.Value, `"`)
}

func writeDomainDoc(filename, domain string, tools []Tool) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s Tools\n\n", cases.Title(language.English).String(domain))
	fmt.Fprintf(&b, "This document describes the %d read-only MCP tools in the %s domain.\n\n", len(tools), domain)

	var enhancedTools []Tool
	for _, t := range tools {
		if t.Enhanced {
			enhancedTools = append(enhancedTools, t)
		}
	}
	if len(enhancedTools) > 0 {
		b.WriteString("## Enhanced Tools\n\n")
		b.WriteString("These tools combine multiple Yunxiao OpenAPI calls into single, user-centric operations. Prefer them when available.\n\n")
		b.WriteString("| Tool | Description |\n")
		b.WriteString("|------|-------------|\n")
		for _, t := range enhancedTools {
			fmt.Fprintf(&b, "| `%s` | %s |\n", t.Name, t.Description)
		}
		b.WriteString("\n")
	}

	paginationModes := map[string]struct{}{}
	for _, t := range tools {
		if t.PaginationMode != "" {
			paginationModes[t.PaginationMode] = struct{}{}
		}
	}
	if len(paginationModes) > 0 {
		modes := make([]string, 0, len(paginationModes))
		for m := range paginationModes {
			modes = append(modes, m)
		}
		b.WriteString("## Pagination\n\n")
		b.WriteString("Tools in this domain use the following pagination scheme(s):\n\n")
		for _, m := range modes {
			fmt.Fprintf(&b, "- %s\n", m)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Tool Inventory\n\n")
	if len(enhancedTools) > 0 {
		b.WriteString("Tools marked in **bold** are enhanced aggregation tools.\n\n")
	}
	b.WriteString("| Tool | Description |\n")
	b.WriteString("|------|-------------|\n")
	for _, t := range tools {
		name := fmt.Sprintf("`%s`", t.Name)
		if t.Enhanced {
			name = fmt.Sprintf("**`%s`**", t.Name)
		}
		fmt.Fprintf(&b, "| %s | %s |\n", name, t.Description)
	}
	b.WriteString("\n")

	for _, t := range tools {
		fmt.Fprintf(&b, "### %s\n\n", t.Name)
		fmt.Fprintf(&b, "**Description**: %s\n\n", t.Description)
		if t.Enhanced {
			b.WriteString("**Type**: Enhanced aggregation tool\n\n")
		}
		if t.PaginationMode != "" {
			fmt.Fprintf(&b, "**Pagination**: %s\n\n", t.PaginationMode)
		}
		if len(t.Params) > 0 {
			b.WriteString("**Parameters**:\n\n")
			b.WriteString("| Name | Type | Required | Description |\n")
			b.WriteString("|------|------|----------|-------------|\n")
			for _, p := range t.Params {
				req := "No"
				if p.Required {
					req = "Yes"
				}
				desc := p.Description
				if desc == "" {
					desc = "-"
				}
				fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n", p.Name, p.Type, req, desc)
			}
			b.WriteString("\n")
		} else {
			b.WriteString("**Parameters**: None\n\n")
		}
	}

	return os.WriteFile(filename, []byte(b.String()), 0644)
}
