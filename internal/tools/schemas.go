package tools

// withIDForUpdate conditionally adds an "id" field to a schema
// properties map and adjusts the required fields. When includeID
// is true, the schema becomes an update schema (id required);
// otherwise the original createRequired fields are used.
func withIDForUpdate(
	props map[string]interface{},
	idDesc string,
	includeID bool,
	createRequired []string,
) []string {
	if includeID {
		props["id"] = map[string]interface{}{
			"type":        "integer",
			"description": idDesc,
		}
		return []string{"id"}
	}
	return createRequired
}

const matchingAlgorithmDesc = "Matching algorithm: " +
	"0=None, 1=Any word, 2=All words, " +
	"3=Exact match, 4=Regex, " +
	"5=Fuzzy word, 6=Automatic"

// addMatchableProps adds the shared matching fields (match,
// matching_algorithm, is_insensitive) to a schema properties
// map.
func addMatchableProps(
	props map[string]interface{},
) {
	props["match"] = map[string]interface{}{
		"type":        "string",
		"description": "Match pattern for auto-assignment",
	}
	props["matching_algorithm"] = map[string]interface{}{
		"type":        "integer",
		"description": matchingAlgorithmDesc,
	}
	props["is_insensitive"] = map[string]interface{}{
		"type":        "boolean",
		"description": "Case-insensitive matching",
	}
}

// addPaginationProps adds page and page_size fields to a
// schema properties map.
func addPaginationProps(
	props map[string]interface{},
) {
	props["page"] = map[string]interface{}{
		"type":        "integer",
		"description": "Page number (default 1)",
	}
	props["page_size"] = map[string]interface{}{
		"type":        "integer",
		"description": "Results per page (default 25)",
	}
}

// emptySchema returns an input schema with no parameters.
func emptySchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// idOnlySchema returns an input schema with a single required "id" field.
func idOnlySchema(desc string) map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": desc,
			},
		},
		"required": []string{"id"},
	}
}

// paginationOnlySchema returns an input schema with only page and
// page_size parameters (no name filter).
func paginationOnlySchema() map[string]interface{} {
	props := map[string]interface{}{}
	addPaginationProps(props)
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
}

// paginatedListSchema returns an input schema for paginated list endpoints.
func paginatedListSchema() map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type": "string",
			"description": "Filter by name " +
				"(case-insensitive contains)",
		},
	}
	addPaginationProps(props)
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
}

// matchableResourceSchema returns an input schema for resources with matching
// fields (name, match, matching_algorithm, is_insensitive). Set includeID
// to true for update tools, false for create tools.
func matchableResourceSchema(
	resourceName string,
	includeID bool,
) map[string]interface{} {
	props := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": resourceName + " name",
		},
	}
	addMatchableProps(props)

	required := withIDForUpdate(
		props,
		resourceName+" ID to update",
		includeID,
		[]string{"name"},
	)

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}
