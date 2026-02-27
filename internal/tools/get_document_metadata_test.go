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

const documentMetadataResponse = `{
	"original_checksum": "abc123def456",
	"original_size": 102400,
	"original_mime_type": "application/pdf",
	"media_filename": "documents/0000001.pdf",
	"original_filename": "invoice-2024.pdf",
	"original_metadata": [],
	"archive_checksum": "xyz789",
	"archive_size": 204800,
	"archive_media_filename": "documents/0000001-archive.pdf",
	"archive_metadata": [],
	"has_archive_version": true,
	"lang": "en"
}`

func TestGetDocumentMetadata_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/metadata/" {
					t.Errorf(
						"Expected /api/documents/1/metadata/, got %s",
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
				w.Write([]byte(documentMetadataResponse))
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
		"Document Metadata (ID: 1)",
		"Filename: invoice-2024.pdf",
		"MIME Type: application/pdf",
		"Checksum: abc123def456",
		"Has Archive Version: true",
		"OCR Language: en",
		"100.00 KB",
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

func TestGetDocumentMetadata_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentMetadata(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing id")
	}
}

func TestGetDocumentMetadata_Execute_ServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
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

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err == nil {
		t.Fatal("Expected error for server error response")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"Error should mention status code, got: %s",
			err.Error(),
		)
	}
}

func TestGetDocumentMetadata_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentMetadata(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestGetDocumentMetadata_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetDocumentMetadata(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf(
			"Schema type = %v, want object",
			schema["type"],
		)
	}
}
