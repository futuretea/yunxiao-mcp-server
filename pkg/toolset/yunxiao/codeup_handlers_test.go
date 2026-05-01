package yunxiao

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListChangeRequestsBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oapi/v1/codeup/organizations/org-1/changeRequests" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("projectIds") != "group/repo,repo-2" ||
			r.URL.Query().Get("state") != "opened" ||
			r.URL.Query().Get("orderBy") != "updated_at" ||
			r.URL.Query().Get("sort") != "desc" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-total", "1")
		_, _ = w.Write([]byte(`[{"localId":1}]`))
	})

	result, err := handleListChangeRequests(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"projectIds":     "group/repo,repo-2",
		"state":          "opened",
		"orderBy":        "updated_at",
		"sort":           "desc",
	})
	if err != nil {
		t.Fatalf("handleListChangeRequests() error = %v", err)
	}
	if !strings.Contains(result, `"pagination"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetChangeRequestBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"localId":12}`))
	})

	if _, err := handleGetChangeRequest(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"localId":        "12",
	}); err != nil {
		t.Fatalf("handleGetChangeRequest() error = %v", err)
	}
}

func TestHandleListChangeRequestPatchSetsBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12/diffs/patches") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"patchSetBizId":"p1"}]`))
	})

	if _, err := handleListChangeRequestPatchSets(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"localId":        "12",
	}); err != nil {
		t.Fatalf("handleListChangeRequestPatchSets() error = %v", err)
	}
}

func TestHandleGetChangeRequestTreeBuildsPathAndQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12/diffs/changeTree") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		if r.URL.Query().Get("fromPatchSetId") != "from-1" ||
			r.URL.Query().Get("toPatchSetId") != "to-1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"files":[]}`))
	})

	if _, err := handleGetChangeRequestTree(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"localId":        "12",
		"fromPatchSetId": "from-1",
		"toPatchSetId":   "to-1",
	}); err != nil {
		t.Fatalf("handleGetChangeRequestTree() error = %v", err)
	}
}

func TestHandleListChangeRequestCommentsBuildsBody(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12/comments/list") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["comment_type"] != "INLINE_COMMENT" || body["state"] != "OPENED" || body["resolved"] != true || body["file_path"] != "src/main.go" {
			t.Fatalf("body = %#v", body)
		}
		if body["patchset_biz_id_list"] != "p1,p2" {
			t.Fatalf("patchset_biz_id_list = %#v", body["patchset_biz_id_list"])
		}
		_, _ = w.Write([]byte(`[{"commentBizId":"c1"}]`))
	})

	if _, err := handleListChangeRequestComments(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"localId":        "12",
		"patchSetBizIds": "p1,p2",
		"commentType":    "INLINE_COMMENT",
		"resolved":       true,
		"filePath":       "src/main.go",
	}); err != nil {
		t.Fatalf("handleListChangeRequestComments() error = %v", err)
	}
}

func TestHandleGetChangeRequestCommentBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if !strings.Contains(r.RequestURI, "/repositories/group%2Frepo/changeRequests/12/comments/comment-1") {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"commentBizId":"comment-1"}`))
	})

	if _, err := handleGetChangeRequestComment(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"repositoryId":   "group/repo",
		"localId":        "12",
		"commentBizId":   "comment-1",
	}); err != nil {
		t.Fatalf("handleGetChangeRequestComment() error = %v", err)
	}
}
