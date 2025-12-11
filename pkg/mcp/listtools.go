package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListTools lists available tools via the MCP API and prints them.
func ListTools(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
	sessionID string,
) error {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	resp, body, err := DoMCPRequest(ctx, client, url, token, sessionID, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tools/list HTTP status %s, body: %s", resp.Status, string(body))
	}

	jsonBytes, err := ParseSSEOrJSON(body)
	if err != nil {
		return fmt.Errorf("parse tools/list body: %w", err)
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(jsonBytes, &rpcResp); err != nil {
		return fmt.Errorf("unmarshal tools/list JSON-RPC: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("tools/list JSON-RPC error: code=%d msg=%s",
			rpcResp.Error.Code, rpcResp.Error.Message)
	}

	var result ListToolsResult
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return fmt.Errorf("unmarshal tools/list result: %w", err)
	}

	fmt.Println("Available tools:")
	for _, t := range result.Tools {
		fmt.Printf("  - %s: %s\n", t.Name, t.Description)
	}

	return nil
}
