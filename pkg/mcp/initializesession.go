package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func InitializeSession(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
) (string, error) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities: map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
				"prompts":   map[string]any{},
			},
			ClientInfo: ClientInfo{
				Name:    "mcpclient",
				Version: "0.1.0",
			},
		},
	}

	resp, body, err := DoMCPRequest(ctx, client, url, token, "", req)
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

	jsonBytes, err := ParseSSEOrJSON(body)
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
