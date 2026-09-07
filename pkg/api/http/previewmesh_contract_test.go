package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stefanprodan/podinfo/pkg/version"
)

func TestPreviewMeshContract(t *testing.T) {
	sha := strings.Repeat("a", 40)
	t.Setenv("PREVIEW_COMMIT_SHA", sha)
	atomic.StoreInt32(&healthy, 1)
	t.Cleanup(func() { atomic.StoreInt32(&healthy, 0) })

	srv := NewMockServer()
	srv.registerHandlers()

	root := httptest.NewRecorder()
	srv.router.ServeHTTP(root, httptest.NewRequest(http.MethodGet, "/", nil))
	if root.Code != http.StatusOK {
		t.Fatalf("GET /: got %d, want %d", root.Code, http.StatusOK)
	}
	var rootBody RuntimeResponse
	if err := json.Unmarshal(root.Body.Bytes(), &rootBody); err != nil {
		t.Fatalf("GET /: invalid JSON: %v", err)
	}
	if rootBody.App != "podinfo" || rootBody.CommitSHA != sha {
		t.Fatalf("GET /: got app=%q commit_sha=%q", rootBody.App, rootBody.CommitSHA)
	}

	health := httptest.NewRecorder()
	srv.router.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("GET /health: got %d, want %d", health.Code, http.StatusOK)
	}
	var healthBody map[string]string
	if err := json.Unmarshal(health.Body.Bytes(), &healthBody); err != nil {
		t.Fatalf("GET /health: invalid JSON: %v", err)
	}
	if healthBody["status"] != "ok" || healthBody["commit_sha"] != sha {
		t.Fatalf("GET /health: got %#v", healthBody)
	}

	unknown := httptest.NewRecorder()
	srv.router.ServeHTTP(unknown, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if unknown.Code != http.StatusNotFound {
		t.Fatalf("GET /missing: got %d, want %d", unknown.Code, http.StatusNotFound)
	}

	t.Setenv("PREVIEW_COMMIT_SHA", "")
	fallback := httptest.NewRecorder()
	srv.router.ServeHTTP(fallback, httptest.NewRequest(http.MethodGet, "/health", nil))
	var fallbackBody map[string]string
	if err := json.Unmarshal(fallback.Body.Bytes(), &fallbackBody); err != nil {
		t.Fatalf("fallback health: invalid JSON: %v", err)
	}
	want := version.REVISION
	if want == "" {
		want = "unknown"
	}
	if fallbackBody["commit_sha"] != want {
		t.Fatalf("fallback health: got commit_sha=%q, want %q", fallbackBody["commit_sha"], want)
	}
}
