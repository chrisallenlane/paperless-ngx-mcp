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

const updateConfigResponse = `{
	"id": 1,
	"user_args": null,
	"barcode_tag_mapping": null,
	"output_type": "pdfa",
	"pages": null,
	"language": "eng+deu",
	"mode": null,
	"skip_archive_file": null,
	"image_dpi": null,
	"unpaper_clean": null,
	"deskew": true,
	"rotate_pages": null,
	"rotate_pages_threshold": null,
	"max_image_pixels": null,
	"color_conversion_strategy": null,
	"app_title": "My Docs",
	"app_logo": null,
	"barcodes_enabled": null,
	"barcode_enable_tiff_support": null,
	"barcode_string": null,
	"barcode_retain_split_pages": null,
	"barcode_enable_asn": null,
	"barcode_asn_prefix": null,
	"barcode_upscale": null,
	"barcode_dpi": null,
	"barcode_max_pages": null,
	"barcode_enable_tag": null
}`

func TestUpdateConfig_Execute(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/config/1/" {
					t.Errorf(
						"Expected /api/config/1/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "PATCH" {
					t.Errorf("Expected PATCH, got %s", r.Method)
				}

				// Verify request body contains only
				// the fields we sent
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %v", err)
				}

				var patch map[string]interface{}
				if err := json.Unmarshal(
					body,
					&patch,
				); err != nil {
					t.Fatalf(
						"Failed to parse body: %v",
						err,
					)
				}

				// Should have 3 fields, not id
				if _, ok := patch["id"]; ok {
					t.Error("Body should not contain id")
				}

				if patch["output_type"] != "pdfa" {
					t.Errorf(
						"output_type = %v, want pdfa",
						patch["output_type"],
					)
				}

				if patch["deskew"] != true {
					t.Errorf(
						"deskew = %v, want true",
						patch["deskew"],
					)
				}

				if patch["app_title"] != "My Docs" {
					t.Errorf(
						"app_title = %v, want My Docs",
						patch["app_title"],
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(updateConfigResponse))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewUpdateConfig(c)

	args := `{
		"id": 1,
		"output_type": "pdfa",
		"deskew": true,
		"app_title": "My Docs"
	}`

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Paperless-NGX Configuration (ID: 1)",
		"Output Type: pdfa",
		"Deskew: true",
		"Title: My Docs",
		"Language: eng+deu",
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
