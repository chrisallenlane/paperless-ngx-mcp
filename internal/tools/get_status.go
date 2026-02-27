package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// GetStatus retrieves the system status from Paperless-NGX.
type GetStatus struct {
	client *client.Client
}

// NewGetStatus creates a new GetStatus tool instance.
func NewGetStatus(c *client.Client) *GetStatus {
	return &GetStatus{client: c}
}

// Description returns a description of what this tool does.
func (t *GetStatus) Description() string {
	return "Get the current system status of the Paperless-NGX server"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *GetStatus) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// Execute runs the tool and returns a formatted status summary.
func (t *GetStatus) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, "/api/status/")
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	var status models.SystemStatus
	if err := ParseJSONResponse(body, &status); err != nil {
		return "", fmt.Errorf(
			"failed to parse status response: %w",
			err,
		)
	}

	return formatStatus(&status), nil
}

func formatStatus(s *models.SystemStatus) string {
	totalTB := float64(s.Storage.Total) / (1024 * 1024 * 1024 * 1024)
	availTB := float64(s.Storage.Available) / (1024 * 1024 * 1024 * 1024)

	out := fmt.Sprintf(
		"Paperless-NGX Status\nVersion: %s\nOS: %s\nInstall: %s\n\n"+
			"Storage: %.2f TB available of %.2f TB\n\n"+
			"Database: %s - %s\n"+
			"Redis: %s\n"+
			"Celery: %s\n",
		s.PNGXVersion,
		s.ServerOS,
		s.InstallType,
		availTB,
		totalTB,
		s.Database.Type,
		s.Database.Status,
		s.Tasks.RedisStatus,
		s.Tasks.CeleryStatus,
	)

	out += formatTaskLine(
		"Index",
		s.Tasks.IndexStatus,
		"last modified",
		s.Tasks.IndexLastModified,
	)
	out += formatTaskLine(
		"Classifier",
		s.Tasks.ClassifierStatus,
		"last trained",
		s.Tasks.ClassifierLastTrained,
	)
	out += formatTaskLine(
		"Sanity Check",
		s.Tasks.SanityCheckStatus,
		"last run",
		s.Tasks.SanityCheckLastRun,
	)

	return out
}

func formatTaskLine(
	name, status, dateLabel string,
	dateValue *string,
) string {
	if dateValue != nil && *dateValue != "" {
		date := formatDate(*dateValue)
		return fmt.Sprintf(
			"%s: %s (%s: %s)\n",
			name,
			status,
			dateLabel,
			date,
		)
	}
	return fmt.Sprintf("%s: %s\n", name, status)
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}
