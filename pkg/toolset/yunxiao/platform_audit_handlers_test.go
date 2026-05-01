package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleListAuditLogsBuildsPathQueryAndNextToken(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/platform/organizations/org-1/auditLogs" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("actionTimeStart") != "2026-05-01T00:00:00Z" ||
			r.URL.Query().Get("actionTimeEnd") != "2026-05-01T01:00:00Z" ||
			r.URL.Query().Get("userIds") != "user-1,user-2" ||
			r.URL.Query().Get("apps") != "codeup,flow" ||
			r.URL.Query().Get("perPage") != "50" ||
			r.URL.Query().Get("nextToken") != "token-1" {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		w.Header().Set("x-next-token", "token-2")
		_, _ = w.Write([]byte(`[{"action":"repo.read"}]`))
	})

	result, err := handleListAuditLogs(context.Background(), client, map[string]any{
		"organizationId":  "org-1",
		"actionTimeStart": "2026-05-01T00:00:00Z",
		"actionTimeEnd":   "2026-05-01T01:00:00Z",
		"userIds":         "user-1,user-2",
		"apps":            "codeup,flow",
		"perPage":         float64(50),
		"nextToken":       "token-1",
	})
	if err != nil {
		t.Fatalf("handleListAuditLogs() error = %v", err)
	}
	if !strings.Contains(result, `"nextToken": "token-2"`) {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleListAuditLogsRequiresActionTimeStart(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	_, err := handleListAuditLogs(context.Background(), client, map[string]any{
		"organizationId": "org-1",
	})
	if err == nil {
		t.Fatal("handleListAuditLogs() expected missing actionTimeStart error")
	}
}
