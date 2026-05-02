package yunxiao

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestHandleGetPipelineScanReportURLBuildsQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/flow/organizations/org-1/pipelines/getPipelineScanReportUrl?reportPath=%2Freports%2Fscan.html" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`"https://example.com/scan.html"`))
	})

	result, err := handleGetPipelineScanReportURL(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"reportPath":     "/reports/scan.html",
	})
	if err != nil {
		t.Fatalf("handleGetPipelineScanReportURL() error = %v", err)
	}
	if !strings.Contains(result, "scan.html") {
		t.Fatalf("result = %q", result)
	}
}

func TestHandleGetPipelineArtifactURLBuildsQuery(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/flow/organizations/org-1/pipelines/getArtifactDownloadUrl?fileName=build.tgz&filePath=%2Fartifacts%2Fbuild.tgz" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`"https://example.com/build.tgz"`))
	})

	if _, err := handleGetPipelineArtifactURL(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"filePath":       "/artifacts/build.tgz",
		"fileName":       "build.tgz",
	}); err != nil {
		t.Fatalf("handleGetPipelineArtifactURL() error = %v", err)
	}
}

func TestHandleGetPipelineEmasArtifactURLPreservesInt64IDs(t *testing.T) {
	const maxInt64 = "9223372036854775807"

	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/oapi/v1/flow/organizations/org-1/pipelines/getEmasArtifactDownloadUrl" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("emasJobInstanceId") != "emas-1" ||
			r.URL.Query().Get("md5") != "abc123" ||
			r.URL.Query().Get("pipelineId") != maxInt64 ||
			r.URL.Query().Get("pipelineRunId") != maxInt64 ||
			r.URL.Query().Get("serviceConnectionId") != maxInt64 {
			t.Fatalf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`"https://example.com/app.apk"`))
	})

	if _, err := handleGetPipelineEmasArtifactURL(context.Background(), client, map[string]any{
		"organizationId":      "org-1",
		"emasJobInstanceId":   "emas-1",
		"md5":                 "abc123",
		"pipelineId":          maxInt64,
		"pipelineRunId":       maxInt64,
		"serviceConnectionId": maxInt64,
	}); err != nil {
		t.Fatalf("handleGetPipelineEmasArtifactURL() error = %v", err)
	}
}

func TestHandleListPipelineRelationsBuildsEscapedPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/flow/organizations/org-1/pipelines/pipe%2F1/pipelineObjRel/VARIABLE_GROUP/list" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`[{"refObjectId":"vg-1"}]`))
	})

	if _, err := handleListPipelineRelations(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe/1",
		"relObjectType":  "VARIABLE_GROUP",
	}); err != nil {
		t.Fatalf("handleListPipelineRelations() error = %v", err)
	}
}

func TestHandleGetLastInstanceBuildsPath(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.RequestURI != "/oapi/v1/flow/organizations/org-1/createServiceConnection/pipe%2F1/getLastInstance" {
			t.Fatalf("RequestURI = %q", r.RequestURI)
		}
		_, _ = w.Write([]byte(`{"last":1,"content":"ok","more":false}`))
	})

	if _, err := handleGetLastInstance(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"pipelineId":     "pipe/1",
	}); err != nil {
		t.Fatalf("handleGetLastInstance() error = %v", err)
	}
}

func TestFlowArtifactRelationHandlersRequireParams(t *testing.T) {
	client := newHandlerTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.RequestURI)
	})

	if _, err := handleGetPipelineScanReportURL(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	_, err := handleGetPipelineArtifactURL(context.Background(), client, map[string]any{
		"organizationId": "org-1",
		"filePath":       "/artifacts/build.tgz",
	})
	if err == nil {
		t.Fatal("handleGetPipelineArtifactURL() expected missing fileName error")
	}
	_, err = handleGetPipelineEmasArtifactURL(context.Background(), client, map[string]any{
		"organizationId":    "org-1",
		"emasJobInstanceId": "emas-1",
		"md5":               "abc123",
		"pipelineId":        "pipe-1",
		"pipelineRunId":     "run-1",
	})
	if err == nil {
		t.Fatal("handleGetPipelineEmasArtifactURL() expected missing serviceConnectionId error")
	}
	if _, err := handleListPipelineRelations(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
	if _, err := handleGetLastInstance(context.Background(), client, map[string]any{}); err == nil {
		t.Fatal("expected missing params error")
	}
}
