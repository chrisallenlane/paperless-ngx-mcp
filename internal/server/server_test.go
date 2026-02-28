package server

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
)

func TestHandleInitialize(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	s := New(c)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Response JSONRPC = %s, want 2.0", resp.JSONRPC)
	}

	if resp.ID != 1 {
		t.Errorf("Response ID = %v, want 1", resp.ID)
	}

	if resp.Error != nil {
		t.Errorf("Unexpected error: %+v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("Result should not be nil")
	}

	// Verify result structure
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	if result["protocolVersion"] != MCPProtocolVersion {
		t.Errorf(
			"Protocol version = %v, want %s",
			result["protocolVersion"],
			MCPProtocolVersion,
		)
	}

	serverInfo, ok := result["serverInfo"].(map[string]string)
	if !ok {
		t.Fatal("serverInfo should be a map")
	}

	if serverInfo["name"] != ServerName {
		t.Errorf(
			"Server name = %s, want %s",
			serverInfo["name"],
			ServerName,
		)
	}

	if serverInfo["version"] != ServerVersion {
		t.Errorf(
			"Server version = %s, want %s",
			serverInfo["version"],
			ServerVersion,
		)
	}
}

func TestHandleListTools(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	s := New(c)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error != nil {
		t.Errorf("Unexpected error: %+v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("Result should not be nil")
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	tools, ok := result["tools"].([]map[string]interface{})
	if !ok {
		t.Fatal("tools should be a slice")
	}

	if len(tools) == 0 {
		t.Fatal("Expected at least one registered tool")
	}

	// Verify tool structure for any registered tools
	for _, tool := range tools {
		if _, ok := tool["name"]; !ok {
			t.Error("Tool should have a name")
		}
		if _, ok := tool["description"]; !ok {
			t.Error("Tool should have a description")
		}
		if _, ok := tool["inputSchema"]; !ok {
			t.Error("Tool should have an inputSchema")
		}
	}
}

func TestHandleUnknownMethod(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	s := New(c)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "unknown/method",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for unknown method")
	}

	if resp.Error.Code != -32601 {
		t.Errorf("Error code = %d, want -32601", resp.Error.Code)
	}

	if !strings.Contains(resp.Error.Message, "Method not found") {
		t.Errorf(
			"Error message should mention 'Method not found', got: %s",
			resp.Error.Message,
		)
	}

	if resp.Result != nil {
		t.Error("Result should be nil for error response")
	}
}

func TestHandleCallTool_InvalidTool(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	s := New(c)

	params := map[string]interface{}{
		"name":      "nonexistent_tool",
		"arguments": json.RawMessage(`{}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      4,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for nonexistent tool")
	}

	if resp.Error.Code != -32603 {
		t.Errorf("Error code = %d, want -32603", resp.Error.Code)
	}
}

func TestHandleCallTool_MalformedParams(t *testing.T) {
	c := client.New("http://localhost", "test-token")
	s := New(c)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      5,
		Method:  "tools/call",
		Params:  json.RawMessage(`{invalid json}`),
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for malformed params")
	}

	if resp.Error.Code != -32603 {
		t.Errorf("Error code = %d, want -32603", resp.Error.Code)
	}

	if !strings.Contains(
		resp.Error.Message,
		"failed to parse tool call params",
	) {
		t.Errorf(
			"Error message should mention parsing failure, got: %s",
			resp.Error.Message,
		)
	}
}

func TestJSONRPCRequest_Unmarshal(t *testing.T) {
	jsonData := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`

	var req JSONRPCRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %s, want 2.0", req.JSONRPC)
	}

	if req.Method != "initialize" {
		t.Errorf("Method = %s, want initialize", req.Method)
	}
}

// runServer is a helper that writes lines to a server's stdin, runs the server,
// and returns the decoded responses from stdout.
func runServer(
	t *testing.T,
	lines []string,
) []map[string]interface{} {
	t.Helper()

	c := client.New("http://localhost", "test-token")
	s := New(c)

	var stdin bytes.Buffer
	for _, line := range lines {
		stdin.WriteString(line + "\n")
	}

	var stdout bytes.Buffer
	if err := s.Run(context.Background(), &stdin, &stdout); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	// Decode one JSON object per output line.
	var responses []map[string]interface{}
	decoder := json.NewDecoder(&stdout)
	for decoder.More() {
		var resp map[string]interface{}
		if err := decoder.Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		responses = append(responses, resp)
	}

	return responses
}

func TestRun_WellFormedRequest(t *testing.T) {
	line := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`
	responses := runServer(t, []string{line})

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, got %d", len(responses))
	}

	resp := responses[0]

	if resp["jsonrpc"] != "2.0" {
		t.Errorf("jsonrpc = %v, want 2.0", resp["jsonrpc"])
	}

	// id is decoded as float64 from JSON.
	if resp["id"].(float64) != 1 {
		t.Errorf("id = %v, want 1", resp["id"])
	}

	if resp["error"] != nil {
		t.Errorf("Unexpected error: %v", resp["error"])
	}

	if resp["result"] == nil {
		t.Error("Result should not be nil")
	}
}

func TestRun_MalformedJSON(t *testing.T) {
	line := `{not valid json`
	responses := runServer(t, []string{line})

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, got %d", len(responses))
	}

	resp := responses[0]

	errObj, ok := resp["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected an error object in response")
	}

	if errObj["code"].(float64) != -32700 {
		t.Errorf("Error code = %v, want -32700", errObj["code"])
	}
}

func TestRun_MultipleRequests(t *testing.T) {
	lines := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	}
	responses := runServer(t, lines)

	if len(responses) != 2 {
		t.Fatalf("Expected 2 responses, got %d", len(responses))
	}

	for i, resp := range responses {
		if resp["jsonrpc"] != "2.0" {
			t.Errorf(
				"response[%d] jsonrpc = %v, want 2.0",
				i,
				resp["jsonrpc"],
			)
		}
		if resp["error"] != nil {
			t.Errorf("response[%d] unexpected error: %v", i, resp["error"])
		}
		if resp["result"] == nil {
			t.Errorf("response[%d] result should not be nil", i)
		}
	}

	// Verify IDs match their respective requests.
	if responses[0]["id"].(float64) != 1 {
		t.Errorf("response[0] id = %v, want 1", responses[0]["id"])
	}

	if responses[1]["id"].(float64) != 2 {
		t.Errorf("response[1] id = %v, want 2", responses[1]["id"])
	}
}

func TestJSONRPCError_Marshal(t *testing.T) {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Error: &JSONRPCError{
			Code:    -32600,
			Message: "Invalid Request",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	errorObj, ok := decoded["error"].(map[string]interface{})
	if !ok {
		t.Fatal("error should be an object")
	}

	if errorObj["code"].(float64) != -32600 {
		t.Errorf(
			"error code = %v, want -32600",
			errorObj["code"],
		)
	}

	if errorObj["message"] != "Invalid Request" {
		t.Errorf(
			"error message = %v, want Invalid Request",
			errorObj["message"],
		)
	}
}
