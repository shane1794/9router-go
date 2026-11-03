package executor

import (
	"net/url"
	"strings"

	json "9router/proxy/internal/fastjson"
)

var senseNovaHosts = map[string]struct{}{
	"token.sensenova.ai": {},
	"token.sensenova.cn": {},
}

// StripSenseNovaReasoning removes response-side reasoning carriers from
// assistant history before it reaches SenseNova. SenseNova's Chat Completions
// endpoint returns reasoning, but its strict request schema does not accept
// reasoning_content, reasoning, or reasoning_details on replayed messages.
//
// The URL check is intentional: custom OpenAI-compatible nodes can point at
// SenseNova while retaining a generated provider id.
func StripSenseNovaReasoning(body []byte, rawURL string) []byte {
	if !isSenseNovaURL(rawURL) {
		return body
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil || request == nil {
		return body
	}
	messages, ok := request["messages"].([]any)
	if !ok {
		return body
	}

	changed := false
	for _, rawMessage := range messages {
		message, ok := rawMessage.(map[string]any)
		if !ok {
			continue
		}
		if role, _ := message["role"].(string); role != "assistant" {
			continue
		}
		for _, key := range []string{"reasoning_content", "reasoning", "reasoning_details"} {
			if _, exists := message[key]; exists {
				delete(message, key)
				changed = true
			}
		}
	}
	if !changed {
		return body
	}

	cleaned, err := json.Marshal(request)
	if err != nil {
		return body
	}
	return cleaned
}

func isSenseNovaURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	_, ok := senseNovaHosts[strings.ToLower(u.Hostname())]
	return ok
}
