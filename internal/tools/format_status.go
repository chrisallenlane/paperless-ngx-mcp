package tools

import (
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

func formatStatus(s *models.SystemStatus) string {
	out := fmt.Sprintf(
		"Paperless-NGX Status\nVersion: %s\nOS: %s\nInstall: %s\n\n"+
			"Storage: %s available of %s\n\n"+
			"Database: %s - %s\n"+
			"Redis: %s\n"+
			"Celery: %s\n",
		s.PNGXVersion,
		s.ServerOS,
		s.InstallType,
		formatFileSize(s.Storage.Available),
		formatFileSize(s.Storage.Total),
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
