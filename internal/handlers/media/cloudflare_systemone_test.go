package media

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"9router/proxy/internal/db"
)

// Cloudflare has no /ai/v1/systemone route: the lane must send a
// cloudflare-ai decision model to the classic run endpoint with the short
// model name in the body, or Workers AI answers HTTP 400 code 7000
// ("No route for that URI").
func TestHandleSystemone_Cloudflare_Mock(t *testing.T) {
	var gotPath, gotModel string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var body struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotModel = body.Model
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"model":"clef-flash","answers":{"is_urgent":{"type":"noul","noul":0.04}}}`))
	}))
	defer server.Close()

	sqlDB, cleanup := setupEmbeddingsTestDB(t)
	defer cleanup()
	repo := db.NewRepo(sqlDB)
	p := 1
	err := repo.CreateProviderConnectionFull("conn-cf-1", "cloudflare-ai", "apikey", "Cloudflare", &p, `{"baseUrl":"`+server.URL+`/client/v4/accounts/test123/ai/v1/chat/completions","apiKey":"test-key"}`)
	if err != nil {
		t.Fatalf("failed to create connection: %v", err)
	}

	handler := newTestMediaHandler(repo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/systemone",
		strings.NewReader(`{"model":"cf/@cf/cloudflare/clef-flash","state":"The server is responding normally.","questions":{"is_urgent":{"type":"noul","instructions":"Urgent?"}}}`),
	)
	handler.HandleSystemone(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200, body: %s", rec.Code, rec.Body.String())
	}
	if want := "/client/v4/accounts/test123/ai/run/@cf/cloudflare/clef-flash"; gotPath != want {
		t.Errorf("upstream path = %q, want %q", gotPath, want)
	}
	if gotModel != "clef-flash" {
		t.Errorf("upstream body model = %q, want %q", gotModel, "clef-flash")
	}
}
