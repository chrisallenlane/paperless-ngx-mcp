package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/models"
)

// noArgGetTool is a data-driven GET tool with no input parameters.
type noArgGetTool[T any] struct {
	client *client.Client
	desc   string
	path   string
	format func(*T) string
}

func (t *noArgGetTool[T]) Description() string {
	return t.desc
}

func (t *noArgGetTool[T]) InputSchema() map[string]interface{} {
	return emptySchema()
}

func (t *noArgGetTool[T]) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, t.path)
	if err != nil {
		return "", err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return t.format(&result), nil
}

// noArgGetToolRaw is a data-driven GET tool with no input
// parameters and custom response handling. Use this when the
// response cannot be unmarshaled into a single typed struct
// (e.g., JSON arrays, untyped maps).
type noArgGetToolRaw struct {
	client  *client.Client
	desc    string
	path    string
	process func([]byte) (string, error)
}

func (t *noArgGetToolRaw) Description() string {
	return t.desc
}

func (t *noArgGetToolRaw) InputSchema() map[string]interface{} {
	return emptySchema()
}

func (t *noArgGetToolRaw) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	body, err := doAPIRequest(ctx, t.client, t.path)
	if err != nil {
		return "", err
	}
	return t.process(body)
}

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

// createTool is a data-driven POST tool for creating resources.
type createTool[T any] struct {
	client   *client.Client
	desc     string
	schema   map[string]interface{}
	path     string
	validate func(json.RawMessage) (interface{}, error)
	format   func(*T) string
}

func (t *createTool[T]) Description() string {
	return t.desc
}

func (t *createTool[T]) InputSchema() map[string]interface{} {
	return t.schema
}

func (t *createTool[T]) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	reqBody, err := t.validate(args)
	if err != nil {
		return "", err
	}

	body, err := doPostRequest(
		ctx,
		t.client,
		t.path,
		reqBody,
	)
	if err != nil {
		return "", err
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf(
			"failed to parse response: %w",
			err,
		)
	}

	return t.format(&result), nil
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
