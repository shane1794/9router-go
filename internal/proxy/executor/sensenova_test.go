package executor

import (
	json "encoding/json/v2"
	"strings"
	"testing"
)

func TestStripSenseNovaReasoning(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantChange bool
		wantFields []string
	}{
		{
			name:       "international host strips all replay carriers",
			url:        "https://token.sensenova.ai/v1/chat/completions",
			wantChange: true,
		},
		{
			name:       "mainland host strips all replay carriers",
			url:        "https://token.sensenova.cn/v1/chat/completions",
			wantChange: true,
		},
		{
			name:       "other host is untouched",
			url:        "https://api.example.com/v1/chat/completions",
			wantFields: []string{"reasoning_content", "reasoning", "reasoning_details"},
		},
	}

	input := []byte(`{"model":"sensenova-6.8-flash-lite","messages":[` +
		`{"role":"assistant","content":"answer","reasoning_content":"trace","reasoning":"trace 2","reasoning_details":[{"type":"summary"}]},` +
		`{"role":"user","content":"continue","reasoning_content":"keep user untouched"}]}`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripSenseNovaReasoning(input, tt.url)
			if tt.wantChange && string(got) == string(input) {
				t.Fatal("expected SenseNova reasoning carriers to be removed")
			}
			var request map[string]any
			if err := json.Unmarshal(got, &request); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			messages := request["messages"].([]any)
			assistant := messages[0].(map[string]any)
			for _, key := range tt.wantFields {
				if _, exists := assistant[key]; !exists {
					t.Errorf("other host lost assistant field %q", key)
				}
			}
			for _, key := range []string{"reasoning_content", "reasoning", "reasoning_details"} {
				if tt.wantChange {
					if _, exists := assistant[key]; exists {
						t.Errorf("SenseNova assistant retained %q", key)
					}
				}
			}
			user := messages[1].(map[string]any)
			if user["reasoning_content"] != "keep user untouched" {
				t.Errorf("user reasoning_content changed: %v", user["reasoning_content"])
			}
			if tt.wantChange && !strings.Contains(string(got), `"content":"answer"`) {
				t.Error("assistant content was not preserved")
			}
		})
	}
}
