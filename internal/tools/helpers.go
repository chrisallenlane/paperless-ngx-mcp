// Package tools provides MCP tool implementations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// doAPIRequest performs a GET API request and returns the response body.
// It handles common patterns: making the request, checking status, reading body.
// Includes response body in error messages when status is not OK.
func doAPIRequest(
	ctx context.Context,
	c *client.Client,
	path string,
) ([]byte, error) {
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, http.StatusOK)
}

// doPatchRequest performs a PATCH API request and returns the response body.
func doPatchRequest(
	ctx context.Context,
	c *client.Client,
	path string,
	body interface{},
) ([]byte, error) {
	resp, err := c.Patch(ctx, path, body)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, http.StatusOK)
}

// doPostRequest performs a POST API request and returns the response body.
func doPostRequest(
	ctx context.Context,
	c *client.Client,
	path string,
	body interface{},
) ([]byte, error) {
	resp, err := c.Post(ctx, path, body)
	if err != nil {
		return nil, err
	}
	return readResponse(resp, http.StatusCreated)
}

// doDeleteRequest performs a DELETE API request.
func doDeleteRequest(
	ctx context.Context,
	c *client.Client,
	path string,
) error {
	resp, err := c.Delete(ctx, path)
	if err != nil {
		return err
	}
	_, err = readResponse(resp, http.StatusNoContent)
	return err
}

// parseIDArg extracts and validates a positive integer "id" from JSON args.
func parseIDArg(args json.RawMessage) (int, error) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return 0, fmt.Errorf("failed to parse arguments: %w", err)
	}
	if params.ID <= 0 {
		return 0, fmt.Errorf("id must be a positive integer")
	}
	return params.ID, nil
}

// parsePatchArgs extracts a positive integer "id" and builds a patch body
// from the remaining fields in the JSON args.
func parsePatchArgs(
	args json.RawMessage,
) (int, map[string]json.RawMessage, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(args, &raw); err != nil {
		return 0, nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	idRaw, ok := raw["id"]
	if !ok {
		return 0, nil, fmt.Errorf("id is required")
	}

	var id int
	if err := json.Unmarshal(idRaw, &id); err != nil {
		return 0, nil, fmt.Errorf("failed to parse id: %w", err)
	}

	if id <= 0 {
		return 0, nil, fmt.Errorf("id must be a positive integer")
	}

	patchBody := make(map[string]json.RawMessage)
	for k, v := range raw {
		if k != "id" {
			patchBody[k] = v
		}
	}

	return id, patchBody, nil
}

// emptySchema returns an input schema with no parameters.
func emptySchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// idOnlySchema returns an input schema with a single required "id" field.
func idOnlySchema(desc string) map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": desc,
			},
		},
		"required": []string{"id"},
	}
}

// paginatedListSchema returns an input schema for paginated list endpoints.
func paginatedListSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (default 1)",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Results per page (default 25)",
			},
			"name": map[string]interface{}{
				"type": "string",
				"description": "Filter by name " +
					"(case-insensitive contains)",
			},
		},
	}
}

// matchableResourceSchema returns an input schema for resources with matching
// fields (name, match, matching_algorithm, is_insensitive). Set includeID
// to true for update tools, false for create tools.
func matchableResourceSchema(
	resourceName string,
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": resourceName + " name",
		},
		"match": map[string]interface{}{
			"type":        "string",
			"description": "Match pattern for auto-assignment",
		},
		"matching_algorithm": map[string]interface{}{
			"type": "integer",
			"description": "Matching algorithm: " +
				"0=None, 1=Any word, 2=All words, " +
				"3=Exact match, 4=Regex, " +
				"5=Fuzzy word, 6=Automatic",
		},
		"is_insensitive": map[string]interface{}{
			"type":        "boolean",
			"description": "Case-insensitive matching",
		},
	}

	required := []string{"name"}
	if includeID {
		props["id"] = map[string]interface{}{
			"type":        "integer",
			"description": resourceName + " ID to update",
		}
		required = []string{"id"}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// customFieldSchema returns an input schema for custom field tools.
// Set includeID to true for update tools, false for create tools.
func customFieldSchema(
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Custom field name",
		},
		"data_type": map[string]interface{}{
			"type": "string",
			"description": "Data type: string, url, " +
				"date, boolean, integer, float, " +
				"monetary, documentlink, " +
				"select, longtext",
		},
		"extra_data": map[string]interface{}{
			"type": "object",
			"description": "Additional field " +
				"configuration (JSON object)",
		},
	}

	required := []string{"name", "data_type"}
	if includeID {
		props["id"] = map[string]interface{}{
			"type":        "integer",
			"description": "Custom field ID to update",
		}
		required = []string{"id"}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

// documentUpdateSchema returns an input schema for the document update tool.
func documentUpdateSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID to update",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Document title (max 128 chars)",
				"maxLength":   128,
			},
			"correspondent": map[string]interface{}{
				"type": "integer",
				"description": "Correspondent ID " +
					"(null to clear)",
			},
			"document_type": map[string]interface{}{
				"type": "integer",
				"description": "Document type ID " +
					"(null to clear)",
			},
			"storage_path": map[string]interface{}{
				"type": "integer",
				"description": "Storage path ID " +
					"(null to clear)",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer",
				},
				"description": "Replace all tags " +
					"with these tag IDs",
			},
			"archive_serial_number": map[string]interface{}{
				"type":        "integer",
				"description": "Archive serial number (null to clear)",
				"minimum":     0,
			},
			"created": map[string]interface{}{
				"type":        "string",
				"description": "Creation date (ISO format)",
			},
			"custom_fields": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field": map[string]interface{}{
							"type":        "integer",
							"description": "Custom field ID",
						},
						"value": map[string]interface{}{
							"description": "Field value " +
								"(type varies by field)",
						},
					},
					"required": []string{
						"field",
						"value",
					},
				},
				"description": "Custom field values to set",
			},
		},
		"required": []string{"id"},
	}
}

// configUpdateSchema returns an input schema for the config update tool.
func configUpdateSchema() map[string]interface{} {
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

// deleteByID parses an ID from args and performs a DELETE request.
func deleteByID(
	ctx context.Context,
	c *client.Client,
	args json.RawMessage,
	pathFmt string,
	resourceName string,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf(pathFmt, id)

	if err := doDeleteRequest(ctx, c, path); err != nil {
		return "", fmt.Errorf(
			"failed to delete %s: %w",
			resourceName,
			err,
		)
	}

	return fmt.Sprintf(
		"%s %d deleted successfully.",
		resourceName,
		id,
	), nil
}

// fetchByID parses an ID, fetches a resource, and unmarshals the response.
func fetchByID[T any](
	ctx context.Context,
	c *client.Client,
	args json.RawMessage,
	pathFmt string,
) (*T, int, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return nil, 0, err
	}

	path := fmt.Sprintf(pathFmt, id)

	body, err := doAPIRequest(ctx, c, path)
	if err != nil {
		return nil, 0, err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return &result, id, nil
}

// patchByID parses patch args, performs a PATCH, and unmarshals the response.
func patchByID[T any](
	ctx context.Context,
	c *client.Client,
	args json.RawMessage,
	pathFmt string,
) (*T, error) {
	id, patchBody, err := parsePatchArgs(args)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(pathFmt, id)

	body, err := doPatchRequest(ctx, c, path, patchBody)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return &result, nil
}

// listResources builds a list path, fetches, and unmarshals a paginated list.
func listResources[T any](
	ctx context.Context,
	c *client.Client,
	basePath string,
	args json.RawMessage,
) (*models.PaginatedList[T], error) {
	path, err := buildListPath(basePath, args)
	if err != nil {
		return nil, err
	}

	body, err := doAPIRequest(ctx, c, path)
	if err != nil {
		return nil, err
	}

	var list models.PaginatedList[T]
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return &list, nil
}

// createMatchable creates a matchable resource (correspondent, document type).
func createMatchable[T any](
	ctx context.Context,
	c *client.Client,
	args json.RawMessage,
	path string,
) (*T, error) {
	var params matchableCreateParams
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	body, err := doPostRequest(ctx, c, path, params)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return &result, nil
}

// matchableCreateParams holds common parameters for creating matchable
// resources (correspondents, document types).
type matchableCreateParams struct {
	Name              string `json:"name"`
	Match             string `json:"match,omitempty"`
	MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
	IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
}

// listParams holds common pagination and filter parameters.
type listParams struct {
	Page     *int    `json:"page"`
	PageSize *int    `json:"page_size"`
	Name     *string `json:"name"`
}

// buildListPath constructs a paginated API path with query parameters.
func buildListPath(
	basePath string,
	args json.RawMessage,
) (string, error) {
	var params listParams
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	q := url.Values{}
	if params.Page != nil {
		q.Set("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.PageSize != nil {
		q.Set(
			"page_size",
			fmt.Sprintf("%d", *params.PageSize),
		)
	}
	if params.Name != nil {
		q.Set("name__icontains", *params.Name)
	}

	if encoded := q.Encode(); encoded != "" {
		return basePath + "?" + encoded, nil
	}
	return basePath, nil
}

// validateFilePath checks that a file path is safe and absolute.
func validateFilePath(path string) error {
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return fmt.Errorf(
			"file path must be absolute: %s",
			path,
		)
	}

	if strings.Contains(cleaned, "..") {
		return fmt.Errorf(
			"file path must not contain '..': %s",
			path,
		)
	}

	return nil
}

// readResponse reads and validates an HTTP response, returning the body bytes.
func readResponse(
	resp *http.Response,
	expectedStatus int,
) ([]byte, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf(
			"unexpected status code %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil
}
