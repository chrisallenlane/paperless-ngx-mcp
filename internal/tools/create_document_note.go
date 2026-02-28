package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// CreateDocumentNote adds a note to a document.
type CreateDocumentNote struct {
	client *client.Client
}

// NewCreateDocumentNote creates a new CreateDocumentNote tool
// instance.
func NewCreateDocumentNote(
	c *client.Client,
) *CreateDocumentNote {
	return &CreateDocumentNote{client: c}
}

// Description returns a description of what this tool does.
func (t *CreateDocumentNote) Description() string {
	return "Add a note to a document in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *CreateDocumentNote) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID",
			},
			"note": map[string]interface{}{
				"type":        "string",
				"description": "Note text to add",
			},
		},
		"required": []string{"id", "note"},
	}
}

// Execute runs the tool and returns the updated notes list.
func (t *CreateDocumentNote) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		ID   int    `json:"id"`
		Note string `json:"note"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.ID <= 0 {
		return "", fmt.Errorf(
			"id must be a positive integer",
		)
	}

	if params.Note == "" {
		return "", fmt.Errorf("note is required")
	}

	path := fmt.Sprintf(
		"/api/documents/%d/notes/",
		params.ID,
	)

	// POST to notes endpoint returns 200 with updated
	// notes list (not 201).
	resp, err := t.client.Post(
		ctx,
		path,
		map[string]string{"note": params.Note},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create note: %w",
			err,
		)
	}

	body, err := readResponse(resp, http.StatusOK)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create note: %w",
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

	return fmt.Sprintf(
		"Note added to document %d.\n\n%s",
		params.ID,
		formatNoteList(params.ID, notes),
	), nil
}
