package providers

import "strings"

// Cloudflare Workers AI has no /ai/v1/systemone route: decision models are
// served by the classic run endpoint, with the model in the path and the
// short model name in the body. The generic systemone lane swaps
// /chat/completions for /systemone, which produces
// .../ai/v1/systemone for cloudflare-ai, and Cloudflare answers that with
// HTTP 400 code 7000 ("No route for that URI").
//
// CloudflareSystemoneURL rewrites a cloudflare-ai base URL to the run
// endpoint for the given model. It returns false when the base URL is not
// the Workers AI chat lane (for example a proxy override), so the caller
// falls back to the generic systemone rules instead of guessing a path.
func CloudflareSystemoneURL(baseURL, model string) (string, bool) {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	const marker = "/ai/v1/chat/completions"
	idx := strings.Index(base, marker)
	if idx < 0 {
		return "", false
	}
	pathModel := model
	if !strings.HasPrefix(pathModel, "@cf/") && !strings.Contains(pathModel, "/") {
		pathModel = "@cf/cloudflare/" + pathModel
	}
	return base[:idx] + "/ai/run/" + pathModel, true
}

// CloudflareSystemoneModel returns the body model name Workers AI accepts
// for a decision model: "clef" or "clef-flash", not the catalog form
// "@cf/cloudflare/clef-flash".
func CloudflareSystemoneModel(model string) string {
	return strings.TrimPrefix(model, "@cf/cloudflare/")
}
