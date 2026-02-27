package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

const paginationHint = "\n(more results available — use page parameter)"

var matchingAlgorithmNames = map[int]string{
	0: "None",
	1: "Any word",
	2: "All words",
	3: "Exact match",
	4: "Regex",
	5: "Fuzzy word",
	6: "Automatic",
}

func matchingAlgorithmName(algo int) string {
	name := matchingAlgorithmNames[algo]
	if name == "" {
		return "Unknown"
	}
	return name
}

func formatOpt[T any](label string, v *T) string {
	if v != nil {
		return fmt.Sprintf("%s: %v\n", label, *v)
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptJSON(label string, v json.RawMessage) string {
	if v != nil && string(v) != "null" {
		return fmt.Sprintf("%s: %s\n", label, string(v))
	}
	return fmt.Sprintf("%s: (default)\n", label)
}

func formatOptInt(v *int) string {
	if v != nil {
		return fmt.Sprintf("%d", *v)
	}
	return "(none)"
}

func formatOptStr(v *string) string {
	if v != nil && *v != "" {
		return *v
	}
	return "(none)"
}

func formatOptDate(v *string) string {
	if v != nil && *v != "" {
		return formatDate(*v)
	}
	return "(none)"
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func formatFileSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
	)

	switch {
	case bytes >= tb:
		return fmt.Sprintf(
			"%.2f TB",
			float64(bytes)/float64(tb),
		)
	case bytes >= gb:
		return fmt.Sprintf(
			"%.2f GB",
			float64(bytes)/float64(gb),
		)
	case bytes >= mb:
		return fmt.Sprintf(
			"%.2f MB",
			float64(bytes)/float64(mb),
		)
	case bytes >= kb:
		return fmt.Sprintf(
			"%.2f KB",
			float64(bytes)/float64(kb),
		)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func formatIntSlice(ids []int) string {
	if len(ids) == 0 {
		return "(none)"
	}
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(strs, ", ")
}

func formatStringSlice(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	return strings.Join(items, ", ")
}

func formatPaginatedList[T any](
	list *models.PaginatedList[T],
	emptyMsg, header string,
	formatItem func(T) string,
) string {
	if list.Count == 0 {
		return emptyMsg
	}

	out := fmt.Sprintf("%s: %d total\n\n", header, list.Count)
	for _, item := range list.Results {
		out += formatItem(item)
	}

	if list.Next != nil {
		out += paginationHint
	}

	return out
}

func formatMatchableFields(
	label string,
	id int,
	name, slug, match string,
	algo int,
	isInsensitive bool,
	docCount int,
) string {
	algoName := matchingAlgorithmName(algo)
	matchDisplay := match
	if matchDisplay == "" {
		matchDisplay = "(none)"
	}

	out := fmt.Sprintf("%s (ID: %d)\n", label, id)
	out += fmt.Sprintf("  Name: %s\n", name)
	out += fmt.Sprintf("  Slug: %s\n", slug)
	out += fmt.Sprintf("  Match: %s\n", matchDisplay)
	out += fmt.Sprintf(
		"  Matching Algorithm: %d (%s)\n",
		algo,
		algoName,
	)
	out += fmt.Sprintf("  Case Insensitive: %v\n", isInsensitive)
	out += fmt.Sprintf("  Document Count: %d\n", docCount)

	return out
}
