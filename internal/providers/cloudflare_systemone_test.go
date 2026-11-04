package providers

import "testing"

func TestCloudflareSystemoneURL(t *testing.T) {
	const base = "https://api.cloudflare.com/client/v4/accounts/acct123/ai/v1/chat/completions"
	cases := []struct {
		name, baseURL, model, want string
		wantOK                     bool
	}{
		{"catalog model", base, "@cf/cloudflare/clef-flash", "https://api.cloudflare.com/client/v4/accounts/acct123/ai/run/@cf/cloudflare/clef-flash", true},
		{"catalog clef", base, "@cf/cloudflare/clef", "https://api.cloudflare.com/client/v4/accounts/acct123/ai/run/@cf/cloudflare/clef", true},
		{"bare model", base, "clef-flash", "https://api.cloudflare.com/client/v4/accounts/acct123/ai/run/@cf/cloudflare/clef-flash", true},
		{"trailing slash base", base + "/", "@cf/cloudflare/clef-flash", "https://api.cloudflare.com/client/v4/accounts/acct123/ai/run/@cf/cloudflare/clef-flash", true},
		{"non workers-ai base falls back", "https://proxy.example/v1/chat/completions", "@cf/cloudflare/clef-flash", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := CloudflareSystemoneURL(tc.baseURL, tc.model)
			if ok != tc.wantOK {
				t.Fatalf("CloudflareSystemoneURL() ok = %v, want %v", ok, tc.wantOK)
			}
			if got != tc.want {
				t.Errorf("CloudflareSystemoneURL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCloudflareSystemoneModel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"@cf/cloudflare/clef-flash", "clef-flash"},
		{"@cf/cloudflare/clef", "clef"},
		{"clef-flash", "clef-flash"},
	}
	for _, tc := range cases {
		if got := CloudflareSystemoneModel(tc.in); got != tc.want {
			t.Errorf("CloudflareSystemoneModel(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
