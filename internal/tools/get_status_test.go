package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

const statusResponse = `{
	"pngx_version": "2.20.8",
	"server_os": "Linux-4.4.302+-x86_64-with-glibc2.41",
	"install_type": "docker",
	"storage": {
		"total": 11518122557440,
		"available": 8525312483328
	},
	"database": {
		"type": "sqlite",
		"url": "/data/db.sqlite3",
		"status": "OK",
		"error": null,
		"migration_status": {
			"latest_migration": "documents.0042_auto",
			"unapplied_migrations": []
		}
	},
	"tasks": {
		"redis_url": "redis://redis:6379",
		"redis_status": "OK",
		"redis_error": null,
		"celery_status": "OK",
		"celery_url": "celery@worker",
		"celery_error": null,
		"index_status": "OK",
		"index_last_modified": "2026-02-27T12:00:00Z",
		"index_error": null,
		"classifier_status": "OK",
		"classifier_last_trained": "2026-02-27T10:00:00Z",
		"classifier_error": null,
		"sanity_check_status": "OK",
		"sanity_check_last_run": "2026-02-22T06:00:00Z",
		"sanity_check_error": null
	}
}`

func TestGetStatus_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/status/" {
					t.Errorf(
						"Expected /api/status/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(statusResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetStatus(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify key parts of the output
	checks := []string{
		"Paperless-NGX Status",
		"Version: 2.20.8",
		"OS: Linux-4.4.302+-x86_64-with-glibc2.41",
		"Install: docker",
		"Database: sqlite - OK",
		"Redis: OK",
		"Celery: OK",
		"Index: OK (last modified: 2026-02-27)",
		"Classifier: OK (last trained: 2026-02-27)",
		"Sanity Check: OK (last run: 2026-02-22)",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf(
				"Output missing %q.\nGot:\n%s",
				check,
				result,
			)
		}
	}
}
