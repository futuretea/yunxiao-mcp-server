package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestYunxiaoCLIPipelineResourceMemberListPrintsTableWithDefaultOrganization(t *testing.T) {
	var gotPath string
	requests := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path]++
		switch r.URL.Path {
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
		case "/oapi/v1/flow/organizations/org-1/resourceMembers/resourceTypes/pipeline/resourceIds/pipe-1":
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`[{"id":"user-1","name":"Alice","role":"Owner"}]`))
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
		"pipeline", "resource-member", "list",
		"--resource-id", "pipe-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"ID", "NAME", "ROLE", "user-1", "Alice", "Owner"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
	if gotPath != "/oapi/v1/flow/organizations/org-1/resourceMembers/resourceTypes/pipeline/resourceIds/pipe-1" {
		t.Fatalf("path = %q", gotPath)
	}
	if requests["/oapi/v1/platform/organizations"] != 1 || requests["/oapi/v1/flow/organizations/org-1/resourceMembers/resourceTypes/pipeline/resourceIds/pipe-1"] != 1 {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestYunxiaoCLIPipelineResourceMembersAliasListPrintsJSONWithExplicitOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-2/resourceMembers/resourceTypes/pipeline/resourceIds/pipe-1":
			_, _ = w.Write([]byte(`[{"userId":"u-1","displayName":"Bob","roleType":"Admin"}]`))
		case "/oapi/v1/platform/organizations":
			t.Fatal("should not resolve default organization when organizationId is provided")
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
		"pipelines", "resource-members", "list",
		"--organization-id", "org-2",
		"--resource-id", "pipe-1",
		"--json",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"userId": "u-1"`) {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestYunxiaoCLIPipelineResourceMemberListWithCustomType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oapi/v1/flow/organizations/org-1/resourceMembers/resourceTypes/hostGroup/resourceIds/hg-1":
			_, _ = w.Write([]byte(`[{"id":"u-1"}]`))
		case "/oapi/v1/platform/organizations":
			_, _ = w.Write([]byte(`[{"id":"org-1"}]`))
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
		"pipeline", "resource-member", "list",
		"--resource-type", "hostGroup",
		"--resource-id", "hg-1",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestYunxiaoCLIPipelineResourceMemberListRequiresResourceID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"pipeline", "resource-member", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected resource-id error")
	}
	if !strings.Contains(err.Error(), "resource-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLIPipelineResourceMemberListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "flow", "pipeline", "resource-member", "list", "--resource-id", "pipe-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_resource_members"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestPipelineResourceMemberListOptionsParamsDefaultsResourceTypeToPipeline(t *testing.T) {
	params, err := (pipelineResourceMemberListOptions{
		OrganizationID: " org-1 ",
		ResourceType:   " pipeline ",
		ResourceID:     " pipe-1 ",
	}).params()
	if err != nil {
		t.Fatalf("params() error = %v", err)
	}

	wants := map[string]any{
		"organizationId": "org-1",
		"resourceType":   "pipeline",
		"resourceId":     "pipe-1",
	}
	for key, want := range wants {
		if got := params[key]; got != want {
			t.Fatalf("params[%q] = %#v, want %#v", key, got, want)
		}
	}
}

func TestPipelineResourceMemberListOptionsParamsRequiresResourceID(t *testing.T) {
	if _, err := (pipelineResourceMemberListOptions{ResourceType: "pipeline"}).params(); err == nil {
		t.Fatal("params() expected resource-id error")
	}
}

func TestPipelineResourceMemberListOptionsParamsRequiresResourceType(t *testing.T) {
	if _, err := (pipelineResourceMemberListOptions{ResourceID: "pipe-1"}).params(); err == nil {
		t.Fatal("params() expected resource-type error")
	}
}

func TestPrintPipelineResourceMemberListFallsBackToRawJSON(t *testing.T) {
	var out bytes.Buffer
	raw := `{"data":{"total":0}}`
	if err := printPipelineResourceMemberList(&out, raw); err != nil {
		t.Fatalf("printPipelineResourceMemberList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != raw {
		t.Fatalf("stdout = %q, want raw JSON", out.String())
	}
}

func TestPrintPipelineResourceMemberListPrintsHeaderForEmptyList(t *testing.T) {
	var out bytes.Buffer
	if err := printPipelineResourceMemberList(&out, `[]`); err != nil {
		t.Fatalf("printPipelineResourceMemberList() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "ID  NAME  ROLE" {
		t.Fatalf("stdout = %q", out.String())
	}
}

func TestPipelineResourceMemberRowsFromJSONExtractsAlternateFields(t *testing.T) {
	rows := pipelineResourceMemberRowsFromJSON(`{"result":{"items":[{"accountId":"a-1","username":"alice","accessLevel":"admin"}]}}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineResourceMemberRow{ID: "a-1", Name: "alice", Role: "admin"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineResourceMemberRowsFromJSONSkipsNonObjectRows(t *testing.T) {
	rows := pipelineResourceMemberRowsFromJSON(`{"data":["skip",{"memberId":"m-1","nickName":"Bob","roleName":"Viewer"}]}`)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	want := pipelineResourceMemberRow{ID: "m-1", Name: "Bob", Role: "Viewer"}
	if rows[0] != want {
		t.Fatalf("row = %#v, want %#v", rows[0], want)
	}
}

func TestPipelineResourceMemberRowsFromJSONReturnsNilForInvalidPayload(t *testing.T) {
	if rows := pipelineResourceMemberRowsFromJSON(`not-json`); len(rows) != 0 {
		t.Fatalf("rows = %#v, want empty", rows)
	}
}
