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

func TestGetCustomField_Execute_NullExtraData(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 2,
					"name": "Due Date",
					"data_type": "date",
					"extra_data": null,
					"document_count": 0
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
	tool := NewGetCustomField(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 2}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "Extra Data: (none)") {
		t.Errorf(
			"Output missing 'Extra Data: (none)'.\nGot:\n%s",
			result,
		)
	}
}
