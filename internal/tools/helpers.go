// Package tools provides MCP tool implementations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
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
