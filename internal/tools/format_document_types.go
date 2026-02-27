package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatDocumentType(dt *models.DocumentType) string {
	return formatMatchableFields(
		"Document Type",
		dt.ID,
		dt.Name,
		dt.Slug,
		dt.Match,
		dt.MatchingAlgorithm,
		dt.IsInsensitive,
		dt.DocumentCount,
	)
}

func formatDocumentTypeList(
	list *models.PaginatedList[models.DocumentType],
) string {
	return formatPaginatedList(
		list,
		"No document types found.",
		"Document Types",
		func(dt models.DocumentType) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				dt.ID,
				dt.Name,
				dt.ID,
				dt.DocumentCount,
			)
		},
	)
}
