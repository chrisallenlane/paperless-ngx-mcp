// Package models defines the data structures for the paperless-ngx-mcp server.
package models

import "encoding/json"

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

// ApplicationConfiguration represents the response from the Paperless-NGX
// /api/config/ endpoint. All fields except ID are nullable (nil means the
// server default is in use).
type ApplicationConfiguration struct {
	ID int `json:"id"`

	// OCR settings
	OutputType              *string         `json:"output_type"`
	Pages                   *int64          `json:"pages"`
	Language                *string         `json:"language"`
	Mode                    *string         `json:"mode"`
	SkipArchiveFile         *string         `json:"skip_archive_file"`
	ImageDPI                *int64          `json:"image_dpi"`
	UnpaperClean            *string         `json:"unpaper_clean"`
	Deskew                  *bool           `json:"deskew"`
	RotatePages             *bool           `json:"rotate_pages"`
	RotatePagesThreshold    *float64        `json:"rotate_pages_threshold"`
	MaxImagePixels          *float64        `json:"max_image_pixels"`
	ColorConversionStrategy *string         `json:"color_conversion_strategy"`
	UserArgs                json.RawMessage `json:"user_args"`

	// App settings
	AppTitle *string `json:"app_title"`
	AppLogo  *string `json:"app_logo"`

	// Barcode settings
	BarcodesEnabled          *bool           `json:"barcodes_enabled"`
	BarcodeEnableTiffSupport *bool           `json:"barcode_enable_tiff_support"`
	BarcodeString            *string         `json:"barcode_string"`
	BarcodeRetainSplitPages  *bool           `json:"barcode_retain_split_pages"`
	BarcodeEnableASN         *bool           `json:"barcode_enable_asn"`
	BarcodeASNPrefix         *string         `json:"barcode_asn_prefix"`
	BarcodeUpscale           *float64        `json:"barcode_upscale"`
	BarcodeDPI               *int64          `json:"barcode_dpi"`
	BarcodeMaxPages          *int64          `json:"barcode_max_pages"`
	BarcodeEnableTag         *bool           `json:"barcode_enable_tag"`
	BarcodeTagMapping        json.RawMessage `json:"barcode_tag_mapping"`
}

// Correspondent represents a Paperless-NGX correspondent.
type Correspondent struct {
	ID                 int     `json:"id"`
	Slug               string  `json:"slug"`
	Name               string  `json:"name"`
	Match              string  `json:"match"`
	MatchingAlgorithm  int     `json:"matching_algorithm"`
	IsInsensitive      bool    `json:"is_insensitive"`
	DocumentCount      int     `json:"document_count"`
	LastCorrespondence *string `json:"last_correspondence"`
}

// PaginatedList is a generic paginated API response.
type PaginatedList[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	All      []int   `json:"all"`
	Results  []T     `json:"results"`
}

// CustomField represents a Paperless-NGX custom field definition.
type CustomField struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	DataType      string          `json:"data_type"`
	ExtraData     json.RawMessage `json:"extra_data"`
	DocumentCount int             `json:"document_count"`
}
