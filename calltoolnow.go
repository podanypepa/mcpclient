package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func callNowTool(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
	sessionID string,
) error {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]any{
			"name": "now",
			"arguments": map[string]any{
				"format": time.RFC3339,
			},
		},
	}

	resp, body, err := doMCPRequest(ctx, client, url, token, sessionID, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tools/call(now) HTTP status %s, body: %s", resp.Status, string(body))
	}

	jsonBytes, err := parseSSEOrJSON(body)
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

	fmt.Println("Result from tool 'now':")
	for i, c := range result.Content {
		fmt.Printf("  [%d] type=%s text=%s\n", i, c.Type, c.Text)
	}

	return nil
}
