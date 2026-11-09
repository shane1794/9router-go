package media

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"9router/proxy/internal/db"
	"9router/proxy/internal/providers"
)

// The Exa regression (2026-11-07): the direct /v1/search lane appended the
// gateway route to the provider BaseURL (api.exa.ai/search/v1/search) instead
// of using SearchURL, and /v1/web/fetch forwarded {"url"} while Exa /contents
// needs {"urls": [...]}. The combo lane already had the SearchURL override;
// the direct lane did not.
func TestHandleSearch_Exa_UsesSearchURL(t *testing.T) {
	var gotPath, gotMethod, gotKey, gotBody string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotKey = r.Header.Get("x-api-key")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer upstream.Close()

	orig := providers.KnownProviders["exa"]
	exa := orig
	exa.SearchURL = upstream.URL + "/search"
	exa.FetchURL = upstream.URL + "/contents"
	providers.KnownProviders["exa"] = exa
	defer func() { providers.KnownProviders["exa"] = orig }()

	sqlDB, cleanup := setupEmbeddingsTestDB(t)
	defer cleanup()
	repo := db.NewRepo(sqlDB)
	p := 1
	if err := repo.CreateProviderConnectionFull("conn-exa-1", "exa", "apikey", "Exa", &p, `{"apiKey":"exa-key-123"}`); err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	handler := newTestMediaHandler(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/search",
		strings.NewReader(`{"model":"exa/search","query":"latest AI news"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.HandleSearch(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	if gotPath != "/search" {
		t.Errorf("upstream path = %q, want %q (the /v1/search suffix must not be appended)", gotPath, "/search")
	}
	if gotMethod != http.MethodPost {
		t.Errorf("upstream method = %q, want POST", gotMethod)
	}
	if gotKey != "exa-key-123" {
		t.Errorf("x-api-key header = %q, want %q", gotKey, "exa-key-123")
	}
	if !strings.Contains(gotBody, `"query":"latest AI news"`) {
		t.Errorf("upstream body lost the query: %s", gotBody)
	}
}

func TestHandleWebFetch_Exa_TranslatesURLToURLs(t *testing.T) {
	var gotPath, gotBody string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer upstream.Close()

	orig := providers.KnownProviders["exa"]
	exa := orig
	exa.SearchURL = upstream.URL + "/search"
	exa.FetchURL = upstream.URL + "/contents"
	providers.KnownProviders["exa"] = exa
	defer func() { providers.KnownProviders["exa"] = orig }()

	sqlDB, cleanup := setupEmbeddingsTestDB(t)
	defer cleanup()
	repo := db.NewRepo(sqlDB)
	p := 1
	if err := repo.CreateProviderConnectionFull("conn-exa-2", "exa", "apikey", "Exa", &p, `{"apiKey":"exa-key-123"}`); err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	handler := newTestMediaHandler(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/web/fetch",
		strings.NewReader(`{"url":"https://example.com","model":"exa/fetch"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.HandleWebFetch(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	if gotPath != "/contents" {
		t.Errorf("upstream path = %q, want %q", gotPath, "/contents")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(gotBody), &body); err != nil {
		t.Fatalf("upstream body is not JSON: %v", err)
	}
	if _, ok := body["url"]; ok {
		t.Errorf("upstream body still carries singular \"url\": %s", gotBody)
	}
	urls, ok := body["urls"].([]any)
	if !ok || len(urls) != 1 || urls[0] != "https://example.com" {
		t.Errorf("upstream body urls = %v, want [https://example.com]: %s", body["urls"], gotBody)
	}
}

// Upstream parity (2026-11-09): the direct /v1/search lane forwarded the
// dashboard body as-is — {"model": "", "query", "search_type",
// "max_results"} plus the dashboard's Authorization: Bearer <gateway key>
// header — while upstream's buildExaRequest sends a rebuilt body and exactly
// Content-Type + x-api-key. Exa kept rejecting the forwarded shape, so the
// lane now rebuilds the body and strips the leaked Authorization header.
func TestHandleSearch_Exa_UpstreamParity(t *testing.T) {
	var gotAuth, gotKey, gotBody string
	var hits int
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		gotAuth = r.Header.Get("Authorization")
		gotKey = r.Header.Get("x-api-key")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer upstream.Close()

	orig := providers.KnownProviders["exa"]
	exa := orig
	exa.SearchURL = upstream.URL + "/search"
	exa.FetchURL = upstream.URL + "/contents"
	providers.KnownProviders["exa"] = exa
	defer func() { providers.KnownProviders["exa"] = orig }()

	sqlDB, cleanup := setupEmbeddingsTestDB(t)
	defer cleanup()
	repo := db.NewRepo(sqlDB)
	p := 1
	if err := repo.CreateProviderConnectionFull("conn-exa-3", "exa", "apikey", "Exa", &p, `{"apiKey":"exa-key-123"}`); err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	handler := newTestMediaHandler(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/search",
		strings.NewReader(`{"model":"exa","query":"latest AI news","search_type":"web","max_results":5}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-gateway-key")
	handler.HandleSearch(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	if hits != 1 {
		t.Fatalf("upstream hits = %d, want 1", hits)
	}
	if gotAuth != "" {
		t.Errorf("upstream Authorization header = %q, want empty (the gateway client key must not leak to Exa)", gotAuth)
	}
	if gotKey != "exa-key-123" {
		t.Errorf("x-api-key header = %q, want %q", gotKey, "exa-key-123")
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(gotBody), &body); err != nil {
		t.Fatalf("upstream body is not JSON: %v", err)
	}
	if body["query"] != "latest AI news" {
		t.Errorf("upstream body query = %v, want %q", body["query"], "latest AI news")
	}
	if body["numResults"] != float64(5) {
		t.Errorf("upstream body numResults = %v, want 5", body["numResults"])
	}
	if body["type"] != "auto" || body["text"] != true || body["highlights"] != true {
		t.Errorf("upstream body missing Exa fields (type auto, text/highlights true): %s", gotBody)
	}
	for _, k := range []string{"model", "search_type", "max_results"} {
		if _, ok := body[k]; ok {
			t.Errorf("upstream body still carries gateway field %q: %s", k, gotBody)
		}
	}
}

// A key-less Exa connection must fail fast at the gateway (401 "no API key
// found") and never reach Exa — upstream's tryDedicatedProvider parity.
func TestHandleSearch_Exa_KeylessFailsFast(t *testing.T) {
	var hits int
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer upstream.Close()

	orig := providers.KnownProviders["exa"]
	exa := orig
	exa.SearchURL = upstream.URL + "/search"
	exa.FetchURL = upstream.URL + "/contents"
	providers.KnownProviders["exa"] = exa
	defer func() { providers.KnownProviders["exa"] = orig }()

	sqlDB, cleanup := setupEmbeddingsTestDB(t)
	defer cleanup()
	repo := db.NewRepo(sqlDB)
	p := 1
	if err := repo.CreateProviderConnectionFull("conn-exa-4", "exa", "apikey", "Exa", &p, `{}`); err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	handler := newTestMediaHandler(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/search",
		strings.NewReader(`{"model":"exa","query":"latest AI news"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	handler.HandleSearch(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want 401, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "no API key found") {
		t.Errorf("gateway error does not name the missing key: %s", rec.Body.String())
	}
	if hits != 0 {
		t.Errorf("upstream hits = %d, want 0 (a key-less request must never reach Exa)", hits)
	}
}

func TestBuildExaSearchBody(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"dashboard body becomes Exa body", `{"model":"exa","query":"latest AI news","search_type":"web","max_results":5}`, `{"query":"latest AI news","numResults":5,"type":"auto","text":true,"highlights":true}`},
		{"missing max_results defaults to 5", `{"model":"exa","query":"q"}`, `{"query":"q","numResults":5,"type":"auto","text":true,"highlights":true}`},
		{"max_results capped at 100", `{"query":"q","max_results":500}`, `{"query":"q","numResults":100,"type":"auto","text":true,"highlights":true}`},
		{"news adds category", `{"query":"q","search_type":"news","max_results":3}`, `{"query":"q","numResults":3,"type":"auto","text":true,"highlights":true,"category":"news"}`},
		{"domain filter splits includes and excludes", `{"query":"q","domain_filter":["example.com","-spam.com"]}`, `{"query":"q","numResults":5,"type":"auto","text":true,"highlights":true,"includeDomains":["example.com"],"excludeDomains":["spam.com"]}`},
		{"no query passes through", `{"model":"exa","max_results":5}`, `{"model":"exa","max_results":5}`},
		{"not JSON passes through", `not json`, `not json`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(buildExaSearchBody([]byte(tc.in)))
			var gotJSON, wantJSON any
			_ = json.Unmarshal([]byte(got), &gotJSON)
			_ = json.Unmarshal([]byte(tc.want), &wantJSON)
			if gotJSON != nil && wantJSON != nil {
				gotB, _ := json.Marshal(gotJSON)
				wantB, _ := json.Marshal(wantJSON)
				if string(gotB) != string(wantB) {
					t.Errorf("got %s, want %s", got, tc.want)
				}
				return
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTranslateFetchURLsBody(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"singular url becomes urls array", `{"url":"https://example.com"}`, `{"urls":["https://example.com"]}`},
		{"other fields preserved", `{"url":"https://example.com","format":"markdown"}`, `{"format":"markdown","urls":["https://example.com"]}`},
		{"urls already present passes through", `{"urls":["https://example.com"]}`, `{"urls":["https://example.com"]}`},
		{"no url passes through", `{"query":"ping"}`, `{"query":"ping"}`},
		{"not JSON passes through", `not json`, `not json`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(translateFetchURLsBody([]byte(tc.in)))
			var gotJSON, wantJSON any
			_ = json.Unmarshal([]byte(got), &gotJSON)
			_ = json.Unmarshal([]byte(tc.want), &wantJSON)
			if gotJSON != nil && wantJSON != nil {
				gotB, _ := json.Marshal(gotJSON)
				wantB, _ := json.Marshal(wantJSON)
				if string(gotB) != string(wantB) {
					t.Errorf("got %s, want %s", got, tc.want)
				}
				return
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
