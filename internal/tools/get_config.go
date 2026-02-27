package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetConfig retrieves the application configuration from Paperless-NGX.
type GetConfig struct {
	client *client.Client
}

// NewGetConfig creates a new GetConfig tool instance.
func NewGetConfig(c *client.Client) *GetConfig {
	return &GetConfig{client: c}
}

// Description returns a description of what this tool does.
func (t *GetConfig) Description() string {
	return "Get the current application configuration of the Paperless-NGX server"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetConfig) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// Execute runs the tool and returns a formatted configuration summary.
func (t *GetConfig) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, "/api/config/")
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	var configs []models.ApplicationConfiguration
	if err := json.Unmarshal(body, &configs); err != nil {
		return "", fmt.Errorf(
			"failed to parse config response: %w",
			err,
		)
	}

	if len(configs) == 0 {
		return "No configuration found.", nil
	}

	return formatConfig(&configs[0]), nil
}

func formatConfig(c *models.ApplicationConfiguration) string {
	out := fmt.Sprintf(
		"Paperless-NGX Configuration (ID: %d)\n",
		c.ID,
	)

	out += "\nOCR Settings:\n"
	out += formatOpt("  Output Type", c.OutputType)
	out += formatOpt("  Pages", c.Pages)
	out += formatOpt("  Language", c.Language)
	out += formatOpt("  Mode", c.Mode)
	out += formatOpt("  Skip Archive File", c.SkipArchiveFile)
	out += formatOpt("  Image DPI", c.ImageDPI)
	out += formatOpt("  Unpaper Clean", c.UnpaperClean)
	out += formatOpt("  Deskew", c.Deskew)
	out += formatOpt("  Rotate Pages", c.RotatePages)
	out += formatOpt("  Rotate Pages Threshold", c.RotatePagesThreshold)
	out += formatOpt("  Max Image Pixels", c.MaxImagePixels)
	out += formatOpt(
		"  Color Conversion Strategy",
		c.ColorConversionStrategy,
	)
	out += formatOptJSON("  User Args", c.UserArgs)

	out += "\nApp Settings:\n"
	out += formatOpt("  Title", c.AppTitle)
	out += formatOpt("  Logo", c.AppLogo)

	out += "\nBarcode Settings:\n"
	out += formatOpt("  Enabled", c.BarcodesEnabled)
	out += formatOpt(
		"  TIFF Support",
		c.BarcodeEnableTiffSupport,
	)
	out += formatOpt("  String", c.BarcodeString)
	out += formatOpt(
		"  Retain Split Pages",
		c.BarcodeRetainSplitPages,
	)
	out += formatOpt("  Enable ASN", c.BarcodeEnableASN)
	out += formatOpt("  ASN Prefix", c.BarcodeASNPrefix)
	out += formatOpt("  Upscale", c.BarcodeUpscale)
	out += formatOpt("  DPI", c.BarcodeDPI)
	out += formatOpt("  Max Pages", c.BarcodeMaxPages)
	out += formatOpt("  Enable Tag", c.BarcodeEnableTag)
	out += formatOptJSON("  Tag Mapping", c.BarcodeTagMapping)

	return out
}

func formatOpt[T any](label string, v *T) string {
	if v != nil {
		return fmt.Sprintf("%s: %v\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptJSON(label string, v json.RawMessage) string {
	if v != nil && string(v) != "null" {
		return fmt.Sprintf("%s: %s\n", label, string(v))
	}
	return fmt.Sprintf("%s: (default)\n", label)
}
