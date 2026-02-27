package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

var ruleTypeNames = map[int]string{
	0:  "Title contains",
	1:  "Content contains",
	2:  "ASN is",
	3:  "Correspondent is",
	4:  "Document type is",
	5:  "Is in inbox",
	6:  "Has tag",
	7:  "Has any tag",
	8:  "Created before",
	9:  "Created after",
	10: "Created year",
	11: "Created month",
	12: "Created day",
	13: "Added before",
	14: "Added after",
	15: "Modified before",
	16: "Modified after",
	17: "Does not have tag",
	18: "Does not have ASN",
	19: "Title or content contains",
	20: "Fulltext query",
	21: "More like document",
	22: "Has tags in",
	23: "ASN greater than",
	24: "ASN less than",
	25: "Storage path is",
	26: "Has correspondent in",
	27: "Does not have correspondent in",
	28: "Has document type in",
	29: "Does not have document type in",
	30: "Has storage path in",
	31: "Does not have storage path in",
	32: "Has tags all",
	33: "Owner is",
	34: "Owner is in",
	35: "Does not have owner in",
	36: "Correspondent starts with",
	37: "Correspondent ends with",
	38: "Title starts with",
	39: "Title ends with",
	44: "Has custom field value",
	45: "Custom field query",
	46: "Is shared by me",
	47: "Has custom fields in",
}

// ruleTypeName returns a human-readable name for a filter
// rule type.
func ruleTypeName(ruleType int) string {
	if name, ok := ruleTypeNames[ruleType]; ok {
		return name
	}
	return fmt.Sprintf("Rule type %d", ruleType)
}

func formatSavedView(v *models.SavedView) string {
	out := fmt.Sprintf("Saved View (ID: %d)\n", v.ID)
	out += fmt.Sprintf("  Name: %s\n", v.Name)
	out += fmt.Sprintf(
		"  Show on Dashboard: %v\n",
		v.ShowOnDashboard,
	)
	out += fmt.Sprintf(
		"  Show in Sidebar: %v\n",
		v.ShowInSidebar,
	)
	out += fmt.Sprintf(
		"  Sort Field: %s\n",
		formatOptStr(v.SortField),
	)
	out += fmt.Sprintf(
		"  Sort Reverse: %v\n",
		v.SortReverse,
	)
	out += fmt.Sprintf(
		"  Page Size: %s\n",
		formatOptInt(v.PageSize),
	)
	out += fmt.Sprintf(
		"  Display Mode: %s\n",
		formatOptStr(v.DisplayMode),
	)

	if len(v.FilterRules) == 0 {
		out += "  Filter Rules: (none)\n"
	} else {
		out += "  Filter Rules:\n"
		for _, r := range v.FilterRules {
			val := "(null)"
			if r.Value != nil {
				val = *r.Value
			}
			out += fmt.Sprintf(
				"    - %s: %s\n",
				ruleTypeName(r.RuleType),
				val,
			)
		}
	}

	return out
}

func formatSavedViewList(
	list *models.PaginatedList[models.SavedView],
) string {
	return formatPaginatedList(
		list,
		"No saved views found.",
		"Saved Views",
		func(v models.SavedView) string {
			flags := ""
			if v.ShowOnDashboard {
				flags += " [dashboard]"
			}
			if v.ShowInSidebar {
				flags += " [sidebar]"
			}
			return fmt.Sprintf(
				"%d. %s (ID: %d) — "+
					"%d filter rules%s\n",
				v.ID,
				v.Name,
				v.ID,
				len(v.FilterRules),
				flags,
			)
		},
	)
}
