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

func TestCreateSavedView_MissingName(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateSavedView(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"show_on_dashboard": true, "show_in_sidebar": false, "filter_rules": []}`,
		),
	)
	if err == nil {
		t.Fatal("Expected error for missing name")
	}

	if !strings.Contains(
		err.Error(),
		"name is required",
	) {
		t.Errorf(
			"Error should mention name is required, got: %s",
			err.Error(),
		)
	}
}

func TestCreateSavedView_MissingShowOnDashboard(
	t *testing.T,
) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateSavedView(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Test", "show_in_sidebar": false, "filter_rules": []}`,
		),
	)
	if err == nil {
		t.Fatal(
			"Expected error for missing show_on_dashboard",
		)
	}

	if !strings.Contains(
		err.Error(),
		"show_on_dashboard is required",
	) {
		t.Errorf(
			"Error should mention show_on_dashboard, got: %s",
			err.Error(),
		)
	}
}

func TestCreateSavedView_MissingShowInSidebar(
	t *testing.T,
) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateSavedView(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Test", "show_on_dashboard": true, "filter_rules": []}`,
		),
	)
	if err == nil {
		t.Fatal(
			"Expected error for missing show_in_sidebar",
		)
	}

	if !strings.Contains(
		err.Error(),
		"show_in_sidebar is required",
	) {
		t.Errorf(
			"Error should mention show_in_sidebar, got: %s",
			err.Error(),
		)
	}
}

func TestCreateSavedView_MissingFilterRules(
	t *testing.T,
) {
	c := client.New("http://localhost", "test-token")
	tool := NewCreateSavedView(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Test", "show_on_dashboard": true, "show_in_sidebar": false}`,
		),
	)
	if err == nil {
		t.Fatal(
			"Expected error for missing filter_rules",
		)
	}

	if !strings.Contains(
		err.Error(),
		"filter_rules is required",
	) {
		t.Errorf(
			"Error should mention filter_rules, got: %s",
			err.Error(),
		)
	}
}

func TestCreateSavedView_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				if r.URL.Path != "/api/saved_views/" {
					t.Errorf(
						"Expected /api/saved_views/, got %s",
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

				if req["name"] != "Unpaid Invoices" {
					t.Errorf(
						"name = %v, want Unpaid Invoices",
						req["name"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"id": 1,
					"name": "Unpaid Invoices",
					"show_on_dashboard": true,
					"show_in_sidebar": false,
					"sort_field": "created",
					"sort_reverse": true,
					"filter_rules": [
						{"rule_type": 4, "value": "3"}
					],
					"page_size": null,
					"display_mode": "table"
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
	tool := NewCreateSavedView(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{
				"name": "Unpaid Invoices",
				"show_on_dashboard": true,
				"show_in_sidebar": false,
				"sort_field": "created",
				"sort_reverse": true,
				"filter_rules": [
					{"rule_type": 4, "value": "3"}
				],
				"display_mode": "table"
			}`,
		),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Saved View (ID: 1)",
		"Name: Unpaid Invoices",
		"Show on Dashboard: true",
		"Show in Sidebar: false",
		"Sort Field: created",
		"Display Mode: table",
		"Document type is: 3",
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
