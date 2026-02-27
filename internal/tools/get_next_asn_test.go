package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestGetNextASN_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/next_asn/" {
					t.Errorf(
						"Expected /api/documents/next_asn/, got %s",
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
				w.Write([]byte(`43`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetNextASN(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "Next available ASN: 43"
	if result != expected {
		t.Errorf("Result = %q, want %q", result, expected)
	}
}
