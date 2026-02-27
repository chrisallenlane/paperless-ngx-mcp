package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSystemStatusUnmarshal(t *testing.T) {
	jsonData := `{
		"pngx_version": "2.20.8",
		"server_os": "Linux-4.4.302+-x86_64-with-glibc2.41",
		"install_type": "docker",
		"storage": {
			"total": 11518122557440,
			"available": 8525312483328
		},
		"database": {
			"type": "sqlite",
			"url": "/usr/src/paperless/data/db.sqlite3",
			"status": "OK",
			"error": null,
			"migration_status": {
				"latest_migration": "documents.0042_auto",
				"unapplied_migrations": []
			}
		},
		"tasks": {
			"redis_url": "redis://redis:6379",
			"redis_status": "OK",
			"redis_error": null,
			"celery_status": "OK",
			"celery_url": "celery@worker",
			"celery_error": null,
			"index_status": "OK",
			"index_last_modified": "2026-02-27T12:00:00Z",
			"index_error": null,
			"classifier_status": "OK",
			"classifier_last_trained": "2026-02-27T10:00:00Z",
			"classifier_error": null,
			"sanity_check_status": "OK",
			"sanity_check_last_run": "2026-02-22T06:00:00Z",
			"sanity_check_error": null
		}
	}`

	var status SystemStatus
	if err := json.Unmarshal([]byte(jsonData), &status); err != nil {
		t.Fatalf("Failed to unmarshal SystemStatus: %v", err)
	}

	if status.PNGXVersion != "2.20.8" {
		t.Errorf(
			"PNGXVersion = %s, want 2.20.8",
			status.PNGXVersion,
		)
	}

	if status.ServerOS != "Linux-4.4.302+-x86_64-with-glibc2.41" {
		t.Errorf("ServerOS = %s", status.ServerOS)
	}

	if status.InstallType != "docker" {
		t.Errorf("InstallType = %s, want docker", status.InstallType)
	}

	if status.Storage.Total != 11518122557440 {
		t.Errorf(
			"Storage.Total = %d, want 11518122557440",
			status.Storage.Total,
		)
	}

	if status.Storage.Available != 8525312483328 {
		t.Errorf(
			"Storage.Available = %d, want 8525312483328",
			status.Storage.Available,
		)
	}

	if status.Database.Type != "sqlite" {
		t.Errorf(
			"Database.Type = %s, want sqlite",
			status.Database.Type,
		)
	}

	if status.Database.Status != "OK" {
		t.Errorf(
			"Database.Status = %s, want OK",
			status.Database.Status,
		)
	}

	if status.Database.Error != nil {
		t.Errorf(
			"Database.Error = %v, want nil",
			status.Database.Error,
		)
	}

	if status.Database.MigrationStatus.LatestMigration != "documents.0042_auto" {
		t.Errorf(
			"LatestMigration = %s",
			status.Database.MigrationStatus.LatestMigration,
		)
	}

	if len(status.Database.MigrationStatus.UnappliedMigrations) != 0 {
		t.Errorf(
			"UnappliedMigrations len = %d, want 0",
			len(status.Database.MigrationStatus.UnappliedMigrations),
		)
	}

	if status.Tasks.RedisStatus != "OK" {
		t.Errorf(
			"Tasks.RedisStatus = %s, want OK",
			status.Tasks.RedisStatus,
		)
	}

	if status.Tasks.CeleryStatus != "OK" {
		t.Errorf(
			"Tasks.CeleryStatus = %s, want OK",
			status.Tasks.CeleryStatus,
		)
	}

	if status.Tasks.IndexStatus != "OK" {
		t.Errorf(
			"Tasks.IndexStatus = %s, want OK",
			status.Tasks.IndexStatus,
		)
	}

	if status.Tasks.IndexLastModified == nil ||
		*status.Tasks.IndexLastModified != "2026-02-27T12:00:00Z" {
		t.Errorf(
			"Tasks.IndexLastModified = %v",
			status.Tasks.IndexLastModified,
		)
	}

	if status.Tasks.ClassifierStatus != "OK" {
		t.Errorf(
			"Tasks.ClassifierStatus = %s, want OK",
			status.Tasks.ClassifierStatus,
		)
	}

	if status.Tasks.SanityCheckStatus != "OK" {
		t.Errorf(
			"Tasks.SanityCheckStatus = %s, want OK",
			status.Tasks.SanityCheckStatus,
		)
	}
}

func TestSystemStatusUnmarshalWithErrors(t *testing.T) {
	dbErr := "Error connecting to database"
	jsonData := `{
		"pngx_version": "2.20.8",
		"server_os": "Linux",
		"install_type": "docker",
		"storage": {"total": 100, "available": 50},
		"database": {
			"type": "sqlite",
			"url": "/data/db.sqlite3",
			"status": "ERROR",
			"error": "Error connecting to database",
			"migration_status": {
				"latest_migration": "",
				"unapplied_migrations": []
			}
		},
		"tasks": {
			"redis_url": "redis://localhost",
			"redis_status": "OK",
			"redis_error": null,
			"celery_status": "OK",
			"celery_url": null,
			"celery_error": null,
			"index_status": "OK",
			"index_last_modified": null,
			"index_error": null,
			"classifier_status": "WARNING",
			"classifier_last_trained": null,
			"classifier_error": null,
			"sanity_check_status": "OK",
			"sanity_check_last_run": null,
			"sanity_check_error": null
		}
	}`

	var status SystemStatus
	if err := json.Unmarshal([]byte(jsonData), &status); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if status.Database.Error == nil || *status.Database.Error != dbErr {
		t.Errorf("Database.Error = %v, want %q", status.Database.Error, dbErr)
	}

	if status.Tasks.IndexLastModified != nil {
		t.Errorf(
			"IndexLastModified = %v, want nil",
			status.Tasks.IndexLastModified,
		)
	}

	if status.Tasks.CeleryURL != nil {
		t.Errorf("CeleryURL = %v, want nil", status.Tasks.CeleryURL)
	}
}

func TestApplicationConfigurationUnmarshal(t *testing.T) {
	jsonData := `[{
		"id": 1,
		"output_type": "pdfa",
		"pages": 5,
		"language": "eng+deu",
		"mode": "skip",
		"skip_archive_file": "with_text",
		"image_dpi": 300,
		"unpaper_clean": "clean",
		"deskew": true,
		"rotate_pages": false,
		"rotate_pages_threshold": 12.5,
		"max_image_pixels": 500000000.0,
		"color_conversion_strategy": "RGB",
		"user_args": {"--deskew": true},
		"app_title": "My Paperless",
		"app_logo": "/media/logo/custom.png",
		"barcodes_enabled": true,
		"barcode_enable_tiff_support": false,
		"barcode_string": "PATCHT",
		"barcode_retain_split_pages": true,
		"barcode_enable_asn": true,
		"barcode_asn_prefix": "ASN",
		"barcode_upscale": 1.5,
		"barcode_dpi": 200,
		"barcode_max_pages": 10,
		"barcode_enable_tag": false,
		"barcode_tag_mapping": {"ASN": "tag1"}
	}]`

	var configs []ApplicationConfiguration
	if err := json.Unmarshal([]byte(jsonData), &configs); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(configs))
	}

	c := configs[0]

	if c.ID != 1 {
		t.Errorf("ID = %d, want 1", c.ID)
	}

	if c.OutputType == nil || *c.OutputType != "pdfa" {
		t.Errorf("OutputType = %v, want pdfa", c.OutputType)
	}

	if c.Pages == nil || *c.Pages != 5 {
		t.Errorf("Pages = %v, want 5", c.Pages)
	}

	if c.Language == nil || *c.Language != "eng+deu" {
		t.Errorf("Language = %v, want eng+deu", c.Language)
	}

	if c.Deskew == nil || *c.Deskew != true {
		t.Errorf("Deskew = %v, want true", c.Deskew)
	}

	if c.RotatePages == nil || *c.RotatePages != false {
		t.Errorf("RotatePages = %v, want false", c.RotatePages)
	}

	if c.RotatePagesThreshold == nil || *c.RotatePagesThreshold != 12.5 {
		t.Errorf(
			"RotatePagesThreshold = %v, want 12.5",
			c.RotatePagesThreshold,
		)
	}

	if c.AppTitle == nil || *c.AppTitle != "My Paperless" {
		t.Errorf("AppTitle = %v, want My Paperless", c.AppTitle)
	}

	if c.BarcodesEnabled == nil || *c.BarcodesEnabled != true {
		t.Errorf(
			"BarcodesEnabled = %v, want true",
			c.BarcodesEnabled,
		)
	}

	if c.BarcodeUpscale == nil || *c.BarcodeUpscale != 1.5 {
		t.Errorf(
			"BarcodeUpscale = %v, want 1.5",
			c.BarcodeUpscale,
		)
	}

	if c.UserArgs == nil ||
		!strings.Contains(string(c.UserArgs), "--deskew") {
		t.Errorf("UserArgs = %v, want JSON with --deskew", c.UserArgs)
	}

	if c.BarcodeTagMapping == nil ||
		!strings.Contains(string(c.BarcodeTagMapping), "ASN") {
		t.Errorf(
			"BarcodeTagMapping = %v, want JSON with ASN",
			c.BarcodeTagMapping,
		)
	}
}

func TestCorrespondentUnmarshal(t *testing.T) {
	jsonData := `{
		"id": 1,
		"slug": "acme-corp",
		"name": "ACME Corp",
		"match": "acme",
		"matching_algorithm": 1,
		"is_insensitive": true,
		"document_count": 5,
		"last_correspondence": "2026-02-15T10:00:00Z"
	}`

	var c Correspondent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Failed to unmarshal Correspondent: %v", err)
	}

	if c.ID != 1 {
		t.Errorf("ID = %d, want 1", c.ID)
	}

	if c.Name != "ACME Corp" {
		t.Errorf("Name = %s, want ACME Corp", c.Name)
	}

	if c.Slug != "acme-corp" {
		t.Errorf("Slug = %s, want acme-corp", c.Slug)
	}

	if c.Match != "acme" {
		t.Errorf("Match = %s, want acme", c.Match)
	}

	if c.MatchingAlgorithm != 1 {
		t.Errorf(
			"MatchingAlgorithm = %d, want 1",
			c.MatchingAlgorithm,
		)
	}

	if !c.IsInsensitive {
		t.Error("IsInsensitive = false, want true")
	}

	if c.DocumentCount != 5 {
		t.Errorf("DocumentCount = %d, want 5", c.DocumentCount)
	}

	if c.LastCorrespondence == nil ||
		*c.LastCorrespondence != "2026-02-15T10:00:00Z" {
		t.Errorf(
			"LastCorrespondence = %v, want 2026-02-15T10:00:00Z",
			c.LastCorrespondence,
		)
	}
}

func TestCorrespondentUnmarshalNullLastCorrespondence(t *testing.T) {
	jsonData := `{
		"id": 2,
		"slug": "john-doe",
		"name": "John Doe",
		"match": "",
		"matching_algorithm": 6,
		"is_insensitive": true,
		"document_count": 0,
		"last_correspondence": null
	}`

	var c Correspondent
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if c.LastCorrespondence != nil {
		t.Errorf(
			"LastCorrespondence = %v, want nil",
			c.LastCorrespondence,
		)
	}
}

func TestPaginatedListCorrespondentUnmarshal(t *testing.T) {
	jsonData := `{
		"count": 1,
		"next": null,
		"previous": null,
		"all": [1],
		"results": [{
			"id": 1,
			"slug": "acme-corp",
			"name": "ACME Corp",
			"match": "",
			"matching_algorithm": 1,
			"is_insensitive": true,
			"document_count": 0
		}]
	}`

	var list PaginatedList[Correspondent]
	if err := json.Unmarshal([]byte(jsonData), &list); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if list.Count != 1 {
		t.Errorf("Count = %d, want 1", list.Count)
	}

	if list.Next != nil {
		t.Errorf("Next = %v, want nil", list.Next)
	}

	if len(list.All) != 1 || list.All[0] != 1 {
		t.Errorf("All = %v, want [1]", list.All)
	}

	if len(list.Results) != 1 {
		t.Fatalf("Results len = %d, want 1", len(list.Results))
	}

	if list.Results[0].Name != "ACME Corp" {
		t.Errorf(
			"Results[0].Name = %s, want ACME Corp",
			list.Results[0].Name,
		)
	}
}

func TestCustomFieldUnmarshal(t *testing.T) {
	jsonData := `{
		"id": 1,
		"name": "Invoice Number",
		"data_type": "string",
		"extra_data": {"select_options": ["a", "b"]},
		"document_count": 10
	}`

	var cf CustomField
	if err := json.Unmarshal([]byte(jsonData), &cf); err != nil {
		t.Fatalf("Failed to unmarshal CustomField: %v", err)
	}

	if cf.ID != 1 {
		t.Errorf("ID = %d, want 1", cf.ID)
	}

	if cf.Name != "Invoice Number" {
		t.Errorf("Name = %s, want Invoice Number", cf.Name)
	}

	if cf.DataType != "string" {
		t.Errorf("DataType = %s, want string", cf.DataType)
	}

	if cf.DocumentCount != 10 {
		t.Errorf("DocumentCount = %d, want 10", cf.DocumentCount)
	}

	if cf.ExtraData == nil ||
		!strings.Contains(string(cf.ExtraData), "select_options") {
		t.Errorf(
			"ExtraData = %v, want JSON with select_options",
			cf.ExtraData,
		)
	}
}

func TestCustomFieldUnmarshalNullExtraData(t *testing.T) {
	jsonData := `{
		"id": 2,
		"name": "Due Date",
		"data_type": "date",
		"extra_data": null,
		"document_count": 0
	}`

	var cf CustomField
	if err := json.Unmarshal([]byte(jsonData), &cf); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if string(cf.ExtraData) != "null" {
		t.Errorf("ExtraData = %s, want null", string(cf.ExtraData))
	}
}

func TestApplicationConfigurationUnmarshalAllNulls(t *testing.T) {
	jsonData := `[{
		"id": 1,
		"output_type": null,
		"pages": null,
		"language": null,
		"mode": null,
		"skip_archive_file": null,
		"image_dpi": null,
		"unpaper_clean": null,
		"deskew": null,
		"rotate_pages": null,
		"rotate_pages_threshold": null,
		"max_image_pixels": null,
		"color_conversion_strategy": null,
		"user_args": null,
		"app_title": null,
		"app_logo": null,
		"barcodes_enabled": null,
		"barcode_enable_tiff_support": null,
		"barcode_string": null,
		"barcode_retain_split_pages": null,
		"barcode_enable_asn": null,
		"barcode_asn_prefix": null,
		"barcode_upscale": null,
		"barcode_dpi": null,
		"barcode_max_pages": null,
		"barcode_enable_tag": null,
		"barcode_tag_mapping": null
	}]`

	var configs []ApplicationConfiguration
	if err := json.Unmarshal([]byte(jsonData), &configs); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	c := configs[0]

	if c.ID != 1 {
		t.Errorf("ID = %d, want 1", c.ID)
	}

	if c.OutputType != nil {
		t.Errorf("OutputType = %v, want nil", c.OutputType)
	}

	if c.Pages != nil {
		t.Errorf("Pages = %v, want nil", c.Pages)
	}

	if c.Deskew != nil {
		t.Errorf("Deskew = %v, want nil", c.Deskew)
	}

	if c.AppTitle != nil {
		t.Errorf("AppTitle = %v, want nil", c.AppTitle)
	}

	if c.BarcodesEnabled != nil {
		t.Errorf("BarcodesEnabled = %v, want nil", c.BarcodesEnabled)
	}
}
