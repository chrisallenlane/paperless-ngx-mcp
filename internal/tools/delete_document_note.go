package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// DeleteDocumentNote deletes a note from a document.
type DeleteDocumentNote struct {
	client *client.Client
}

// NewDeleteDocumentNote creates a new DeleteDocumentNote tool
// instance.
func NewDeleteDocumentNote(
	c *client.Client,
) *DeleteDocumentNote {
	return &DeleteDocumentNote{client: c}
}

// Description returns a description of what this tool does.
func (t *DeleteDocumentNote) Description() string {
	return "Delete a note from a document " +
		"in Paperless-NGX"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *DeleteDocumentNote) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"document_id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID",
			},
			"note_id": map[string]interface{}{
				"type":        "integer",
				"description": "Note ID to delete",
			},
		},
		"required": []string{
			"document_id",
			"note_id",
		},
	}
}

// Execute runs the tool and returns the updated notes list.
func (t *DeleteDocumentNote) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		DocumentID int `json:"document_id"`
		NoteID     int `json:"note_id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.DocumentID <= 0 {
		return "", fmt.Errorf(
			"document_id must be a positive integer",
		)
	}

	if params.NoteID <= 0 {
		return "", fmt.Errorf(
			"note_id must be a positive integer",
		)
	}

	// DELETE uses query parameter for note ID.
	path := fmt.Sprintf(
		"/api/documents/%d/notes/?id=%d",
		params.DocumentID,
		params.NoteID,
	)

	// DELETE to notes endpoint returns 200 with updated
	// notes list (not 204).
	resp, err := t.client.Delete(ctx, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete note: %w",
			err,
		)
	}

	body, err := readResponse(resp, http.StatusOK)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete note: %w",
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
		"Note %d deleted from document %d.\n\n%s",
		params.NoteID,
		params.DocumentID,
		formatNoteList(params.DocumentID, notes),
	), nil
}
