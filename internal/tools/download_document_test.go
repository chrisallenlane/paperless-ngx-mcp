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

func TestDownloadDocument_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	savePath := filepath.Join(tmpDir, "downloaded.pdf")

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/documents/1/download/" {
					t.Errorf(
						"Expected /api/documents/1/download/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}

				w.Header().Set(
					"Content-Type",
					"application/pdf",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("fake pdf content here"))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewDownloadDocument(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"id": 1, "save_path": %q}`,
			savePath,
		)),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Document downloaded successfully",
		savePath,
		"application/pdf",
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

	content, err := os.ReadFile(savePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != "fake pdf content here" {
		t.Errorf(
			"File content = %q, want %q",
			string(content),
			"fake pdf content here",
		)
	}
}

func TestDownloadDocument_Execute_Original(t *testing.T) {
	tmpDir := t.TempDir()
	savePath := filepath.Join(tmpDir, "original.pdf")

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("original") != "true" {
					t.Errorf(
						"Expected original=true query param, got %s",
						r.URL.RawQuery,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/pdf",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("original content"))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"id": 1, "original": true, "save_path": %q}`,
			savePath,
		)),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDownloadDocument_Execute_MissingID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"save_path": "/tmp/test.pdf"}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing id")
	}

	if !strings.Contains(err.Error(), "positive integer") {
		t.Errorf(
			"Error should mention positive integer, got: %s",
			err.Error(),
		)
	}
}

func TestDownloadDocument_Execute_MissingSavePath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1}`),
	)
	if err == nil {
		t.Fatal("Expected error for missing save_path")
	}

	if !strings.Contains(err.Error(), "save_path is required") {
		t.Errorf(
			"Error should mention save_path, got: %s",
			err.Error(),
		)
	}
}

func TestDownloadDocument_Execute_RelativePath(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "save_path": "relative/path.pdf"}`,
		),
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

func TestDownloadDocument_Execute_NonexistentDir(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 1, "save_path": "/nonexistent/dir/file.pdf"}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for nonexistent directory")
	}

	if !strings.Contains(
		err.Error(),
		"parent directory does not exist",
	) {
		t.Errorf(
			"Error should mention parent directory, got: %s",
			err.Error(),
		)
	}
}

func TestDownloadDocument_Execute_ServerError(t *testing.T) {
	tmpDir := t.TempDir()
	savePath := filepath.Join(tmpDir, "test.pdf")

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Not Found"))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewDownloadDocument(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(fmt.Sprintf(
			`{"id": 999, "save_path": %q}`,
			savePath,
		)),
	)
	if err == nil {
		t.Fatal("Expected error for server error response")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf(
			"Error should mention status code, got: %s",
			err.Error(),
		)
	}
}
