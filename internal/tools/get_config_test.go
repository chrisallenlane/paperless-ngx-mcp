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

const configResponseAllDefaults = `[{
	"id": 1,
	"user_args": null,
	"barcode_tag_mapping": null,
	"output_type": null,
	"pages": null,
	"language": null,
	"mode": null,
	"skip_archive_file": null,
	"image_dpi": null,
	"unpaper_clean": null,
	"deskew": null,
	"rotate_pages": null,
	"rotate_pages_threshold": null,
	"max_image_pixels": null,
	"color_conversion_strategy": null,
	"app_title": null,
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
}]`

const configResponseWithValues = `[{
	"id": 1,
	"user_args": {"--deskew": true},
	"barcode_tag_mapping": {"ASN": "tag1"},
	"output_type": "pdfa",
	"pages": 5,
	"language": "eng+deu",
	"mode": "skip",
	"skip_archive_file": "with_text",
	"image_dpi": 300,
	"unpaper_clean": "clean",
	"deskew": true,
	"rotate_pages": false,
	"rotate_pages_threshold": 12.5,
	"max_image_pixels": 500000000.0,
	"color_conversion_strategy": "RGB",
	"app_title": "My Paperless",
	"app_logo": "/media/logo/custom.png",
	"barcodes_enabled": true,
	"barcode_enable_tiff_support": false,
	"barcode_string": "PATCHT",
	"barcode_retain_split_pages": true,
	"barcode_enable_asn": true,
	"barcode_asn_prefix": "ASN",
	"barcode_upscale": 1.5,
	"barcode_dpi": 200,
	"barcode_max_pages": 10,
	"barcode_enable_tag": false
}]`

func TestGetConfig_Execute_AllDefaults(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/config/" {
					t.Errorf(
						"Expected /api/config/, got %s",
						r.URL.Path,
					)
				}
				if r.Method != "GET" {
					t.Errorf("Expected GET, got %s", r.Method)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(configResponseAllDefaults))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetConfig(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Paperless-NGX Configuration (ID: 1)",
		"OCR Settings:",
		"  Output Type: (default)",
		"  Pages: (default)",
		"  Language: (default)",
		"  Mode: (default)",
		"App Settings:",
		"  Title: (default)",
		"  Logo: (default)",
		"Barcode Settings:",
		"  Enabled: (default)",
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

func TestGetConfig_Execute_WithValues(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(configResponseWithValues))
			},
		),
	)
	defer server.Close()

	c := client.NewWithHTTPClient(
		server.URL,
		"test-token",
		server.Client(),
	)
	tool := NewGetConfig(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checks := []string{
		"Output Type: pdfa",
		"Pages: 5",
		"Language: eng+deu",
		"Mode: skip",
		"Skip Archive File: with_text",
		"Image DPI: 300",
		"Unpaper Clean: clean",
		"Deskew: true",
		"Rotate Pages: false",
		"Rotate Pages Threshold: 12.5",
		"Max Image Pixels: 5e+08",
		"Color Conversion Strategy: RGB",
		"Title: My Paperless",
		"Logo: /media/logo/custom.png",
		"Enabled: true",
		"TIFF Support: false",
		"String: PATCHT",
		"Retain Split Pages: true",
		"Enable ASN: true",
		"ASN Prefix: ASN",
		"Upscale: 1.5",
		"DPI: 200",
		"Max Pages: 10",
		"Enable Tag: false",
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

func TestGetConfig_Execute_EmptyArray(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
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
	tool := NewGetConfig(c)

	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "No configuration found." {
		t.Errorf("Expected empty result message, got: %s", result)
	}
}

func TestGetConfig_Execute_ServerError(t *testing.T) {
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
	tool := NewGetConfig(c)

	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
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

func TestGetConfig_Description(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetConfig(c)

	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestGetConfig_InputSchema(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	tool := NewGetConfig(c)

	schema := tool.InputSchema()
	if schema == nil {
		t.Fatal("InputSchema should not be nil")
	}

	schemaType, ok := schema["type"].(string)
	if !ok || schemaType != "object" {
		t.Errorf("Schema type = %v, want object", schema["type"])
	}
}
