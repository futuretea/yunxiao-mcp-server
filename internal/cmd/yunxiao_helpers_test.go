package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	"github.com/futuretea/yunxiao-mcp-server/pkg/toolset"
)

func TestYunxiaoCLIFlagBindersReportMissingFlags(t *testing.T) {
	if err := bindYunxiaoCLIFlags(viper.New(), &cobra.Command{}); err == nil {
		t.Fatal("bindYunxiaoCLIFlags() expected error for command without persistent flags")
	}
	if err := bindYunxiaoMCPFlags(viper.New(), &cobra.Command{}); err == nil {
		t.Fatal("bindYunxiaoMCPFlags() expected error for command without mcp flags")
	}
}

func TestNewSDKClientFromConfigAllowsInsecureTLS(t *testing.T) {
	client, err := newSDKClientFromConfig(&config.StaticConfig{
		BaseURL:               "https://example.com",
		AccessToken:           "token",
		RequestTimeoutSeconds: 1,
		InsecureSkipTLSVerify: true,
	})
	if err != nil {
		t.Fatalf("newSDKClientFromConfig() error = %v", err)
	}
	if client == nil {
		t.Fatal("newSDKClientFromConfig() returned nil client")
	}
}

func TestNewToolSummaryMarksWriteTools(t *testing.T) {
	summary := newToolSummary(toolset.ServerTool{
		Tool:   mcp.NewTool("create_workitem", mcp.WithDescription("create a work item")),
		Domain: "projex",
	})
	if summary.Access != "write" {
		t.Fatalf("Access = %q, want write", summary.Access)
	}
}

func TestParseToolParamsAcceptsEmptyInput(t *testing.T) {
	params, err := parseToolParamsWithInput("   ", "", nil)
	if err != nil {
		t.Fatalf("parseToolParamsWithInput() error = %v", err)
	}
	if len(params) != 0 {
		t.Fatalf("params = %#v, want empty map", params)
	}
}

func TestParseToolParamsReadsStdin(t *testing.T) {
	params, err := parseToolParamsWithInput("{}", "-", strings.NewReader(`{"page":4}`))
	if err != nil {
		t.Fatalf("parseToolParamsWithInput() error = %v", err)
	}
	if got := params["page"]; got != float64(4) {
		t.Fatalf("params[page] = %#v, want 4", got)
	}
}

func TestParseToolParamsStdinRequiresInput(t *testing.T) {
	if _, err := parseToolParamsWithInput("{}", "-", nil); err == nil {
		t.Fatal("parseToolParamsWithInput() expected stdin input error")
	}
}

func TestTaskListOptionsParamsIncludesFiltersAndPagination(t *testing.T) {
	params, err := (taskListOptions{
		OrganizationID: " org-1 ",
		ProjectID:      " project-1 ",
		Category:       " Bug ",
		Subject:        " crash ",
		Status:         " open,closed ",
		AssignedTo:     " user-1 ",
		Creator:        " user-2 ",
		Sprint:         " sprint-1 ",
		OrderBy:        " updatedAt ",
		Sort:           " desc ",
		Page:           2,
		PerPage:        50,
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"projectId":      "project-1",
		"category":       "Bug",
		"subject":        "crash",
		"status":         "open,closed",
		"assignedTo":     "user-1",
		"creator":        "user-2",
		"sprint":         "sprint-1",
		"orderBy":        "updatedAt",
		"sort":           "desc",
		"page":           2,
		"perPage":        50,
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPrintTaskListShowsNoResultsWhenRowsEmpty(t *testing.T) {
	var out bytes.Buffer
	raw := "No results found."
	if err := printTaskList(&out, raw); err != nil {
		t.Fatalf("printTaskList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want \"No results found.\"", out.String())
	}
}

func TestTaskRowsFromJSONExtractsNestedAndTypedValues(t *testing.T) {
	rows := taskRowsFromJSON(`{
		"result": {
			"items": [
				{
					"workItemId": 123,
					"title": "Fix crash",
					"status": {"name": "Open"},
					"assignee": {"displayName": "Alice"}
				},
				"skip non-object rows"
			]
		}
	}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := taskRow{ID: "123", Subject: "Fix crash", Status: "Open", Assignee: "Alice"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestTaskRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := taskRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}

func TestRowsFromJSONWithPresenceDistinguishesEmptyLists(t *testing.T) {
	rows, ok := rowsFromJSONWithPresence(`{"data":[]}`)
	if !ok {
		t.Fatal("rowsFromJSONWithPresence() ok = false, want true")
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}

	if _, ok := rowsFromJSONWithPresence(`{"data":{"total":0}}`); ok {
		t.Fatal("rowsFromJSONWithPresence() ok = true, want false for non-list payload")
	}
}

func TestStringifyCLIValueFormatsBools(t *testing.T) {
	if got := stringifyCLIValue(true); got != "true" {
		t.Fatalf("stringifyCLIValue(true) = %q, want true", got)
	}
	if got := stringifyCLIValue(false); got != "false" {
		t.Fatalf("stringifyCLIValue(false) = %q, want false", got)
	}
}

func TestStringifyCLIValueReturnsEmptyForUnsupportedTypes(t *testing.T) {
	if got := stringifyCLIValue([]string{"unsupported"}); got != "" {
		t.Fatalf("stringifyCLIValue(slice) = %q, want empty string", got)
	}
}
