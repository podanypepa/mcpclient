package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func CallTool(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
	sessionID string,
	toolName string,
	toolArgs map[string]any,
) error {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params: CallToolParams{
			Name:      toolName,
			Arguments: toolArgs,
		},
	}

	resp, body, err := DoMCPRequest(ctx, client, url, token, sessionID, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tools/call(%s) HTTP status %s, body: %s", toolName, resp.Status, string(body))
	}

	jsonBytes, err := ParseSSEOrJSON(body)
	if err != nil {
		return fmt.Errorf("parse tools/call body: %w", err)
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(jsonBytes, &rpcResp); err != nil {
		return fmt.Errorf("unmarshal tools/call JSON-RPC: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("tools/call JSON-RPC error: code=%d msg=%s",
			rpcResp.Error.Code, rpcResp.Error.Message)
	}

	var result CallToolResult
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return fmt.Errorf("unmarshal tools/call result: %w", err)
	}

	fmt.Printf("Result from tool '%s':\n", toolName)
	for i, c := range result.Content {
		fmt.Printf("  [%d] type=%s text=%s\n", i, c.Type, c.Text)
	}

	return nil
}