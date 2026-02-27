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

	return appendQuery(basePath, q), nil
}

// appendQuery appends encoded query parameters to a base path.
// Returns the base path unchanged if q is empty.
func appendQuery(
	basePath string,
	q url.Values,
) string {
	if encoded := q.Encode(); encoded != "" {
		return basePath + "?" + encoded
	}
	return basePath
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
