package yunxiao

import (
	"context"
	"strings"
	"testing"
)

func TestHandleDescribeToolset(t *testing.T) {
	result, err := handleDescribeToolset(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("handleDescribeToolset() error = %v", err)
	}
	if result == "" {
		t.Fatal("handleDescribeToolset() returned empty result")
	}

	// Should contain the capability tool itself and tools from multiple domains.
	if !strings.Contains(result, "describe_toolset") {
		t.Error("result should mention describe_toolset")
	}
	if !strings.Contains(result, "projex") {
		t.Error("result should contain projex domain")
	}
	if !strings.Contains(result, "flow") {
		t.Error("result should contain flow domain")
	}
	if !strings.Contains(result, "platform") {
		t.Error("result should contain platform domain")
	}
}

func TestHandleDescribeToolsetWithDomainFilter(t *testing.T) {
	result, err := handleDescribeToolset(context.Background(), nil, map[string]any{"domain": "meta"})
	if err != nil {
		t.Fatalf("handleDescribeToolset() error = %v", err)
	}
	if !strings.Contains(result, "describe_toolset") {
		t.Error("result should mention describe_toolset")
	}
	if strings.Contains(result, "projex") {
		t.Error("result should not contain other domains when filtered to meta")
	}
}

func TestIsReadOnlyTool(t *testing.T) {
	ts := &Toolset{ReadOnly: false}
	allTools := ts.GetTools(nil)
	readOnlyCount := 0
	for _, tool := range allTools {
		if isReadOnlyTool(tool.Tool) {
			readOnlyCount++
		}
	}
	if readOnlyCount == 0 {
		t.Error("expected some read-only tools")
	}
	if readOnlyCount == len(allTools) {
		t.Error("expected some non-read-only tools")
	}
}
