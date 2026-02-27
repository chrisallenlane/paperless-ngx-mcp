package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestReadResponse_OK(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"id": 1}`)),
	}

	body, err := readResponse(resp, http.StatusOK)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(body) != `{"id": 1}` {
		t.Errorf("Body = %s, want {\"id\": 1}", string(body))
	}
}

func TestReadResponse_Created(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(strings.NewReader(`{"id": 1}`)),
	}

	body, err := readResponse(resp, http.StatusCreated)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(body) != `{"id": 1}` {
		t.Errorf("Body = %s, want {\"id\": 1}", string(body))
	}
}

func TestReadResponse_NoContent(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	_, err := readResponse(resp, http.StatusNoContent)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestReadResponse_Error(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body: io.NopCloser(
			strings.NewReader("Internal Server Error"),
		),
	}

	_, err := readResponse(resp, http.StatusOK)
	if err == nil {
		t.Fatal("Expected error for non-200 response")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"Error should mention status code, got: %s",
			err.Error(),
		)
	}
}

func TestBuildListPath_WithPageParams(t *testing.T) {
	args := json.RawMessage(
		`{"page": 2, "page_size": 10}`,
	)

	path, err := buildListPath("/api/test/", args)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(path, "page=2") {
		t.Errorf("Path should contain page=2, got: %s", path)
	}

	if !strings.Contains(path, "page_size=10") {
		t.Errorf(
			"Path should contain page_size=10, got: %s",
			path,
		)
	}
}

func TestBuildListPath_NoParams(t *testing.T) {
	args := json.RawMessage(`{}`)

	path, err := buildListPath("/api/test/", args)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if path != "/api/test/" {
		t.Errorf("Path = %s, want /api/test/", path)
	}
}
