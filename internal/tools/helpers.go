// Package tools provides MCP tool implementations.
package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"

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
