package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatNote(note *models.Note) string {
	out := fmt.Sprintf("Note (ID: %d)\n", note.ID)
	out += fmt.Sprintf(
		"  Author: %s\n",
		note.User.Username,
	)
	out += fmt.Sprintf(
		"  Created: %s\n",
		formatDate(note.Created),
	)
	out += fmt.Sprintf("  Content: %s\n", note.Note)

	return out
}

func formatNoteList(
	docID int,
	list *models.PaginatedList[models.Note],
) string {
	if list.Count == 0 {
		return fmt.Sprintf(
			"No notes found for document %d.",
			docID,
		)
	}

	out := fmt.Sprintf(
		"Notes for Document %d: %d total\n\n",
		docID,
		list.Count,
	)

	for _, note := range list.Results {
		out += fmt.Sprintf(
			"--- Note %d (by %s on %s) ---\n%s\n\n",
			note.ID,
			note.User.Username,
			formatDate(note.Created),
			note.Note,
		)
	}

	if list.Next != nil {
		out += "(More notes available — " +
			"use page parameter)\n"
	}

	return out
}
