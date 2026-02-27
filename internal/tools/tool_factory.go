package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// getTool is a data-driven GET-by-ID tool.
type getTool[T any] struct {
	client  *client.Client
	desc    string
	schema  map[string]interface{}
	pathFmt string
	format  func(*T) string
}

func (t *getTool[T]) Description() string {
	return t.desc
}

func (t *getTool[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *getTool[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	result, _, err := fetchByID[T](
		ctx,
		t.client,
		args,
		t.pathFmt,
	)
	if err != nil {
		return "", err
	}

	return t.format(result), nil
}

// getToolWithID is a GET-by-ID tool that forwards the ID to the
// formatter.
type getToolWithID[T any] struct {
	client  *client.Client
	desc    string
	schema  map[string]interface{}
	pathFmt string
	format  func(int, *T) string
}

func (t *getToolWithID[T]) Description() string {
	return t.desc
}

func (t *getToolWithID[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *getToolWithID[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	result, id, err := fetchByID[T](
		ctx,
		t.client,
		args,
		t.pathFmt,
	)
	if err != nil {
		return "", err
	}

	return t.format(id, result), nil
}

// listTool is a data-driven paginated list tool.
type listTool[T any] struct {
	client   *client.Client
	desc     string
	schema   map[string]interface{}
	basePath string
	format   func(*models.PaginatedList[T]) string
}

func (t *listTool[T]) Description() string {
	return t.desc
}

func (t *listTool[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *listTool[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	list, err := listResources[T](
		ctx,
		t.client,
		t.basePath,
		args,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list resources: %w",
			err,
		)
	}

	return t.format(list), nil
}

// patchTool is a data-driven PATCH-by-ID tool.
type patchTool[T any] struct {
	client  *client.Client
	desc    string
	schema  map[string]interface{}
	pathFmt string
	format  func(*T) string
}

func (t *patchTool[T]) Description() string {
	return t.desc
}

func (t *patchTool[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *patchTool[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	result, err := patchByID[T](
		ctx,
		t.client,
		args,
		t.pathFmt,
	)
	if err != nil {
		return "", err
	}

	return t.format(result), nil
}

// createMatchableTool is a data-driven CREATE tool for matchable
// resources.
type createMatchableTool[T any] struct {
	client *client.Client
	desc   string
	schema map[string]interface{}
	path   string
	format func(*T) string
}

func (t *createMatchableTool[T]) Description() string {
	return t.desc
}

func (t *createMatchableTool[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *createMatchableTool[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	result, err := createMatchable[T](
		ctx,
		t.client,
		args,
		t.path,
	)
	if err != nil {
		return "", err
	}

	return t.format(result), nil
}

// deleteTool is a data-driven DELETE-by-ID tool.
type deleteTool struct {
	client       *client.Client
	desc         string
	schema       map[string]interface{}
	pathFmt      string
	resourceName string
}

func (t *deleteTool) Description() string {
	return t.desc
}

func (t *deleteTool) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *deleteTool) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	return deleteByID(
		ctx,
		t.client,
		args,
		t.pathFmt,
		t.resourceName,
	)
}
