package tools

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ptrStr returns a pointer to the given string value.
func ptrStr(s string) *string {
	return &s
}

// --- Gap 6: Format edge cases ---

func TestFormatOptDate(t *testing.T) {
	tests := []struct {
		name string
		v    *string
		want string
	}{
		{
			name: "nil pointer returns (none)",
			v:    nil,
			want: "(none)",
		},
		{
			name: "empty string pointer returns (none)",
			v:    ptrStr(""),
			want: "(none)",
		},
		{
			name: "valid date string returns truncated date",
			v:    ptrStr("2024-01-15T10:30:00Z"),
			want: "2024-01-15",
		},
		{
			name: "exactly 10 chars returns as-is",
			v:    ptrStr("2024-01-15"),
			want: "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatOptDate(tt.v)
			if got != tt.want {
				t.Errorf(
					"formatOptDate() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string returns empty string",
			input: "",
			want:  "",
		},
		{
			name:  "short string under 10 chars returns as-is",
			input: "2024",
			want:  "2024",
		},
		{
			name:  "9 chars returns as-is",
			input: "2024-01-1",
			want:  "2024-01-1",
		},
		{
			name:  "exactly 10 chars returns all 10",
			input: "2024-01-15",
			want:  "2024-01-15",
		},
		{
			name:  "full ISO timestamp truncated to date part",
			input: "2024-01-15T10:30:00Z",
			want:  "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDate(tt.input)
			if got != tt.want {
				t.Errorf(
					"formatDate(%q) = %q, want %q",
					tt.input,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{
			name:  "bytes",
			bytes: 512,
			want:  "512 bytes",
		},
		{
			name:  "kilobytes",
			bytes: 2 * 1024,
			want:  "2.00 KB",
		},
		{
			name:  "megabytes",
			bytes: 3 * 1024 * 1024,
			want:  "3.00 MB",
		},
		{
			name:  "5 GB",
			bytes: 5 * 1024 * 1024 * 1024,
			want:  "5.00 GB",
		},
		{
			name:  "2 TB",
			bytes: 2 * 1024 * 1024 * 1024 * 1024,
			want:  "2.00 TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFileSize(tt.bytes)
			if got != tt.want {
				t.Errorf(
					"formatFileSize(%d) = %q, want %q",
					tt.bytes,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatNoteList_Empty(t *testing.T) {
	got := formatNoteList(1, []models.Note{})
	expected := "No notes found for document 1."
	if got != expected {
		t.Errorf(
			"formatNoteList(empty) = %q, want %q",
			got,
			expected,
		)
	}
}

func TestFormatNoteList_WithNotes(t *testing.T) {
	notes := []models.Note{
		{
			ID:      1,
			Note:    "Only note",
			Created: "2024-01-15T00:00:00Z",
			User: models.BasicUser{
				ID:       1,
				Username: "alice",
			},
		},
	}

	got := formatNoteList(1, notes)

	if !strings.Contains(got, "1 total") {
		t.Errorf(
			"formatNoteList should contain count, got:\n%s",
			got,
		)
	}
	if !strings.Contains(got, "Only note") {
		t.Errorf(
			"formatNoteList should contain note text, got:\n%s",
			got,
		)
	}
}

func TestRuleTypeName(t *testing.T) {
	tests := []struct {
		name     string
		ruleType int
		want     string
	}{
		{
			name:     "known type 0",
			ruleType: 0,
			want:     "Title contains",
		},
		{
			name:     "known type 1",
			ruleType: 1,
			want:     "Content contains",
		},
		{
			name:     "known type 47",
			ruleType: 47,
			want:     "Has custom fields in",
		},
		{
			name:     "unknown type 999 returns fallback",
			ruleType: 999,
			want:     "Rule type 999",
		},
		{
			name:     "unknown type 100 returns fallback",
			ruleType: 100,
			want:     "Rule type 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ruleTypeName(tt.ruleType)
			if got != tt.want {
				t.Errorf(
					"ruleTypeName(%d) = %q, want %q",
					tt.ruleType,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatStatValue(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{
			name: "integer-valued float64 formats as integer",
			v:    float64(42),
			want: "42",
		},
		{
			name: "non-integer float64 formats with 2 decimal places",
			v:    float64(3.14),
			want: "3.14",
		},
		{
			name: "zero float64 formats as integer",
			v:    float64(0),
			want: "0",
		},
		{
			name: "non-float value formatted with %v",
			v:    "some string",
			want: "some string",
		},
		{
			name: "bool value formatted with %v",
			v:    true,
			want: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatValue(tt.v)
			if got != tt.want {
				t.Errorf(
					"formatStatValue(%v) = %q, want %q",
					tt.v,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatStatSlice(t *testing.T) {
	tests := []struct {
		name  string
		items []interface{}
		want  string
	}{
		{
			name:  "empty slice returns (none)",
			items: []interface{}{},
			want:  "(none)",
		},
		{
			name:  "nil slice returns (none)",
			items: nil,
			want:  "(none)",
		},
		{
			name:  "single scalar item",
			items: []interface{}{float64(7)},
			want:  "7",
		},
		{
			name:  "multiple scalar items joined with semicolons",
			items: []interface{}{float64(1), float64(2)},
			want:  "1; 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatSlice(tt.items)
			if got != tt.want {
				t.Errorf(
					"formatStatSlice(%v) = %q, want %q",
					tt.items,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFormatTaskLine(t *testing.T) {
	tests := []struct {
		name         string
		taskName     string
		status       string
		dateLabel    string
		dateValue    *string
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:         "nil dateValue omits parenthetical",
			taskName:     "Index",
			status:       "OK",
			dateLabel:    "last modified",
			dateValue:    nil,
			wantContains: []string{"Index", "OK"},
			wantAbsent:   []string{"last modified"},
		},
		{
			name:         "empty string dateValue omits parenthetical",
			taskName:     "Classifier",
			status:       "OK",
			dateLabel:    "last trained",
			dateValue:    ptrStr(""),
			wantContains: []string{"Classifier", "OK"},
			wantAbsent:   []string{"last trained"},
		},
		{
			name:      "non-empty dateValue includes parenthetical",
			taskName:  "Sanity Check",
			status:    "OK",
			dateLabel: "last run",
			dateValue: ptrStr("2024-01-15T00:00:00Z"),
			wantContains: []string{
				"Sanity Check",
				"OK",
				"last run",
				"2024-01-15",
			},
			wantAbsent: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTaskLine(
				tt.taskName,
				tt.status,
				tt.dateLabel,
				tt.dateValue,
			)
			for _, s := range tt.wantContains {
				if !strings.Contains(got, s) {
					t.Errorf(
						"formatTaskLine() output missing %q, got: %q",
						s,
						got,
					)
				}
			}
			for _, s := range tt.wantAbsent {
				if strings.Contains(got, s) {
					t.Errorf(
						"formatTaskLine() output should not "+
							"contain %q, got: %q",
						s,
						got,
					)
				}
			}
		})
	}
}

func TestFormatDocument_ContentTruncation(t *testing.T) {
	// Build content that is exactly 501 characters to trigger truncation.
	longContent := strings.Repeat("x", 501)

	doc := &models.Document{
		ID:      42,
		Title:   "Test Doc",
		Content: longContent,
	}

	got := formatDocument(doc)

	if !strings.Contains(got, "...") {
		t.Errorf(
			"formatDocument with content > 500 chars should "+
				"contain '...', got:\n%s",
			got,
		)
	}

	// The truncated content should be 500 chars of "x" plus "...".
	expectedSnippet := strings.Repeat("x", 500) + "..."
	if !strings.Contains(got, expectedSnippet) {
		t.Errorf(
			"formatDocument truncated content does not match "+
				"expected snippet, got:\n%s",
			got,
		)
	}
}

func TestFormatDocument_ContentNotTruncated(t *testing.T) {
	shortContent := strings.Repeat("y", 500)

	doc := &models.Document{
		ID:      1,
		Title:   "Short Doc",
		Content: shortContent,
	}

	got := formatDocument(doc)

	if strings.Contains(got, "...") {
		t.Errorf(
			"formatDocument with content == 500 chars should "+
				"not truncate, got:\n%s",
			got,
		)
	}
}

// --- Gap 9: Validate function JSON parse failures ---

func TestValidateFunctions_MalformedJSON(t *testing.T) {
	badJSON := json.RawMessage("not json")

	tests := []struct {
		name     string
		validate func(json.RawMessage) (interface{}, error)
	}{
		{
			name:     "validateMatchableCreate",
			validate: validateMatchableCreate,
		},
		{
			name:     "validateCreateTag",
			validate: validateCreateTag,
		},
		{
			name:     "validateCreateStoragePath",
			validate: validateCreateStoragePath,
		},
		{
			name:     "validateCreateCustomField",
			validate: validateCreateCustomField,
		},
		{
			name:     "validateCreateSavedView",
			validate: validateCreateSavedView,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.validate(badJSON)
			if err == nil {
				t.Fatalf(
					"%s: expected error for malformed JSON, "+
						"got nil",
					tt.name,
				)
			}
			if !strings.Contains(
				err.Error(),
				"failed to parse arguments",
			) {
				t.Errorf(
					"%s: error = %q, want error containing "+
						"\"failed to parse arguments\"",
					tt.name,
					err.Error(),
				)
			}
		})
	}
}

// --- Gap 10: Path builder malformed JSON ---

func TestBuildDocumentListPath_MalformedJSON(t *testing.T) {
	_, err := buildDocumentListPath(
		json.RawMessage("not json"),
	)
	if err == nil {
		t.Fatal(
			"buildDocumentListPath: expected error for " +
				"malformed JSON, got nil",
		)
	}
}

func TestBuildTaskListPath_MalformedJSON(t *testing.T) {
	_, err := buildTaskListPath(
		json.RawMessage("not json"),
	)
	if err == nil {
		t.Fatal(
			"buildTaskListPath: expected error for " +
				"malformed JSON, got nil",
		)
	}
}
