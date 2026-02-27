package tools

import (
	"fmt"
	"sort"
	"strings"
)

func formatStatistics(
	stats map[string]interface{},
) string {
	if len(stats) == 0 {
		return "No statistics available."
	}

	out := "Paperless-NGX Statistics\n\n"

	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		label := formatStatLabel(key)
		switch v := stats[key].(type) {
		case []interface{}:
			out += fmt.Sprintf(
				"  %s: %s\n",
				label,
				formatStatSlice(v),
			)
		default:
			out += fmt.Sprintf(
				"  %s: %s\n",
				label,
				formatStatValue(v),
			)
		}
	}

	return out
}

func formatStatLabel(key string) string {
	words := strings.Split(key, "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func formatStatValue(v interface{}) string {
	if f, ok := v.(float64); ok {
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%.2f", f)
	}
	return fmt.Sprintf("%v", v)
}

func formatStatSlice(items []interface{}) string {
	if len(items) == 0 {
		return "(none)"
	}

	parts := make([]string, len(items))
	for i, item := range items {
		switch v := item.(type) {
		case map[string]interface{}:
			pairs := make([]string, 0, len(v))
			for k, val := range v {
				pairs = append(
					pairs,
					fmt.Sprintf(
						"%s=%s",
						k,
						formatStatValue(val),
					),
				)
			}
			sort.Strings(pairs)
			parts[i] = strings.Join(pairs, ", ")
		default:
			parts[i] = formatStatValue(item)
		}
	}
	return strings.Join(parts, "; ")
}
