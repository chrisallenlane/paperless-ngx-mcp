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

func TestGetDocumentMetadata_Execute_NoArchive(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"original_checksum": "abc123",
					"original_size": 512,
					"original_mime_type": "image/png",
					"media_filename": "documents/0000002.png",
					"original_filename": "photo.png",
					"original_metadata": [],
					"archive_checksum": "",
					"archive_size": 0,
					"archive_media_filename": "",
					"archive_metadata": [],
					"has_archive_version": false,
					"lang": "en"
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
	tool := NewGetDocumentMetadata(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 2}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(
		result,
		"Has Archive Version: false",
	) {
		t.Errorf(
			"Output should show no archive.\nGot:\n%s",
			result,
		)
	}
}
