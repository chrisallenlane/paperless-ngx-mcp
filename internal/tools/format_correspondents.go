package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatCorrespondent(c *models.Correspondent) string {
	out := formatMatchableFields(
		"Correspondent",
		c.ID,
		c.Name,
		c.Slug,
		c.Match,
		c.MatchingAlgorithm,
		c.IsInsensitive,
		c.DocumentCount,
	)

	lastCorr := "(none)"
	if c.LastCorrespondence != nil {
		lastCorr = formatDate(*c.LastCorrespondence)
	}
	out += fmt.Sprintf("  Last Correspondence: %s\n", lastCorr)

	return out
}

func formatCorrespondentList(
	list *models.PaginatedList[models.Correspondent],
) string {
	return formatPaginatedList(
		list,
		"No correspondents found.",
		"Correspondents",
		func(c models.Correspondent) string {
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents\n",
				c.ID,
				c.Name,
				c.ID,
				c.DocumentCount,
			)
		},
	)
}
