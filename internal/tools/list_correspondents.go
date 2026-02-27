package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListCorrespondents lists correspondents from Paperless-NGX.
type ListCorrespondents struct {
	client *client.Client
}

// NewListCorrespondents creates a new ListCorrespondents tool instance.
func NewListCorrespondents(c *client.Client) *ListCorrespondents {
	return &ListCorrespondents{client: c}
}

// Description returns a description of what this tool does.
func (t *ListCorrespondents) Description() string {
	return "List correspondents in Paperless-NGX " +
		"with optional filtering by name"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *ListCorrespondents) InputSchema() map[string]interface{} {
	return paginatedListSchema()
}

// Execute runs the tool and returns a formatted correspondent list.
func (t *ListCorrespondents) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	list, err := listResources[models.Correspondent](
		ctx,
		t.client,
		"/api/correspondents/",
		args,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list correspondents: %w",
			err,
		)
	}

	return formatCorrespondentList(list), nil
}
