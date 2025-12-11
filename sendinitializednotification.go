package main

import (
	"context"
	"fmt"
	"net/http"
)

func sendInitializedNotification(
	ctx context.Context,
	client *http.Client,
	url string,
	token string,
	sessionID string,
) error {
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}

	resp, body, err := doMCPRequest(ctx, client, url, token, sessionID, req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifications/initialized HTTP status %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("notifications/initialized sent OK (status:", resp.Status, ")")
	return nil
}
