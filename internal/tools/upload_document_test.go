package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestUploadDocument_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(
		testFile,
		[]byte("fake pdf content"),
		0o644,
	); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/post_document/" {
					t.Errorf(
						"Expected /api/documents/post_document/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "POST" {
					t.Errorf("Expected POST, got %s", r.Method)
				}

				ct := r.Header.Get("Content-Type")
				if !strings.HasPrefix(
					ct,
					"multipart/form-data",
				) {
					t.Errorf(
						"Expected multipart/form-data content type, got %s",
						ct,
					)
				}

				if err := r.ParseMultipartForm(
					32 << 20,
				); err != nil {
					t.Fatalf(
						"Failed to parse multipart: %v",
						err,
					)
				}

				file, header, err := r.FormFile("document")
				if err != nil {
					t.Fatalf(
						"Missing document field: %v",
						err,
					)
				}
				defer file.Close()

				if header.Filename != "test.pdf" {
					t.Errorf(
						"Filename = %s, want test.pdf",
						header.Filename,
					)
				}

				if r.FormValue("title") != "My Test Doc" {
					t.Errorf(
						"title = %s, want My Test Doc",
						r.FormValue("title"),
					)
				}

				w.WriteHeader(http.StatusOK)
				w.Write(
					[]byte(`"abc-123-task-id"`),
				)
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewUploadDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"file_path": %q, "title": "My Test Doc"}`,
			testFile,
		)),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document uploaded successfully",
		"test.pdf",
		"Task ID:",
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

func TestUploadDocument_Execute_WithOptionalFields(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invoice.pdf")
	if err := os.WriteFile(
		testFile,
		[]byte("pdf content"),
		0o644,
	); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseMultipartForm(
					32 << 20,
				); err != nil {
					t.Fatalf(
						"Failed to parse multipart: %v",
						err,
					)
				}

				if r.FormValue("correspondent") != "1" {
					t.Errorf(
						"correspondent = %s, want 1",
						r.FormValue("correspondent"),
					)
				}

				if r.FormValue("document_type") != "2" {
					t.Errorf(
						"document_type = %s, want 2",
						r.FormValue("document_type"),
					)
				}

				tags := r.Form["tags"]
				if len(tags) != 2 {
					t.Errorf(
						"Expected 2 tags, got %d",
						len(tags),
					)
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"task-456"`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"file_path": %q, "correspondent": 1, "document_type": 2, "tags": [1, 3]}`,
			testFile,
		)),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestUploadDocument_Execute_MissingFilePath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing file_path")
	}

	if !strings.Contains(err.Error(), "file_path is required") {
		t.Errorf(
			"Error should mention file_path, got: %s",
			err.Error(),
		)
	}
}

func TestUploadDocument_Execute_RelativePath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"file_path": "relative/path.pdf"}`),
	)
	if err == nil {
		t.Fatal("Expected error for relative path")
	}

	if !strings.Contains(err.Error(), "absolute") {
		t.Errorf(
			"Error should mention absolute, got: %s",
			err.Error(),
		)
	}
}

func TestUploadDocument_Execute_TraversalPath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"file_path": "/tmp/../etc/passwd"}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for traversal path")
	}
}

func TestUploadDocument_Execute_NonexistentFile(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"file_path": "/tmp/nonexistent-file.pdf"}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestUploadDocument_Execute_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"file_path": %q}`,
			tmpDir,
		)),
	)
	if err == nil {
		t.Fatal("Expected error for directory path")
	}

	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf(
			"Error should mention directory, got: %s",
			err.Error(),
		)
	}
}

func TestUploadDocument_Execute_ServerError(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(
		testFile,
		[]byte("content"),
		0o644,
	); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

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
	tool := NewUploadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"file_path": %q}`,
			testFile,
		)),
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

func TestUploadDocument_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestUploadDocument_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewUploadDocument(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Schema should have required field")
	}

	foundFilePath := false
	for _, r := range required {
		if r == "file_path" {
			foundFilePath = true
		}
	}
	if !foundFilePath {
		t.Error("file_path should be in required fields")
	}
}
