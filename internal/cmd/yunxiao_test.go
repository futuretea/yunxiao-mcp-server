package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestYunxiaoCLIVersionPrintsCLIName(t *testing.T) {
	var out, errOut bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &errOut})
	command.SetArgs([]string{"version"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "yunxiao version=") {
		t.Fatalf("stdout = %q", out.String())
	}
	if strings.Contains(out.String(), "yunxiao-mcp-server") {
		t.Fatalf("stdout = %q, should use CLI binary name", out.String())
	}
	if errOut.Len() != 0 {
		t.Fatalf("stderr = %q", errOut.String())
	}
}

func TestYunxiaoCLIMCPSubcommandValidatesBeforeServing(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"mcp",
		"--port", "1",
		"--enabled-tools", "not_a_tool",
		"--access-token", "token",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected unknown tool error")
	}
	if !strings.Contains(err.Error(), `unknown MCP tool "not_a_tool"`) {
		t.Fatalf("error = %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLIMCPSubcommandStartsAndShutsDownHTTPServer(t *testing.T) {
	restoreLogger := preserveLogger()
	t.Cleanup(restoreLogger)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"mcp",
		"--port", strconv.Itoa(port),
		"--access-token", "token",
	})
	command.SetContext(ctx)

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestYunxiaoCLIToolsListJSONIsOfflineAndFiltered(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "get_current_user", "tools", "list", "--json"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var summaries []toolSummary
	if err := json.Unmarshal(out.Bytes(), &summaries); err != nil {
		t.Fatalf("unmarshal stdout: %v\n%s", err, out.String())
	}
	if len(summaries) != 1 {
		t.Fatalf("summary count = %d, want 1: %#v", len(summaries), summaries)
	}
	if summaries[0].Name != "get_current_user" || summaries[0].Domain != "platform" || summaries[0].Access != "read-only" {
		t.Fatalf("summary = %#v", summaries[0])
	}
}

func TestYunxiaoCLIToolsListTable(t *testing.T) {
	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--enabled-tools", "get_current_user", "tools", "list"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"NAME", "DOMAIN", "get_current_user", "platform"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
}

func TestYunxiaoCLITaskListPrintsTable(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/projex/organizations/org-1/workitems:search":
			body, _ := io.ReadAll(r.Body)
			gotBody = string(body)
			_, _ = w.Write([]byte(`{"data":[{"id":"wi-1","subject":"Fix bug","status":"Open","assignedTo":{"name":"Alice"}}]}`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"task", "list",
		"--project-id", "project-1",
		"--subject", "Fix",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "SUBJECT", "wi-1", "Fix bug", "Open", "Alice"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	for _, want := range []string{`"category":"Task"`, `"spaceId":"project-1"`} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("body = %q, missing %q", gotBody, want)
		}
	}
}

func TestYunxiaoCLITaskListRequiresProjectID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"task", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected project-id error")
	}
	if !strings.Contains(err.Error(), "project-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsCallInvokesSDKClient(t *testing.T) {
	var gotToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/platform/users:me":
			gotToken = r.Header.Get("x-yunxiao-token")
			_, _ = w.Write([]byte(`{"id":"u-1"}`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"--enabled-tools", "get_current_user",
		"tools", "call", "get_current_user",
		"--params", "{}",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotToken != "token-1" {
		t.Fatalf("token = %q, want token-1", gotToken)
	}
	if !strings.Contains(out.String(), `"id": "u-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIToolsCallRejectsInvalidBaseURL(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", "://invalid-url",
		"tools", "call", "get_current_user",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected invalid base URL error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIToolsCallFillsDefaultOrganizationID(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/platform/organizations/org-1/members":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[{"name":"Alice"}]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"--enabled-tools", "list_organization_members",
		"tools", "call", "list_organization_members",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/platform/organizations/org-1/members" {
		t.Fatalf("path = %q", gotPath)
	}
	if !strings.Contains(out.String(), `"data"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIToolsCallSkipsDefaultOrgWhenOrganizationProvided(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations/org-2/members":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[]`))
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"--enabled-tools", "list_organization_members",
		"tools", "call", "list_organization_members",
		"--params", `{"organizationId":"org-2"}`,
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "/oapi/v1/platform/organizations/org-2/members" {
		t.Fatalf("path = %q", gotPath)
	}
}

func TestYunxiaoCLIToolsCallValidatesToolBeforeNetwork(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--enabled-tools", "create_workitem",
		"tools", "call", "create_workitem",
		"--params", "{}",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected read-only write tool error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLIToolsCallValidatesRequiredParamsBeforeNetwork(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"--enabled-tools", "get_project_overview",
		"tools", "call", "get_project_overview",
		"--params", "{}",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected required params error")
	}
	if !strings.Contains(err.Error(), "[validation]") || !strings.Contains(err.Error(), "projectId is required") {
		t.Fatalf("error = %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
}

func TestYunxiaoCLIToolsCallReportsDefaultOrgResolveError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/platform/organizations" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		http.Error(w, `{"error":"denied"}`, http.StatusUnauthorized)
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "bad-token",
		"--enabled-tools", "list_organization_members",
		"tools", "call", "list_organization_members",
		"--params", "{}",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected auth error")
	}
	if !strings.Contains(err.Error(), "[auth]") || !strings.Contains(err.Error(), "Authentication failed") {
		t.Fatalf("error = %v", err)
	}
	if strings.Contains(err.Error(), "organizationId is required") {
		t.Fatalf("error should preserve auth failure, got %v", err)
	}
}

func TestYunxiaoCLIToolsCallReadsParamsFile(t *testing.T) {
	paramsPath := filepath.Join(t.TempDir(), "params.json")
	if err := os.WriteFile(paramsPath, []byte(`{"page":2}`), 0o600); err != nil {
		t.Fatalf("write params: %v", err)
	}

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/platform/organizations/org-1/members":
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--base-url", server.URL,
		"--access-token", "token-1",
		"--enabled-tools", "list_organization_members",
		"tools", "call", "list_organization_members",
		"--params-file", paramsPath,
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotQuery != "page=2" {
		t.Fatalf("query = %q, want page=2", gotQuery)
	}
}

func TestYunxiaoCLIToolsCallRejectsReadOnlyWriteTool(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--enabled-tools", "create_workitem",
		"tools", "call", "create_workitem",
		"--params", "{}",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected read-only write tool error")
	}
	if !strings.Contains(err.Error(), `unknown MCP tool "create_workitem"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestParseToolParamsRejectsNonObject(t *testing.T) {
	if _, err := parseToolParams("[]", ""); err == nil {
		t.Fatal("parseToolParams() expected error")
	}
}

func TestParseToolParamsRejectsTrailingJSON(t *testing.T) {
	if _, err := parseToolParams(`{} {}`, ""); err == nil {
		t.Fatal("parseToolParams() expected trailing JSON error")
	}
}

func TestParseToolParamsRejectsUnreadableFile(t *testing.T) {
	_, err := parseToolParams("{}", filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("parseToolParams() expected file read error")
	}
}
