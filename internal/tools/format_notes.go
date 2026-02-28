package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatNoteList(
	docID int,
	notes []models.Note,
) string {
	if len(notes) == 0 {
		return fmt.Sprintf(
			"No notes found for document %d.",
			docID,
		)
	}

	out := fmt.Sprintf(
		"Notes for Document %d: %d total\n\n",
		docID,
		len(notes),
	)

	for _, note := range notes {
		out += fmt.Sprintf(
			"--- Note %d (by %s on %s) ---\n%s\n\n",
			note.ID,
			note.User.Username,
			formatDate(note.Created),
			note.Note,
		)
	}

	return out
}
