package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

// DownloadDocument downloads a document file from Paperless-NGX.
type DownloadDocument struct {
	client *client.Client
}

// NewDownloadDocument creates a new DownloadDocument tool instance.
func NewDownloadDocument(
	c *client.Client,
) *DownloadDocument {
	return &DownloadDocument{client: c}
}

// Description returns a description of what this tool does.
func (t *DownloadDocument) Description() string {
	return "Download a document file from Paperless-NGX " +
		"to the local filesystem"
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (t *DownloadDocument) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "integer",
				"description": "Document ID to download",
			},
			"original": map[string]interface{}{
				"type": "boolean",
				"description": "Download original file " +
					"instead of archive version " +
					"(default: false)",
			},
			"save_path": map[string]interface{}{
				"type": "string",
				"description": "Filesystem path where " +
					"the file should be saved",
			},
		},
		"required": []string{"id", "save_path"},
	}
}

// Execute runs the tool and returns a confirmation message.
func (t *DownloadDocument) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		ID       int    `json:"id"`
		Original bool   `json:"original"`
		SavePath string `json:"save_path"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	if params.ID <= 0 {
		return "", fmt.Errorf("id must be a positive integer")
	}

	if params.SavePath == "" {
		return "", fmt.Errorf("save_path is required")
	}

	if err := validateFilePath(params.SavePath); err != nil {
		return "", err
	}

	dir := filepath.Dir(params.SavePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf(
			"parent directory does not exist: %s",
			dir,
		)
	}

	path := fmt.Sprintf(
		"/api/documents/%d/download/",
		params.ID,
	)
	if params.Original {
		path += "?original=true"
	}

	resp, err := t.client.Get(ctx, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to download document: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"download failed with status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	outFile, err := os.Create(params.SavePath)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create output file: %w",
			err,
		)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "unknown"
	}

	return fmt.Sprintf(
		"Document downloaded successfully.\n"+
			"Saved to: %s\n"+
			"Size: %s\n"+
			"Content-Type: %s",
		params.SavePath,
		formatFileSize(int(written)),
		contentType,
	), nil
}
