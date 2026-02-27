// Package client provides an HTTP client for the Paperless-NGX API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPDoer interface allows mocking HTTP requests for testing
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a Paperless-NGX HTTP API client.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient HTTPDoer
}

// New creates a new HTTP client with default settings
func New(baseURL, token string) *Client {
	return NewWithHTTPClient(baseURL, token, &http.Client{
		Timeout: 30 * time.Second,
	})
}

// NewWithHTTPClient creates a new client with a custom HTTP client
// (useful for testing).
func NewWithHTTPClient(
	baseURL, token string,
	httpClient HTTPDoer,
) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: httpClient,
	}
}

// doRequest performs an HTTP request
func (c *Client) doRequest(
	ctx context.Context,
	method, path string,
	body []byte,
) (*http.Response, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Token "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(
	ctx context.Context,
	path string,
) (*http.Response, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(
	ctx context.Context,
	path string,
	body interface{},
) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return c.doRequest(ctx, "POST", path, data)
}

// Patch performs a PATCH request
func (c *Client) Patch(
	ctx context.Context,
	path string,
	body interface{},
) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return c.doRequest(ctx, "PATCH", path, data)
}

// Delete performs a DELETE request
func (c *Client) Delete(
	ctx context.Context,
	path string,
) (*http.Response, error) {
	return c.doRequest(ctx, "DELETE", path, nil)
}

// PostMultipart performs a POST request with a multipart/form-data body.
// The body reader should contain the pre-built multipart content, and
// contentType should include the multipart boundary.
func (c *Client) PostMultipart(
	ctx context.Context,
	path string,
	body io.Reader,
	contentType string,
) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", contentType)
	if c.Token != "" {
		req.Header.Set(
			"Authorization",
			"Token "+c.Token,
		)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
