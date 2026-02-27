package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatCustomField(f *models.CustomField) string {
	extraData := "(none)"
	if f.ExtraData != nil && string(f.ExtraData) != "null" {
		extraData = string(f.ExtraData)
	}

	out := fmt.Sprintf("Custom Field (ID: %d)\n", f.ID)
	out += fmt.Sprintf("  Name: %s\n", f.Name)
	out += fmt.Sprintf("  Data Type: %s\n", f.DataType)
	out += fmt.Sprintf("  Extra Data: %s\n", extraData)
	out += fmt.Sprintf("  Document Count: %d\n", f.DocumentCount)

	return out
}

func formatCustomFieldList(
	list *models.PaginatedList[models.CustomField],
) string {
	return formatPaginatedList(
		list,
		"No custom fields found.",
		"Custom Fields",
		func(f models.CustomField) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — type: %s, %d documents\n",
				f.ID,
				f.Name,
				f.ID,
				f.DataType,
				f.DocumentCount,
			)
		},
	)
}
