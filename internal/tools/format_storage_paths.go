package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatStoragePath(sp *models.StoragePath) string {
	out := formatMatchableFields(
		"Storage Path",
		sp.ID,
		sp.Name,
		sp.Slug,
		sp.Match,
		sp.MatchingAlgorithm,
		sp.IsInsensitive,
		sp.DocumentCount,
	)
	out += fmt.Sprintf("  Path: %s\n", sp.Path)

	return out
}

func formatStoragePathList(
	list *models.PaginatedList[models.StoragePath],
) string {
	return formatPaginatedList(
		list,
		"No storage paths found.",
		"Storage Paths",
		func(sp models.StoragePath) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				sp.ID,
				sp.Name,
				sp.ID,
				sp.DocumentCount,
			)
		},
	)
}
