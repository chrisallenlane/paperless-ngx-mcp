package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatTag(tag *models.Tag) string {
	out := formatMatchableFields(
		"Tag",
		tag.ID,
		tag.Name,
		tag.Slug,
		tag.Match,
		tag.MatchingAlgorithm,
		tag.IsInsensitive,
		tag.DocumentCount,
	)
	out += fmt.Sprintf("  Color: %s\n", tag.Color)
	out += fmt.Sprintf(
		"  Text Color: %s\n",
		tag.TextColor,
	)
	out += fmt.Sprintf(
		"  Is Inbox Tag: %v\n",
		tag.IsInboxTag,
	)
	out += fmt.Sprintf(
		"  Parent: %s\n",
		formatOptInt(tag.Parent),
	)
	out += fmt.Sprintf(
		"  Children: %s\n",
		formatIntSlice(tag.Children),
	)

	return out
}

func formatTagList(
	list *models.PaginatedList[models.Tag],
) string {
	return formatPaginatedList(
		list,
		"No tags found.",
		"Tags",
		func(tag models.Tag) string {
			extra := ""
			if tag.IsInboxTag {
				extra = " [inbox]"
			}
			return fmt.Sprintf(
				"%d. %s (ID: %d) — %d documents%s\n",
				tag.ID,
				tag.Name,
				tag.ID,
				tag.DocumentCount,
				extra,
			)
		},
	)
}
