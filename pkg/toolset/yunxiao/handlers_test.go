package yunxiao

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newHandlerTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL, "token", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func TestBuildProjectConditions(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		want   string
	}{
		{"empty", map[string]any{}, ""},
		{"name only", map[string]any{"name": "demo"}, `{"conditionGroups":[[{"className":"string","fieldIdentifier":"name","format":"input","operator":"CONTAINS","toValue":null,"value":["demo"]}]]}`},
		{"status only", map[string]any{"status": "TODO,DOING"}, `{"conditionGroups":[[{"className":"status","fieldIdentifier":"status","format":"list","operator":"CONTAINS","toValue":null,"value":["TODO","DOING"]}]]}`},
		{"creator only", map[string]any{"creator": "alice"}, `{"conditionGroups":[[{"className":"user","fieldIdentifier":"creator","format":"list","operator":"CONTAINS","toValue":null,"value":["alice"]}]]}`},
		{"multiple", map[string]any{"name": "demo", "status": "TODO"}, `{"conditionGroups":[[{"className":"string","fieldIdentifier":"name","format":"input","operator":"CONTAINS","toValue":null,"value":["demo"]},{"className":"status","fieldIdentifier":"status","format":"list","operator":"CONTAINS","toValue":null,"value":["TODO"]}]]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProjectConditions(tt.params)
			if got != tt.want {
				t.Fatalf("buildProjectConditions() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHandleListRepositoriesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/repositories" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" ||
			r.URL.Query().Get("perPage") != "10" ||
			r.URL.Query().Get("orderBy") != "name" ||
			r.URL.Query().Get("sort") != "asc" ||
			r.URL.Query().Get("search") != "demo" ||
			r.URL.Query().Get("archived") != "false" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"name":"repo"}]`))
	})

	result, err := handleListRepositories(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"page":           float64(2),
		"perPage":        float64(10),
		"orderBy":        "name",
		"sort":           "asc",
		"search":         "demo",
		"archived":       false,
	})
	if err != nil {
		t.Fatalf("handleListRepositories() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetRepositoryEncodesRepositoryPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2FDemo%20Repo") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"repo"}`))
	})

	if _, err := handleGetRepository(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/Demo Repo",
	}); err != nil {
		t.Fatalf("handleGetRepository() error = %v", err)
	}
}

func TestHandleListBranchesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/branches") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		if r.URL.Query().Get("sort") != "updated_desc" || r.URL.Query().Get("search") != "main" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"name":"main"}]`))
	})

	if _, err := handleListBranches(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"sort":           "updated_desc",
		"search":         "main",
	}); err != nil {
		t.Fatalf("handleListBranches() error = %v", err)
	}
}

func TestHandleGetBranchEncodesBranchPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/branches/feature%2Fdemo") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if strings.Contains(r.RequestURI, "%252F") {
			t.Fatalf("RequestURI = %q, contains double-encoded slash", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"name":"feature/demo"}`))
	})

	if _, err := handleGetBranch(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"branchName":     "feature/demo",
	}); err != nil {
		t.Fatalf("handleGetBranch() error = %v", err)
	}
}

func TestHandleListFilesBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/files/tree") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("path") != "cmd" ||
			r.URL.Query().Get("ref") != "main" ||
			r.URL.Query().Get("type") != "RECURSIVE" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"path":"cmd"}]`))
	})

	if _, err := handleListFiles(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"path":           "cmd",
		"ref":            "main",
		"type":           "RECURSIVE",
	}); err != nil {
		t.Fatalf("handleListFiles() error = %v", err)
	}
}

func TestHandleGetFileBlobsEncodesFilePath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/files/src%2Fmain.go") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("ref") != "main" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"content":"package main"}`))
	})

	if _, err := handleGetFileBlobs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"filePath":       "/src/main.go",
		"ref":            "main",
	}); err != nil {
		t.Fatalf("handleGetFileBlobs() error = %v", err)
	}
}

func TestHandleListCommitsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/commits") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("refName") != "main" ||
			r.URL.Query().Get("path") != "cmd/main.go" ||
			r.URL.Query().Get("showSignature") != "false" ||
			r.URL.Query().Get("committerIds") != "user-1,user-2" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"id":"sha-1"}]`))
	})

	result, err := handleListCommits(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"refName":        "main",
		"path":           "cmd/main.go",
		"showSignature":  false,
		"committerIds":   "user-1,user-2",
	})
	if err != nil {
		t.Fatalf("handleListCommits() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetCommitBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/commits/abc123") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"id":"abc123"}`))
	})

	if _, err := handleGetCommit(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"sha":            "abc123",
	}); err != nil {
		t.Fatalf("handleGetCommit() error = %v", err)
	}
}

func TestHandleCompareBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/compares") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("from") != "main" ||
			r.URL.Query().Get("to") != "release" ||
			r.URL.Query().Get("sourceType") != "branch" ||
			r.URL.Query().Get("targetType") != "branch" ||
			r.URL.Query().Get("straight") != "false" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"commits":[]}`))
	})

	if _, err := handleCompare(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"from":           "main",
		"to":             "release",
		"sourceType":     "branch",
		"targetType":     "branch",
		"straight":       "false",
	}); err != nil {
		t.Fatalf("handleCompare() error = %v", err)
	}
}
