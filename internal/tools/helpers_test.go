package tools

import (
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
