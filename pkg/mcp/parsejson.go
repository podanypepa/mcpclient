package mcp

import (
	"fmt"
	"strings"
)

// ParseSSEOrJSON parses a response body that may be either raw JSON or
func ParseSSEOrJSON(raw []byte) ([]byte, error) {
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return nil, fmt.Errorf("empty response body")
	}

	if strings.HasPrefix(text, "{") {
		return []byte(text), nil
	}

	for line := range strings.SplitSeq(text, "\n") {
		line = strings.TrimSpace(line)

		if jsonPart, ok := strings.CutPrefix(line, "data:"); ok {
			jsonPart = strings.TrimSpace(jsonPart)
			if jsonPart == "" {
				continue
			}
			return []byte(jsonPart), nil
		}
	}

	return nil, fmt.Errorf("no JSON data found in SSE/response body")
}
