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
	out += formatOptStr("  Output Type", c.OutputType)
	out += formatOptInt("  Pages", c.Pages)
	out += formatOptStr("  Language", c.Language)
	out += formatOptStr("  Mode", c.Mode)
	out += formatOptStr("  Skip Archive File", c.SkipArchiveFile)
	out += formatOptInt("  Image DPI", c.ImageDPI)
	out += formatOptStr("  Unpaper Clean", c.UnpaperClean)
	out += formatOptBool("  Deskew", c.Deskew)
	out += formatOptBool("  Rotate Pages", c.RotatePages)
	out += formatOptFloat("  Rotate Pages Threshold", c.RotatePagesThreshold)
	out += formatOptFloat("  Max Image Pixels", c.MaxImagePixels)
	out += formatOptStr(
		"  Color Conversion Strategy",
		c.ColorConversionStrategy,
	)
	out += formatOptJSON("  User Args", c.UserArgs)

	out += "\nApp Settings:\n"
	out += formatOptStr("  Title", c.AppTitle)
	out += formatOptStr("  Logo", c.AppLogo)

	out += "\nBarcode Settings:\n"
	out += formatOptBool("  Enabled", c.BarcodesEnabled)
	out += formatOptBool(
		"  TIFF Support",
		c.BarcodeEnableTiffSupport,
	)
	out += formatOptStr("  String", c.BarcodeString)
	out += formatOptBool(
		"  Retain Split Pages",
		c.BarcodeRetainSplitPages,
	)
	out += formatOptBool("  Enable ASN", c.BarcodeEnableASN)
	out += formatOptStr("  ASN Prefix", c.BarcodeASNPrefix)
	out += formatOptFloat("  Upscale", c.BarcodeUpscale)
	out += formatOptInt("  DPI", c.BarcodeDPI)
	out += formatOptInt("  Max Pages", c.BarcodeMaxPages)
	out += formatOptBool("  Enable Tag", c.BarcodeEnableTag)
	out += formatOptJSON("  Tag Mapping", c.BarcodeTagMapping)

	return out
}

func formatOptStr(label string, v *string) string {
	if v != nil {
		return fmt.Sprintf("%s: %s\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptInt(label string, v *int64) string {
	if v != nil {
		return fmt.Sprintf("%s: %d\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptBool(label string, v *bool) string {
	if v != nil {
		return fmt.Sprintf("%s: %t\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptFloat(label string, v *float64) string {
	if v != nil {
		return fmt.Sprintf("%s: %g\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptJSON(label string, v json.RawMessage) string {
	if v != nil && string(v) != "null" {
		return fmt.Sprintf("%s: %s\n", label, string(v))
	}
	return fmt.Sprintf("%s: (default)\n", label)
}
