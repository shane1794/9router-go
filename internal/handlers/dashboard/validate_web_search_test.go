package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"9router/proxy/internal/providers"
)

// The Exa regression (2026-11-07): the dashboard Test/Validate probe treated
// Exa as an Anthropic provider because isAnthropicProbe keys on AuthHeader
// "x-api-key", and POSTed an Anthropic messages body to /search, which Exa
// cannot accept. A web-only provider must be probed with {"query": "ping"}
// and its own auth header instead (upstream probeWebProvider parity).
func TestValidateProviderKey_Exa_WebProbe(t *testing.T) {
	cfg := providers.KnownProviders["exa"]
	calls := withProbeStub(t, func(call int, p capturedProbe) (int, []byte, error) {
		return http.StatusOK, nil, nil
	})
	out := validateProviderKey(context.Background(), "exa", cfg, "test-key", nil)
	if !out.valid || !out.supported {
		t.Fatalf("valid=%v supported=%v, want true/true", out.valid, out.supported)
	}
	if len(*calls) != 1 {
		t.Fatalf("probe calls = %d, want 1", len(*calls))
	}
	call := (*calls)[0]
	if call.method != http.MethodPost || call.url != "https://api.exa.ai/search" {
		t.Errorf("probe = %s %s, want POST https://api.exa.ai/search", call.method, call.url)
	}
	if got := call.headers.Get("x-api-key"); got != "test-key" {
		t.Errorf("x-api-key header = %q, want %q", got, "test-key")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(call.body), &body); err != nil {
		t.Fatalf("probe body is not JSON: %v", err)
	}
	if body["query"] != "ping" {
		t.Errorf("probe body query = %v, want ping", body["query"])
	}
	if _, ok := body["messages"]; ok {
		t.Errorf("probe body carries an Anthropic messages field: %s", call.body)
	}
	if _, ok := body["model"]; ok {
		t.Errorf("probe body carries a model field: %s", call.body)
	}
}

func TestValidateProviderKey_Exa_Unauthorized(t *testing.T) {
	cfg := providers.KnownProviders["exa"]
	staticProbeStub(t, http.StatusUnauthorized)
	out := validateProviderKey(context.Background(), "exa", cfg, "bad-key", nil)
	if out.valid {
		t.Errorf("valid=true for a 401 probe, want false")
	}
	if !out.supported {
		t.Errorf("supported=false, want true")
	}
}

func TestIsWebOnlyProvider(t *testing.T) {
	cases := map[string]bool{
		"exa": true, "tavily": true, "serper": true, "linkup": true,
		"tinyfish": true, "searchapi": true, "brave-search": true,
		"openai": false, "anthropic": false, "ollama": false,
		"cloudflare-ai": false, "unknown-provider": false,
	}
	for provider, want := range cases {
		if got := isWebOnlyProvider(provider); got != want {
			t.Errorf("isWebOnlyProvider(%q) = %v, want %v", provider, got, want)
		}
	}
}
