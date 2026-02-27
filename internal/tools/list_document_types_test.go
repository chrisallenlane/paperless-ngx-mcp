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

func TestListDocumentTypes_Execute_WithPagination(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"count": 30,
					"next": "http://example.com/api/document_types/?page=2",
					"previous": null,
					"all": [1],
					"results": [{
						"id": 1,
						"slug": "invoice",
						"name": "Invoice",
						"match": "",
						"matching_algorithm": 1,
						"is_insensitive": true,
						"document_count": 10
					}]
				}`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocumentTypes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "more results available") {
		t.Errorf(
			"Output should show pagination hint.\nGot:\n%s",
			result,
		)
	}
}
