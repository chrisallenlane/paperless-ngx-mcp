package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

const noteListResponse = `{
	"count": 2,
	"next": null,
	"previous": null,
	"all": [1, 2],
	"results": [
		{
			"id": 1,
			"note": "Payment received",
			"created": "2026-02-15T10:30:00Z",
			"user": {
				"id": 1,
				"username": "admin",
				"first_name": "Admin",
				"last_name": "User"
			}
		},
		{
			"id": 2,
			"note": "Forwarded to accounting",
			"created": "2026-02-16T14:00:00Z",
			"user": {
				"id": 2,
				"username": "jane",
				"first_name": "Jane",
				"last_name": "Doe"
			}
		}
	]
}`

const emptyNoteListResponse = `{
	"count": 0,
	"next": null,
	"previous": null,
	"all": [],
	"results": []
}`

// --- ListDocumentNotes tests ---

func TestListDocumentNotes_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/documents/42/notes/" {
					t.Errorf(
						"Expected /api/documents/42/notes/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf(
						"Expected GET, got %s",
						r.Method,
					)
				}
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(noteListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocumentNotes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 42}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Notes for Document 42: 2 total",
		"Note 1",
		"admin",
		"Payment received",
		"Note 2",
		"jane",
		"Forwarded to accounting",
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

func TestListDocumentNotes_Empty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(emptyNoteListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocumentNotes(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 42}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "No notes found for document 42."
	if result != expected {
		t.Errorf(
			"Expected %q, got: %s",
			expected,
			result,
		)
	}
}

func TestListDocumentNotes_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewListDocumentNotes(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 0}`),
	)
	if err == nil {
		t.Fatal("Expected error for invalid id")
	}

	if !strings.Contains(
		err.Error(),
		"positive integer",
	) {
		t.Errorf(
			"Error should mention positive integer, got: %s",
			err.Error(),
		)
	}
}

func TestListDocumentNotes_Pagination(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				q := r.URL.Query()
				if got := q.Get("page"); got != "2" {
					t.Errorf(
						"Expected page=2, got %s",
						got,
					)
				}
				if got := q.Get("page_size"); got != "5" {
					t.Errorf(
						"Expected page_size=5, got %s",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(emptyNoteListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListDocumentNotes(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 42, "page": 2, "page_size": 5}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// --- CreateDocumentNote tests ---

func TestCreateDocumentNote_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/documents/42/notes/" {
					t.Errorf(
						"Expected /api/documents/42/notes/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "POST" {
					t.Errorf(
						"Expected POST, got %s",
						r.Method,
					)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf(
						"Failed to read body: %v",
						err,
					)
				}

				var req map[string]interface{}
				if err := json.Unmarshal(
					body,
					&req,
				); err != nil {
					t.Fatalf(
						"Failed to parse body: %v",
						err,
					)
				}

				if req["note"] != "Payment received" {
					t.Errorf(
						"note = %v, want Payment received",
						req["note"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(noteListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewCreateDocumentNote(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 42, "note": "Payment received"}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Note added to document 42",
		"Notes for Document 42: 2 total",
		"Payment received",
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

func TestCreateDocumentNote_MissingNote(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentNote(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"id": 1, "note": ""}`),
	)
	if err == nil {
		t.Fatal("Expected error for empty note")
	}

	if !strings.Contains(
		err.Error(),
		"note is required",
	) {
		t.Errorf(
			"Error should mention note is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateDocumentNote_InvalidID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateDocumentNote(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"id": 0, "note": "Test"}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for invalid id")
	}

	if !strings.Contains(
		err.Error(),
		"positive integer",
	) {
		t.Errorf(
			"Error should mention positive integer, got: %s",
			err.Error(),
		)
	}
}

// --- DeleteDocumentNote tests ---

func TestDeleteDocumentNote_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/documents/42/notes/" {
					t.Errorf(
						"Expected /api/documents/42/notes/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "DELETE" {
					t.Errorf(
						"Expected DELETE, got %s",
						r.Method,
					)
				}

				noteID := r.URL.Query().Get("id")
				if noteID != "7" {
					t.Errorf(
						"Expected note id=7, got %s",
						noteID,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(emptyNoteListResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewDeleteDocumentNote(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"document_id": 42, "note_id": 7}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Note 7 deleted from document 42",
		"No notes found for document 42.",
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

func TestDeleteDocumentNote_InvalidDocumentID(
	t *testing.T,
) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocumentNote(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"document_id": 0, "note_id": 1}`,
		),
	)
	if err == nil {
		t.Fatal(
			"Expected error for invalid document_id",
		)
	}

	if !strings.Contains(
		err.Error(),
		"document_id must be a positive integer",
	) {
		t.Errorf(
			"Error should mention document_id, got: %s",
			err.Error(),
		)
	}
}

func TestDeleteDocumentNote_InvalidNoteID(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewDeleteDocumentNote(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"document_id": 1, "note_id": 0}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for invalid note_id")
	}

	if !strings.Contains(
		err.Error(),
		"note_id must be a positive integer",
	) {
		t.Errorf(
			"Error should mention note_id, got: %s",
			err.Error(),
		)
	}
}
