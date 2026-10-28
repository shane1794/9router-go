package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"9router/proxy/internal/providers"
)

// newNousGrantServer starts a token endpoint fixture: it asserts the request
// the bug report was about (the spent grant in BOTH the form body and the
// x-nous-refresh-token header) and answers with the given status/body.
func newNousGrantServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-nous-refresh-token"); got != "spent-token" {
			t.Errorf("x-nous-refresh-token = %q, want spent-token", got)
		}
		if got := r.FormValue("refresh_token"); got != "spent-token" {
			t.Errorf("refresh_token = %q, want spent-token", got)
		}
		if got := r.FormValue("grant_type"); got != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", got)
		}
		if got := r.FormValue("client_id"); got != "hermes-cli" {
			t.Errorf("client_id = %q, want hermes-cli", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// pinNousGrantServer points the refresher (and its agent-key mint) at the
// fixture for the life of the test, then restores production URLs.
func pinNousGrantServer(t *testing.T, srv *httptest.Server) {
	t.Helper()
	prevToken, prevAgent := nousTokenURL, nousAgentKeyURL
	nousTokenURL = srv.URL
	nousAgentKeyURL = srv.URL
	t.Cleanup(func() {
		nousTokenURL, nousAgentKeyURL = prevToken, prevAgent
	})
}

// A spent single-use refresh token must arrive as a typed error the parkers
// recognise: unwrapped fmt text is invisible to IsRefreshGrantDead, so the
// request path and the background tick would spend the dead grant forever.
func TestRefreshNous_SpentTokenIsGrantDead(t *testing.T) {
	server := newNousGrantServer(t, http.StatusBadRequest,
		`{"error":"invalid_grant","error_description":"This refresh token was already rejected; please re-authenticate"}`)
	pinNousGrantServer(t, server)

	_, err := refreshNous(context.Background(),
		&Params{Client: server.Client(), Provider: "nous", RefreshToken: "spent-token"})
	if err == nil {
		t.Fatal("a spent refresh token returned no error")
	}
	if !providers.IsRefreshGrantDead(err) {
		t.Errorf("IsRefreshGrantDead(%v) = false, want true", err)
	}
}

// A blip must stay retryable: only a dead grant may park the account.
func TestRefreshNous_TransientFailureIsNotGrantDead(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"token endpoint down", http.StatusBadGateway, `boom`},
		{"rate limited", http.StatusTooManyRequests, `{"error":"slow_down"}`},
		{"server error with a body", http.StatusInternalServerError, `{"error":"internal"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newNousGrantServer(t, tt.status, tt.body)
			pinNousGrantServer(t, server)

			_, err := refreshNous(context.Background(),
				&Params{Client: server.Client(), Provider: "nous", RefreshToken: "spent-token"})
			if err == nil {
				t.Fatalf("a %d answer returned no error", tt.status)
			}
			if providers.IsRefreshGrantDead(err) {
				t.Errorf("IsRefreshGrantDead(%v) = true, want false", err)
			}
		})
	}
}
