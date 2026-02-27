// Package models defines the data structures for the paperless-ngx-mcp server.
package models

// SystemStatus represents the response from the Paperless-NGX /api/status/
// endpoint.
type SystemStatus struct {
	PNGXVersion string         `json:"pngx_version"`
	ServerOS    string         `json:"server_os"`
	InstallType string         `json:"install_type"`
	Storage     StorageStatus  `json:"storage"`
	Database    DatabaseStatus `json:"database"`
	Tasks       TasksStatus    `json:"tasks"`
}

// StorageStatus represents the storage section of the status response.
type StorageStatus struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
}

// DatabaseStatus represents the database section of the status response.
type DatabaseStatus struct {
	Type            string          `json:"type"`
	URL             string          `json:"url"`
	Status          string          `json:"status"`
	Error           *string         `json:"error"`
	MigrationStatus MigrationStatus `json:"migration_status"`
}

// MigrationStatus represents the migration section of the database status.
type MigrationStatus struct {
	LatestMigration     string   `json:"latest_migration"`
	UnappliedMigrations []string `json:"unapplied_migrations"`
}

// TasksStatus represents the tasks section of the status response.
// The index, classifier, and sanity_check fields are flattened into
// this object rather than being separate nested objects.
type TasksStatus struct {
	RedisURL    string  `json:"redis_url"`
	RedisStatus string  `json:"redis_status"`
	RedisError  *string `json:"redis_error"`

	CeleryStatus string  `json:"celery_status"`
	CeleryURL    *string `json:"celery_url"`
	CeleryError  *string `json:"celery_error"`

	IndexStatus       string  `json:"index_status"`
	IndexLastModified *string `json:"index_last_modified"`
	IndexError        *string `json:"index_error"`

	ClassifierStatus      string  `json:"classifier_status"`
	ClassifierLastTrained *string `json:"classifier_last_trained"`
	ClassifierError       *string `json:"classifier_error"`

	SanityCheckStatus  string  `json:"sanity_check_status"`
	SanityCheckLastRun *string `json:"sanity_check_last_run"`
	SanityCheckError   *string `json:"sanity_check_error"`
}
