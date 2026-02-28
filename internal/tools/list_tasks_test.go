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

const taskArrayResponse = `[
	{
		"id": 1,
		"task_id": "abc-123",
		"task_name": "consume_file",
		"task_file_name": "invoice.pdf",
		"date_created": "2026-02-27T10:00:00Z",
		"date_done": "2026-02-27T10:01:00Z",
		"type": "auto_task",
		"status": "SUCCESS",
		"result": "Success",
		"acknowledged": false,
		"related_document": "42"
	},
	{
		"id": 2,
		"task_id": "def-456",
		"task_name": "train_classifier",
		"task_file_name": null,
		"date_created": "2026-02-27T11:00:00Z",
		"date_done": null,
		"type": "scheduled_task",
		"status": "STARTED",
		"result": null,
		"acknowledged": false,
		"related_document": null
	}
]`

func TestListTasks_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/tasks/" {
					t.Errorf(
						"Expected /api/tasks/, got %s",
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
				w.Write([]byte(taskArrayResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Tasks: 2 total",
		"consume_file",
		"SUCCESS",
		"invoice.pdf",
		"train_classifier",
		"STARTED",
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

func TestListTasks_StatusFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				got := r.URL.Query().Get("status")
				if got != "SUCCESS" {
					t.Errorf(
						"Expected status=SUCCESS, got %s",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"status": "SUCCESS"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListTasks_TaskNameFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				got := r.URL.Query().Get("task_name")
				if got != "consume_file" {
					t.Errorf(
						"Expected task_name=consume_file, got %s",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"task_name": "consume_file"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListTasks_TypeFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				got := r.URL.Query().Get("type")
				if got != "auto_task" {
					t.Errorf(
						"Expected type=auto_task, got %s",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"type": "auto_task"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListTasks_TaskIDFilter(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				got := r.URL.Query().Get("task_id")
				if got != "abc-123" {
					t.Errorf(
						"Expected task_id=abc-123, got %s",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"task_id": "abc-123"}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestListTasks_Empty(t *testing.T) {
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
				w.Write([]byte(`[]`))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewListTasks(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No tasks found." {
		t.Errorf(
			"Expected 'No tasks found.', got: %s",
			result,
		)
	}
}
