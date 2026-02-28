package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListDocumentNotes lists notes for a document.
type ListDocumentNotes struct {
	client *client.Client
}

// NewListDocumentNotes creates a new ListDocumentNotes tool
// instance.
func NewListDocumentNotes(
	c *client.Client,
) *ListDocumentNotes {
	return &ListDocumentNotes{client: c}
}

// Description returns a description of what this tool does.
func (t *ListDocumentNotes) Description() string {
	return "List notes attached to a document " +
		"in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *ListDocumentNotes) InputSchema() map[string]interface{} {
	return idOnlySchema("Document ID")
}

// Execute runs the tool and returns a formatted notes list.
func (t *ListDocumentNotes) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	id, err := parseIDArg(args)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf(
		"/api/documents/%d/notes/",
		id,
	)

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list document notes: %w",
			err,
		)
	}

	var notes []models.Note
	if err := json.Unmarshal(body, &notes); err != nil {
		return "", fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return formatNoteList(id, notes), nil
}
