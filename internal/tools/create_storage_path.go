package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateStoragePath creates a new storage path in Paperless-NGX.
type CreateStoragePath struct {
	client *client.Client
}

// NewCreateStoragePath creates a new CreateStoragePath tool
// instance.
func NewCreateStoragePath(
	c *client.Client,
) *CreateStoragePath {
	return &CreateStoragePath{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateStoragePath) Description() string {
	return "Create a new storage path in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *CreateStoragePath) InputSchema() map[string]interface{} {
	return storagePathSchema(false)
}

// Execute runs the tool and returns a formatted storage path
// summary.
func (t *CreateStoragePath) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Name              string `json:"name"`
		Path              string `json:"path"`
		Match             string `json:"match,omitempty"`
		MatchingAlgorithm *int   `json:"matching_algorithm,omitempty"`
		IsInsensitive     *bool  `json:"is_insensitive,omitempty"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if params.Path == "" {
		return "", fmt.Errorf("path is required")
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		"/api/storage_paths/",
		params,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create storage path: %w",
			err,
		)
	}

	var sp models.StoragePath
	if err := json.Unmarshal(body, &sp); err != nil {
		return "", fmt.Errorf(
			"failed to parse storage path response: %w",
			err,
		)
	}

	return formatStoragePath(&sp), nil
}
