package mcp

import (
	"context"
	"fmt"
	"net/http"
)

// SendInitializedNotification sends a notifications/initialized notification to the MCP API.
func SendInitializedNotification(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
	sessionID string,
) error {
	req := JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}

	resp, body, err := DoMCPRequest(ctx, client, url, token, sessionID, req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifications/initialized HTTP status %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("notifications/initialized sent OK (status:", resp.Status, ")")
	return nil
}
