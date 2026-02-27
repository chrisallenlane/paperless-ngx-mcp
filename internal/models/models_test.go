package models

import (
	"encoding/json"
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
