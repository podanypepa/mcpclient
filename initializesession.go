package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func initializeSession(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
) (string, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
				"prompts":   map[string]any{},
			},
			"clientInfo": map[string]any{
				"name":    "go-mcp-http-test-client",
				"version": "0.1.0",
			},
		},
	}

	resp, body, err := doMCPRequest(ctx, client, url, token, "", req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("initialize HTTP status %s, body: %s", resp.Status, string(body))
	}

	sessionID := resp.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		return "", fmt.Errorf("missing Mcp-Session-Id header in initialize response")
	}

	jsonBytes, err := parseSSEOrJSON(body)
	if err != nil {
		return "", fmt.Errorf("parse initialize body: %w", err)
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(jsonBytes, &rpcResp); err != nil {
		return "", fmt.Errorf("unmarshal JSON-RPC resp: %w", err)
	}

	if rpcResp.Error != nil {
		return "", fmt.Errorf("initialize JSON-RPC error: code=%d msg=%s",
			rpcResp.Error.Code, rpcResp.Error.Message)
	}

	fmt.Println("Initialize OK, session ID:", sessionID)
	return sessionID, nil
}
