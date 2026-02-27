package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// UpdateConfig updates the application configuration in Paperless-NGX.
type UpdateConfig struct {
	client *client.Client
}

// NewUpdateConfig creates a new UpdateConfig tool instance.
func NewUpdateConfig(c *client.Client) *UpdateConfig {
	return &UpdateConfig{client: c}
}

// Description returns a description of what this tool does.
func (t *UpdateConfig) Description() string {
	return "Update the application configuration of the Paperless-NGX server"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UpdateConfig) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Config ID to update (typically 1)",
			},
			"output_type": map[string]interface{}{
				"type": "string",
				"description": "PDF output type: pdf, pdfa, " +
					"pdfa-1, pdfa-2, pdfa-3",
			},
			"pages": map[string]interface{}{
				"type":        "integer",
				"description": "Do OCR from page 1 to this value",
			},
			"language": map[string]interface{}{
				"type": "string",
				"description": "OCR language(s), " +
					"e.g. eng, eng+deu",
			},
			"mode": map[string]interface{}{
				"type": "string",
				"description": "OCR mode: skip, redo, " +
					"force, skip_noarchive",
			},
			"skip_archive_file": map[string]interface{}{
				"type": "string",
				"description": "Archive file generation: " +
					"never, with_text, always",
			},
			"image_dpi": map[string]interface{}{
				"type":        "integer",
				"description": "Image DPI fallback value",
			},
			"unpaper_clean": map[string]interface{}{
				"type": "string",
				"description": "Unpaper cleaning: " +
					"clean, clean-final, none",
			},
			"deskew": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable deskew",
			},
			"rotate_pages": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable page rotation",
			},
			"rotate_pages_threshold": map[string]interface{}{
				"type":        "number",
				"description": "Threshold for page rotation",
			},
			"max_image_pixels": map[string]interface{}{
				"type": "number",
				"description": "Maximum image size " +
					"for decompression",
			},
			"color_conversion_strategy": map[string]interface{}{
				"type": "string",
				"description": "Ghostscript color conversion: " +
					"LeaveColorUnchanged, RGB, " +
					"UseDeviceIndependentColor, " +
					"Gray, CMYK",
			},
			"user_args": map[string]interface{}{
				"type": "object",
				"description": "Additional OCRMyPDF " +
					"user arguments (JSON object)",
			},
			"app_title": map[string]interface{}{
				"type":        "string",
				"description": "Application title",
			},
			"barcodes_enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable barcode scanning",
			},
			"barcode_enable_tiff_support": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable barcode TIFF support",
			},
			"barcode_string": map[string]interface{}{
				"type":        "string",
				"description": "Barcode string pattern",
			},
			"barcode_retain_split_pages": map[string]interface{}{
				"type":        "boolean",
				"description": "Retain pages after barcode split",
			},
			"barcode_enable_asn": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable ASN barcode",
			},
			"barcode_asn_prefix": map[string]interface{}{
				"type":        "string",
				"description": "ASN barcode prefix",
			},
			"barcode_upscale": map[string]interface{}{
				"type":        "number",
				"description": "Barcode upscale factor",
			},
			"barcode_dpi": map[string]interface{}{
				"type":        "integer",
				"description": "Barcode DPI",
			},
			"barcode_max_pages": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum pages for barcode scan",
			},
			"barcode_enable_tag": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable tag barcode",
			},
			"barcode_tag_mapping": map[string]interface{}{
				"type": "object",
				"description": "Tag barcode mapping " +
					"(JSON object)",
			},
		},
		"required": []string{"id"},
	}
}

// Execute runs the tool and returns a formatted configuration summary.
func (t *UpdateConfig) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/api/config/%d/", id)

	body, err := doPatchRequest(
		ctx,
		t.client,
		path,
		patchBody,
	)
	if err != nil {
		return "", fmt.Errorf("failed to update config: %w", err)
	}

	var config models.ApplicationConfiguration
	if err := json.Unmarshal(body, &config); err != nil {
		return "", fmt.Errorf(
			"failed to parse config response: %w",
			err,
		)
	}

	return formatConfig(&config), nil
}
