package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestSetOptionalIntBody(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		key    string
		want   any
		wantOk bool
	}{
		{"float64", map[string]any{"k": float64(42)}, "k", 42, true},
		{"int", map[string]any{"k": int(7)}, "k", 7, true},
		{"int64", map[string]any{"k": int64(99)}, "k", 99, true},
		{"string non-empty", map[string]any{"k": "123"}, "k", "123", true},
		{"string empty", map[string]any{"k": ""}, "k", nil, false},
		{"nil", map[string]any{}, "k", nil, false},
		{"bool", map[string]any{"k": true}, "k", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{}
			setOptionalIntBody(body, tt.params, tt.key)
			if tt.wantOk {
				if fmt.Sprint(body[tt.key]) != fmt.Sprint(tt.want) {
					t.Fatalf("body[%q] = %v, want %v", tt.key, body[tt.key], tt.want)
				}
			} else {
				if _, ok := body[tt.key]; ok {
					t.Fatalf("body[%q] should not be set", tt.key)
				}
			}
		})
	}
}

func TestSetOptionalIntAs(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		from   string
		to     string
		want   string
		wantOk bool
	}{
		{"float64", map[string]any{"k": float64(42)}, "k", "k", "42", true},
		{"int", map[string]any{"k": int(7)}, "k", "k", "7", true},
		{"int64", map[string]any{"k": int64(99)}, "k", "k", "99", true},
		{"string non-empty", map[string]any{"k": "123"}, "k", "k", "123", true},
		{"string empty", map[string]any{"k": ""}, "k", "k", "", false},
		{"nil", map[string]any{}, "k", "k", "", false},
		{"bool", map[string]any{"k": true}, "k", "k", "", false},
		{"rename", map[string]any{"from": "5"}, "from", "to", "5", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values{}
			setOptionalIntAs(query, tt.params, tt.from, tt.to)
			if tt.wantOk {
				if got := query.Get(tt.to); got != tt.want {
					t.Fatalf("query[%q] = %q, want %q", tt.to, got, tt.want)
				}
			} else {
				if query.Get(tt.to) != "" {
					t.Fatalf("query[%q] should not be set", tt.to)
				}
			}
		})
	}
}

func TestSetOptionalStringArrayBody(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		key    string
		want   []string
		wantOk bool
	}{
		{"[]any with values", map[string]any{"k": []any{"a", "b"}}, "k", []string{"a", "b"}, true},
		{"[]any with whitespace", map[string]any{"k": []any{" a ", "  b  "}}, "k", []string{"a", "b"}, true},
		{"[]any mixed types", map[string]any{"k": []any{"a", 1, "b"}}, "k", []string{"a", "b"}, true},
		{"[]any all non-string", map[string]any{"k": []any{1, 2}}, "k", nil, false},
		{"[]any empty", map[string]any{"k": []any{}}, "k", nil, false},
		{"[]string with values", map[string]any{"k": []string{"a", "b"}}, "k", []string{"a", "b"}, true},
		{"[]string with whitespace", map[string]any{"k": []string{" a ", "b"}}, "k", []string{"a", "b"}, true},
		{"[]string all empty", map[string]any{"k": []string{"", "  "}}, "k", nil, false},
		{"[]string empty", map[string]any{"k": []string{}}, "k", nil, false},
		{"nil", map[string]any{}, "k", nil, false},
		{"string not array", map[string]any{"k": "not-array"}, "k", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{}
			setOptionalStringArrayBody(body, tt.params, tt.key)
			if tt.wantOk {
				got, ok := body[tt.key].([]string)
				if !ok {
					t.Fatalf("body[%q] type = %T, want []string", tt.key, body[tt.key])
				}
				if len(got) != len(tt.want) {
					t.Fatalf("body[%q] = %v, want %v", tt.key, got, tt.want)
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Fatalf("body[%q][%d] = %q, want %q", tt.key, i, got[i], tt.want[i])
					}
				}
			} else {
				if _, ok := body[tt.key]; ok {
					t.Fatalf("body[%q] should not be set", tt.key)
				}
			}
		})
	}
}

func TestStartOfDay(t *testing.T) {
	if got := startOfDay("2026-05-03"); got != "2026-05-03 00:00:00" {
		t.Fatalf("startOfDay(date) = %q", got)
	}
	if got := startOfDay("2026-05-03 10:30:00"); got != "2026-05-03 10:30:00" {
		t.Fatalf("startOfDay(datetime) = %q", got)
	}
}

func TestEndOfDay(t *testing.T) {
	if got := endOfDay("2026-05-03"); got != "2026-05-03 23:59:59" {
		t.Fatalf("endOfDay(date) = %q", got)
	}
	if got := endOfDay("2026-05-03 10:30:00"); got != "2026-05-03 10:30:00" {
		t.Fatalf("endOfDay(datetime) = %q", got)
	}
}
