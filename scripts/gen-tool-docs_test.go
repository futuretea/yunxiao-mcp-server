package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractToolsReadOnlyAccessFromAnnotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture_tools.go")
	src := `package fixture

import "github.com/mark3labs/mcp-go/mcp"

func tools() {
	_ = mcp.NewTool("read_tool",
		mcp.WithDescription("Read something."),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	_ = mcp.NewTool("write_tool",
		mcp.WithDescription("Write something."),
	)
}
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	tools, err := extractTools(path, false)
	if err != nil {
		t.Fatalf("extractTools() error = %v", err)
	}
	if len(tools) != 2 {
		t.Fatalf("tool count = %d, want 2", len(tools))
	}

	byName := map[string]Tool{}
	for _, tool := range tools {
		byName[tool.Name] = tool
	}
	if !byName["read_tool"].ReadOnly {
		t.Fatal("read_tool should be read-only")
	}
	if byName["write_tool"].ReadOnly {
		t.Fatal("write_tool should be write-capable when annotation is absent")
	}
}

func TestWriteDomainDocShowsAccessTypes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "projex-tools.md")
	tools := []Tool{
		{Name: "search_workitems", Description: "Search work items.", ReadOnly: true},
		{Name: "create_workitem", Description: "Create a work item.", ReadOnly: false},
	}

	if err := writeDomainDoc(path, "projex", tools); err != nil {
		t.Fatalf("writeDomainDoc() error = %v", err)
	}
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated doc: %v", err)
	}
	content := string(contentBytes)

	mustContain := []string{
		"This document describes the 2 MCP tools in the projex domain.",
		"Access summary: 1 read-only, 1 write-capable.",
		"| Tool | Access | Description |",
		"| `search_workitems` | Read-only | Search work items. |",
		"| `create_workitem` | Write-capable | Create a work item. |",
		"**Access**: Read-only",
		"**Access**: Write-capable (requires `read_only=false`)",
	}
	for _, want := range mustContain {
		if !strings.Contains(content, want) {
			t.Fatalf("generated doc missing %q\n%s", want, content)
		}
	}
	if strings.Contains(content, "2 read-only MCP tools") {
		t.Fatalf("generated doc should not use domain-level all-read-only wording:\n%s", content)
	}
}

func TestGeneratedProjexDocMarksWriteTools(t *testing.T) {
	readTools, err := extractTools(filepath.Join("..", "pkg", "toolset", "yunxiao", "projex_tools.go"), false)
	if err != nil {
		t.Fatalf("extract read tools: %v", err)
	}
	writeTools, err := extractTools(filepath.Join("..", "pkg", "toolset", "yunxiao", "projex_write_tools.go"), false)
	if err != nil {
		t.Fatalf("extract write tools: %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "projex-tools.md")
	tools := append(readTools, writeTools...)
	if err := writeDomainDoc(path, "projex", tools); err != nil {
		t.Fatalf("writeDomainDoc() error = %v", err)
	}
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated doc: %v", err)
	}
	content := string(contentBytes)

	wantWriteTools := []string{
		"create_workitem",
		"update_workitem",
		"update_workitem_status",
		"add_workitem_comment",
	}
	for _, name := range wantWriteTools {
		want := "| `" + name + "` | Write-capable |"
		if !strings.Contains(content, want) {
			t.Fatalf("generated Projex doc missing write-capable access for %s", name)
		}
	}
	if !strings.Contains(content, "| `search_workitems` | Read-only |") {
		t.Fatal("generated Projex doc should keep read-only tools read-only")
	}
	if strings.Contains(content, "read-only MCP tools in the projex domain") {
		t.Fatal("generated Projex doc should not use all-read-only domain wording")
	}
}

func TestWriteDomainDocSortsPaginationModes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mixed-tools.md")
	tools := []Tool{
		{
			Name:           "page_tool",
			Description:    "Uses page pagination.",
			ReadOnly:       true,
			PaginationMode: "Offset (page/perPage)",
		},
		{
			Name:           "token_tool",
			Description:    "Uses token pagination.",
			ReadOnly:       true,
			PaginationMode: "Keyset (nextToken)",
		},
		{
			Name:           "current_tool",
			Description:    "Uses current pagination.",
			ReadOnly:       true,
			PaginationMode: "Offset (current/pageSize)",
		},
	}

	if err := writeDomainDoc(path, "mixed", tools); err != nil {
		t.Fatalf("writeDomainDoc() error = %v", err)
	}
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated doc: %v", err)
	}
	content := string(contentBytes)

	want := "- Keyset (nextToken)\n- Offset (current/pageSize)\n- Offset (page/perPage)"
	if !strings.Contains(content, want) {
		t.Fatalf("pagination modes not sorted, want sequence %q\n%s", want, content)
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"platform_tools.go", "platform"},
		{"platform_enhanced_tools.go", "platform"},
		{"codeup_write_tools.go", "codeup"},
		{"flow_tools.go", "flow"},
		{"flow_write_tools.go", "flow"},
		{"appstack_tools.go", "appstack"},
		{"projex_tools.go", "projex"},
		{"packages_tools.go", "packages"},
		{"lingma_tools.go", "lingma"},
		{"capability_tools.go", ""},
		{"unknown_domain_tools.go", ""},
		{"not_a_tools_file.go", ""},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := extractDomain(tt.filename); got != tt.want {
				t.Fatalf("extractDomain(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}
