package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatTask(t *models.Task) string {
	out := fmt.Sprintf("Task (ID: %d)\n", t.ID)
	out += fmt.Sprintf("  Task UUID: %s\n", t.TaskID)
	out += fmt.Sprintf("  Status: %s\n", t.Status)
	out += fmt.Sprintf("  Type: %s\n", t.Type)
	out += fmt.Sprintf(
		"  Task Name: %s\n",
		formatOptStr(t.TaskName),
	)
	out += fmt.Sprintf(
		"  File Name: %s\n",
		formatOptStr(t.TaskFileName),
	)
	out += fmt.Sprintf(
		"  Created: %s\n",
		formatOptDate(t.DateCreated),
	)
	out += fmt.Sprintf(
		"  Done: %s\n",
		formatOptDate(t.DateDone),
	)
	out += fmt.Sprintf(
		"  Result: %s\n",
		formatOptStr(t.Result),
	)
	out += fmt.Sprintf(
		"  Acknowledged: %v\n",
		t.Acknowledged,
	)
	out += fmt.Sprintf(
		"  Related Document: %s\n",
		formatOptStr(t.RelatedDocument),
	)

	return out
}

func formatTaskArray(tasks []models.Task) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	out := fmt.Sprintf("Tasks: %d total\n\n", len(tasks))
	for _, task := range tasks {
		name := formatOptStr(task.TaskName)
		out += fmt.Sprintf(
			"%d. %s — %s",
			task.ID,
			name,
			task.Status,
		)
		if task.TaskFileName != nil {
			out += fmt.Sprintf(
				" — %s",
				*task.TaskFileName,
			)
		}
		out += "\n"
	}

	return out
}
