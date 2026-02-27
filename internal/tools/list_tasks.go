package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// ListTasks lists background tasks from Paperless-NGX.
type ListTasks struct {
	client *client.Client
}

// NewListTasks creates a new ListTasks tool instance.
func NewListTasks(c *client.Client) *ListTasks {
	return &ListTasks{client: c}
}

// Description returns a description of what this tool does.
func (t *ListTasks) Description() string {
	return "List background tasks in Paperless-NGX " +
		"with optional status, name, and type filters"
}

// InputSchema returns the JSON schema for the tool's input
// parameters.
func (t *ListTasks) InputSchema() map[string]interface{} {
	return taskListSchema()
}

// Execute runs the tool and returns a formatted task list.
func (t *ListTasks) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	path, err := buildTaskListPath(args)
	if err != nil {
		return "", err
	}

	body, err := doAPIRequest(ctx, t.client, path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list tasks: %w",
			err,
		)
	}

	var tasks []models.Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return "", fmt.Errorf(
			"failed to parse tasks response: %w",
			err,
		)
	}

	return formatTaskArray(tasks), nil
}

// buildTaskListPath constructs the API path with query
// parameters for task listing.
func buildTaskListPath(
	args json.RawMessage,
) (string, error) {
	var params struct {
		Status   *string `json:"status"`
		TaskName *string `json:"task_name"`
		Type     *string `json:"type"`
		TaskID   *string `json:"task_id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}

	q := url.Values{}
	if params.Status != nil {
		q.Set("status", *params.Status)
	}
	if params.TaskName != nil {
		q.Set("task_name", *params.TaskName)
	}
	if params.Type != nil {
		q.Set("type", *params.Type)
	}
	if params.TaskID != nil {
		q.Set("task_id", *params.TaskID)
	}

	basePath := "/api/tasks/"
	if encoded := q.Encode(); encoded != "" {
		return basePath + "?" + encoded, nil
	}
	return basePath, nil
}
