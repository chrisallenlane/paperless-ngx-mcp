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

func TestGetDocumentMetadata_Execute_WithArchive(t *testing.T) {
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
					"original_size": 1024,
					"original_mime_type": "application/pdf",
					"media_filename": "documents/0000001.pdf",
					"original_filename": "invoice.pdf",
					"original_metadata": [],
					"archive_checksum": "def456",
					"archive_size": 2048,
					"archive_media_filename": "archive/0000001.pdf",
					"archive_metadata": [],
					"has_archive_version": true,
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
		json.RawMessage(`{"id": 1}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Has Archive Version: true",
		"archive/0000001.pdf",
		"def456",
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
