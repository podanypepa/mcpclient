// Package main implements a simple Go application.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	urlFlag := flag.String("url", "http://127.0.0.1:9996/mcp", "MCP server URL (HTTP endpoint)")
	tokenFlag := flag.String("token", "", "optional bearer token for Authorization header")
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "per-request timeout")

	flag.Parse()

	ctx := context.Background()
	client := &http.Client{
		Timeout: *timeoutFlag,
	}

	// 1) initialize
	sessionID, err := initializeSession(ctx, client, *urlFlag, *tokenFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialize error: %v\n", err)
		os.Exit(1)
	}

	// 2) notifications/initialized
	if err := sendInitializedNotification(ctx, client, *urlFlag, *tokenFlag, sessionID); err != nil {
		fmt.Fprintf(os.Stderr, "notifications/initialized error: %v\n", err)
		os.Exit(1)
	}

	// 3) tools/list
	if err := listTools(ctx, client, *urlFlag, *tokenFlag, sessionID); err != nil {
		fmt.Fprintf(os.Stderr, "tools/list error: %v\n", err)
		os.Exit(1)
	}

	// 4) tools/call now
	if err := callNowTool(ctx, client, *urlFlag, *tokenFlag, sessionID); err != nil {
		fmt.Fprintf(os.Stderr, "tools/call(now) error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
