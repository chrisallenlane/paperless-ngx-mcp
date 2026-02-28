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

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "absolute path",
			path:    "/home/user/documents/file.pdf",
			wantErr: false,
		},
		{
			name:    "absolute path at root",
			path:    "/file.pdf",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "relative/path/file.pdf",
			wantErr: true,
		},
		{
			name:    "relative path with dot prefix",
			path:    "./file.pdf",
			wantErr: true,
		},
		{
			name:    "path with traversal sequences",
			path:    "/foo/../../../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "path with embedded traversal",
			path:    "/safe/../../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "empty string",
			path:    "",
			wantErr: true,
		},
		{
			name:    "bare double-dot",
			path:    "..",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateFilePath(%q) error = %v, wantErr %v",
					tt.path,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestParseIDArg(t *testing.T) {
	tests := []struct {
		name    string
		args    json.RawMessage
		wantID  int
		wantErr bool
	}{
		{
			name:    "valid positive integer",
			args:    json.RawMessage(`{"id": 42}`),
			wantID:  42,
			wantErr: false,
		},
		{
			name:    "id of 1",
			args:    json.RawMessage(`{"id": 1}`),
			wantID:  1,
			wantErr: false,
		},
		{
			name:    "zero returns error",
			args:    json.RawMessage(`{"id": 0}`),
			wantErr: true,
		},
		{
			name:    "negative integer returns error",
			args:    json.RawMessage(`{"id": -5}`),
			wantErr: true,
		},
		{
			name:    "non-integer value returns error",
			args:    json.RawMessage(`{"id": "abc"}`),
			wantErr: true,
		},
		{
			name:    "missing id field returns error",
			args:    json.RawMessage(`{}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := parseIDArg(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"parseIDArg(%s) error = %v, wantErr %v",
					tt.args,
					err,
					tt.wantErr,
				)
			}
			if !tt.wantErr && gotID != tt.wantID {
				t.Errorf(
					"parseIDArg(%s) = %d, want %d",
					tt.args,
					gotID,
					tt.wantID,
				)
			}
		})
	}
}

func TestParsePatchArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         json.RawMessage
		wantID       int
		wantBodyKeys []string
		wantIDInBody bool
		wantErr      bool
	}{
		{
			name: "id and extra fields",
			args: json.RawMessage(
				`{"id": 7, "name": "foo", "color": "#fff"}`,
			),
			wantID:       7,
			wantBodyKeys: []string{"name", "color"},
			wantErr:      false,
		},
		{
			name:    "id only produces empty patch body",
			args:    json.RawMessage(`{"id": 3}`),
			wantID:  3,
			wantErr: false,
		},
		{
			name:    "missing id returns error",
			args:    json.RawMessage(`{"name": "foo"}`),
			wantErr: true,
		},
		{
			name:    "zero id returns error",
			args:    json.RawMessage(`{"id": 0, "name": "foo"}`),
			wantErr: true,
		},
		{
			name:    "negative id returns error",
			args:    json.RawMessage(`{"id": -1, "name": "foo"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotBody, err := parsePatchArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"parsePatchArgs(%s) error = %v, wantErr %v",
					tt.args,
					err,
					tt.wantErr,
				)
				return
			}

			if tt.wantErr {
				return
			}

			if gotID != tt.wantID {
				t.Errorf(
					"parsePatchArgs(%s) id = %d, want %d",
					tt.args,
					gotID,
					tt.wantID,
				)
			}

			// id must never appear in the patch body
			if _, ok := gotBody["id"]; ok {
				t.Errorf(
					"parsePatchArgs(%s) patch body must not contain 'id'",
					tt.args,
				)
			}

			for _, key := range tt.wantBodyKeys {
				if _, ok := gotBody[key]; !ok {
					t.Errorf(
						"parsePatchArgs(%s) patch body missing key %q",
						tt.args,
						key,
					)
				}
			}
		})
	}
}
