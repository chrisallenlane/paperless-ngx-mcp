package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// UploadDocument uploads a document to Paperless-NGX.
type UploadDocument struct {
	client *client.Client
}

// NewUploadDocument creates a new UploadDocument tool instance.
func NewUploadDocument(c *client.Client) *UploadDocument {
	return &UploadDocument{client: c}
}

// Description returns a description of what this tool does.
func (t *UploadDocument) Description() string {
	return "Upload a document file to Paperless-NGX " +
		"for processing and indexing"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *UploadDocument) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{
				"type": "string",
				"description": "Path to the file " +
					"to upload",
			},
			"title": map[string]interface{}{
				"type": "string",
				"description": "Document title " +
					"(defaults to filename)",
			},
			"correspondent": map[string]interface{}{
				"type": "integer",
				"description": "Correspondent ID " +
					"to assign",
			},
			"document_type": map[string]interface{}{
				"type": "integer",
				"description": "Document type ID " +
					"to assign",
			},
			"storage_path": map[string]interface{}{
				"type": "integer",
				"description": "Storage path ID " +
					"to assign",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "integer",
				},
				"description": "Tag IDs to assign",
			},
			"archive_serial_number": map[string]interface{}{
				"type":        "integer",
				"description": "Archive serial number",
				"minimum":     0,
			},
			"created": map[string]interface{}{
				"type": "string",
				"description": "Override creation " +
					"date (ISO format)",
			},
		},
		"required": []string{"file_path"},
	}
}

// Execute runs the tool and returns a confirmation message.
func (t *UploadDocument) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		FilePath            string `json:"file_path"`
		Title               string `json:"title"`
		Correspondent       *int   `json:"correspondent"`
		DocumentType        *int   `json:"document_type"`
		StoragePath         *int   `json:"storage_path"`
		Tags                []int  `json:"tags"`
		ArchiveSerialNumber *int   `json:"archive_serial_number"`
		Created             string `json:"created"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.FilePath == "" {
		return "", fmt.Errorf("file_path is required")
	}

	if err := validateFilePath(params.FilePath); err != nil {
		return "", err
	}

	file, err := os.Open(params.FilePath)
	if err != nil {
		return "", fmt.Errorf(
			"failed to open file: %w",
			err,
		)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf(
			"failed to stat file: %w",
			err,
		)
	}

	if stat.IsDir() {
		return "", fmt.Errorf(
			"file_path must be a file, not a directory",
		)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(
		"document",
		filepath.Base(params.FilePath),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create form file: %w",
			err,
		)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf(
			"failed to write file to form: %w",
			err,
		)
	}

	if params.Title != "" {
		writer.WriteField("title", params.Title)
	}
	if params.Correspondent != nil {
		writer.WriteField(
			"correspondent",
			fmt.Sprintf("%d", *params.Correspondent),
		)
	}
	if params.DocumentType != nil {
		writer.WriteField(
			"document_type",
			fmt.Sprintf("%d", *params.DocumentType),
		)
	}
	if params.StoragePath != nil {
		writer.WriteField(
			"storage_path",
			fmt.Sprintf("%d", *params.StoragePath),
		)
	}
	for _, tagID := range params.Tags {
		writer.WriteField(
			"tags",
			fmt.Sprintf("%d", tagID),
		)
	}
	if params.ArchiveSerialNumber != nil {
		writer.WriteField(
			"archive_serial_number",
			fmt.Sprintf("%d", *params.ArchiveSerialNumber),
		)
	}
	if params.Created != "" {
		writer.WriteField("created", params.Created)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf(
			"failed to finalize multipart form: %w",
			err,
		)
	}

	resp, err := t.client.PostMultipart(
		ctx,
		"/api/documents/post_document/",
		&buf,
		writer.FormDataContentType(),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to upload document: %w",
			err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"failed to read response: %w",
			err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"upload failed with status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	taskID := string(body)
	return fmt.Sprintf(
		"Document uploaded successfully.\n"+
			"File: %s\n"+
			"Size: %s\n"+
			"Task ID: %s",
		filepath.Base(params.FilePath),
		formatFileSize(int(stat.Size())),
		taskID,
	), nil
}
