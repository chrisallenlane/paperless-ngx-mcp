package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// unmarshalable is a value that cannot be serialized to JSON because it
// contains a function field, which json.Marshal rejects.
var unmarshalable = map[string]interface{}{"fn": func() {}}

func TestNew(t *testing.T) {
	c := New("https://api.example.com", "my-token")

	if c.BaseURL != "https://api.example.com" {
		t.Errorf(
			"Expected BaseURL to be https://api.example.com, got %s",
			c.BaseURL,
		)
	}

	if c.Token != "my-token" {
		t.Errorf("Expected Token to be my-token, got %s", c.Token)
	}

	if c.HTTPClient == nil {
		t.Error("Expected HTTPClient to be set")
	}
}

func TestAuthHeaderSent(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Token test-token-123" {
				t.Errorf(
					"Expected Authorization header 'Token test-token-123', got %q",
					authHeader,
				)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(
		server.URL,
		"test-token-123",
		server.Client(),
	)

	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/test" {
				t.Errorf("Expected /test, got %s", r.URL.Path)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": "success"}`))
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(server.URL, "token", server.Client())

	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST, got %s", r.Method)
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf(
					"Expected Content-Type application/json, got %s",
					contentType,
				)
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": 1}`))
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(server.URL, "token", server.Client())

	resp, err := c.Post(
		context.Background(),
		"/test",
		map[string]string{"key": "value"},
	)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201, got %d", resp.StatusCode)
	}
}

func TestPatch(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Expected PATCH, got %s", r.Method)
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader != "Token token" {
				t.Errorf(
					"Expected Authorization header, got %q",
					authHeader,
				)
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf(
					"Expected Content-Type application/json, got %s",
					contentType,
				)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1}`))
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(server.URL, "token", server.Client())

	resp, err := c.Patch(
		context.Background(),
		"/test/1",
		map[string]string{"key": "patched"},
	)
	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Expected DELETE, got %s", r.Method)
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(server.URL, "token", server.Client())

	resp, err := c.Delete(context.Background(), "/test/1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", resp.StatusCode)
	}
}

func TestPostMultipart(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST, got %s", r.Method)
			}

			contentType := r.Header.Get("Content-Type")
			if !strings.HasPrefix(
				contentType,
				"multipart/form-data",
			) {
				t.Errorf(
					"Expected multipart/form-data, got %s",
					contentType,
				)
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader != "Token token" {
				t.Errorf(
					"Expected Authorization header, got %q",
					authHeader,
				)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"task-id"`))
		}),
	)
	defer server.Close()

	c := NewWithHTTPClient(server.URL, "token", server.Client())

	body := strings.NewReader("fake multipart content")
	resp, err := c.PostMultipart(
		context.Background(),
		"/test/upload",
		body,
		"multipart/form-data; boundary=test",
	)
	if err != nil {
		t.Fatalf("PostMultipart failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

// TestTransportErrors verifies that Get, Post, Patch, and Delete each
// return an error when the HTTP client cannot reach the server at all
// (transport-level failure, no response).
func TestTransportErrors(t *testing.T) {
	// Port 1 on loopback is reserved and will refuse connections
	// reliably without requiring an open port.
	c := New("http://127.0.0.1:1", "token")

	tests := []struct {
		name string
		call func() (*http.Response, error)
	}{
		{
			name: "Get",
			call: func() (*http.Response, error) {
				return c.Get(context.Background(), "/api/test/")
			},
		},
		{
			name: "Post",
			call: func() (*http.Response, error) {
				return c.Post(
					context.Background(),
					"/api/test/",
					map[string]string{"key": "value"},
				)
			},
		},
		{
			name: "Patch",
			call: func() (*http.Response, error) {
				return c.Patch(
					context.Background(),
					"/api/test/1/",
					map[string]string{"key": "value"},
				)
			},
		},
		{
			name: "Delete",
			call: func() (*http.Response, error) {
				return c.Delete(context.Background(), "/api/test/1/")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.call()
			if err == nil {
				resp.Body.Close()
				t.Errorf("%s: expected transport error, got nil", tt.name)
			}
		})
	}
}

// TestMarshalErrors verifies that Post and Patch return an error containing
// "failed to marshal request body" when given a body that cannot be
// serialized to JSON (e.g., a map containing a function value).
func TestMarshalErrors(t *testing.T) {
	// The server URL is irrelevant; json.Marshal runs before any
	// network call is made, so we never actually connect.
	c := New("http://127.0.0.1:1", "token")

	tests := []struct {
		name string
		call func() (*http.Response, error)
	}{
		{
			name: "Post",
			call: func() (*http.Response, error) {
				return c.Post(
					context.Background(),
					"/api/test/",
					unmarshalable,
				)
			},
		},
		{
			name: "Patch",
			call: func() (*http.Response, error) {
				return c.Patch(
					context.Background(),
					"/api/test/1/",
					unmarshalable,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.call()
			if err == nil {
				resp.Body.Close()
				t.Fatalf(
					"%s: expected marshal error, got nil",
					tt.name,
				)
			}

			if !strings.Contains(
				err.Error(),
				"failed to marshal request body",
			) {
				t.Errorf(
					"%s: expected error containing %q, got %q",
					tt.name,
					"failed to marshal request body",
					err.Error(),
				)
			}
		})
	}
}
