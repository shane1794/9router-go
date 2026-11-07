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
