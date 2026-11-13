# Changelog (fork `shane1794/9router-go`)

Fork-local history. Every change landing in this fork gets an entry under
`## [Unreleased]`, newest first — in addition to the `.changes/` fragment
required by the upstream convention (see `AGENTS.md` step 6). This file never
merges into upstream `CHANGELOG.md`.

## [Unreleased]

### Upstream sync to 21d1fc5a (2026-11-13)

- Synced upstream luqman-v1/9router-go through `21d1fc5a` (Kenari provider #1155, Go toolchain 1.27.4 #1153, version bump to 1.9.17-exp.1, SOCKS5 #1120, prompt-cache TTL #1121, xiaomi-tokenplan #1119, antigravity revert #1118, kiro #1117, free-pi revert #1116). Trial merge had 44 conflicts, all resolved. Fork behaviour preserved: CodeBuddy fresh request/conversation ids, fork-only `make build` (only `9router-go`), model context tags, Nous catalogue (10 models), Jev tilde alias, SenseNova `.ai` pointing + 4-model list, clef systemone routing, Exa probe/fetch/search translation, media-page Add Connection `apiKey` fix, and the sensenova registry line.
- Codex responses-lite now follows upstream's contract-based selection (`shapeCodexResponsesBody` + the `x-openai-internal-codex-responses-lite` header) instead of the fork's model-list `applyCodexModelShape`, which upstream removed; the web-search/Claude translation helpers (`translateCodexClaudeWebSearch`, `translateCodexClaudeWebSearchChoice`, `registerCodexHostedWebSearch`) now come from upstream unchanged, so `transform.go` no longer needs the fork's `applyCodexModelShape` call. This keeps the codex lane consistent with upstream's `providers.go`, handlers, and tests.
- Deleted fork test `internal/handlers/media/cloudflare_flux2_test.go`: it referenced `isCloudflareFlux2Model` and `flux2ReferenceImage`, which are not defined anywhere in the repo (not even at fork HEAD), so it never compiled.
- Fixed auto-merge artifacts: duplicate `TestConnectionRequiresReauth` in `internal/db/accounts_oauth_lock_test.go`; missing closing brace in `TestCohereChatUsesCompatibilityEndpoint` (`internal/providers/providers_test.go`); duplicated `httptest` server block in `internal/handlers/chat/compatible_base_url_test.go` and a `gotPath` declaration lost with it.
- VERSION is now 1.9.17-exp.1. Verified: `make build` passes (1.9.17-exp.1 embedded), `go build ./...` and `go vet ./...` are clean, and all test packages compile (`go test -run '^$' ./...`). No provider connections were checked or tested live.

### Docs: README install guide is build-from-source only (2026-11-11)

- The README's install guide offered the upstream one-line installers (`install.sh`, `install.ps1`, `install-noavx.sh`), the `luqmenul/9router-go` Docker image, manual release downloads, and the experimental pre-release builds — every one of them installs upstream's build, not this fork. The guide now documents only build from source: clone `shane1794/9router-go`, `make web-build`, `make build`. The no-AVX note is reframed as a build flag (`CGO_ENABLED=0 GOAMD64=v1 make build`) instead of a release asset, and the "Release binary" and "Experimental builds" sections are removed. Fragment `.changes/949-readme-install-build-from-source-only.md`.
- Verified: grep confirms no install path in README still points at upstream (remaining `luqman-v1` links are header badges and docs/credits links, not install steps); both `git clone` lines point at the fork; the edited sections read back cleanly. No code changed, so no build or test run was needed for this docs-only change.

### Fix: Cohere chat now uses the Compatibility API endpoint (2026-11-10)

- Chat requests to Cohere went to `https://api.cohere.ai/v1/chat/completions` and Cohere answered 405: the native API has no OpenAI-shaped route there (its Chat API is `/v1/chat` with `{message, chat_history}` in v1 and `/v2/chat` with a `messages` array in v2). The registry now points cohere at Cohere's own Compatibility API, `https://api.cohere.ai/compatibility/v1/chat/completions` (docs.cohere.com/docs/compatibility-api), which serves the OpenAI Chat Completions shape — including streaming, tools, and structured outputs — with the same Bearer key, so the generic OpenAI forward needs no translation. Upstream (luqman-v1/9router-go at 298fee33, and serenhope's cohere.js) still point at the dead `/v1` path, so this is fork-only until upstream moves. Fragment `.changes/948-cohere-compatibility-endpoint.md`.
- The sibling lanes follow from the same base URL: embeddings now resolve to `compatibility/v1/embeddings` (documented on the same page), and the models probe to `compatibility/v1/models` (verified to exist: it answers 401 without a key, same as the old `/v1/models`). Rerank is the exception — the Compatibility API has no rerank lane (OpenAI has none), so `buildRerankURL` resolves a compatibility base on the native `/v1/rerank` instead of inventing `/compatibility/v1/rerank`.
- Verified: new tests `TestCohereChatUsesCompatibilityEndpoint`, `TestHandleChatCompletions_CohereCompatibilityEndpoint` (mock upstream receives POST /compatibility/v1/chat/completions), plus compatibility cases in `TestBuildRerankURL`, `TestBuildEmbeddingsURL`, and `TestCompatibleNodeModelsURLFollowsTheChatRoute`, pass; `make build` passes (1.9.16-exp.10). No provider connections were checked or tested live.

### Upstream sync to 298fee33 (2026-11-10)

- Synced upstream luqman-v1/9router-go through `298fee33` (29 commits after `c74841e7`: #1050 ImgBot image optimisation, #1049 free-pi provider with OAuth device flow, #1047 Mistral Vibe Code usage + connection routing + expanded catalog, #1044 headroom engine over MCP + conversation prefix replay, #1042 provider-agnostic upstream error classification, #1041 dashboard operator budget read, #1048 opencode abuse-flagged egress retry, #1046 kiro default profile ARN, #1043 gzip identity streaming, #1040 translator response_format onto Gemini native output, #1038 codex-image gpt-5.6-terra fallback, #1030 Claude quota reset catch, #1028 combo Claude redacted thinking, plus version bumps to 1.9.16-exp.9 then 1.9.16-exp.10). Trial merge had ZERO conflicts. Fork behaviour preserved: CodeBuddy fresh request/conversation ids, fork-only `make build` (only `9router-go`), model context tags, Nous catalogue (10 models incl. the 3 serenhope seeds), Jev tilde alias, SenseNova `.ai` pointing + 4-model list, clef systemone routing, Exa probe + fetch translation + search body rebuild (fragment 947), media-page Add Connection `apiKey` fix (fragment 946). Upstream's own `CHANGELOG.md` still carries committed conflict markers from an earlier upstream merge (`1b5688b`) — identical on both sides, left untouched. VERSION is now 1.9.16-exp.10. Verified: `make build` passes (only `9router-go`, version 1.9.16-exp.10 embedded), `go test ./internal/providers/` and `./internal/handlers/dashboard/` pass, web unit 839 pass / 1 fail (the known timezone-only `consoleLog` failure; `consoleLog` untouched upstream this window), svelte-check 74 errors == baseline 74 (warnings 66 vs 60, from upstream's new Svelte files; the ratchet pins errors only). Known pre-existing failures, unchanged: executor `TestForwardQoder_UnknownModelNamesTheModel` (VM egress proxy TLS); media `TestHandleAudioVoices_elevenlabs` and `TestTargetPrivate_RespectsAllowPrivateNetworking`; chat `TestLocalGatewayConnectionPermittedWithOptIn` (fails on pristine upstream too — the VM egress proxy resolves api.deepseek.com to 198.18.25.163, which the #840 SSRF guard blocks). Known gaps: 10 upstream PNGs not pushed (binary — ImgBot re-optimised 7 docs screenshots + 2 evidence images + the new free-pi provider icon this window); `.changes/1009-banned-tab-appeal.md` deleted upstream but stays on the fork (snapshot cannot delete); `ci.yml`/`release.yml` still pre-`c74841e7` (GitHub App cannot write `.github/workflows/`, only Arthur's own git push can land them). No provider connections were checked or tested.

### Fix: Exa /v1/search requests now match upstream's shape (2026-11-09)

- The direct and combo /v1/search lanes forwarded the dashboard body as-is — `{"model": "", "query", "search_type", "max_results"}` — and copied the dashboard's `Authorization: Bearer <gateway client key>` header onto the Exa request alongside `x-api-key`. Upstream's `buildExaRequest` sends a rebuilt body `{query, numResults, type: "auto", text: true, highlights: true}` with exactly `Content-Type` + `x-api-key`, and Exa kept rejecting the forwarded shape. Both lanes now rebuild the Exa search body (`max_results` defaults to 5, capped at 100; `search_type: "news"` adds `category: "news"`; `domain_filter` splits into `includeDomains`/`excludeDomains`), and the `Authorization` header is stripped whenever the provider authenticates on another header. The direct lane's fail-fast (gateway 401 "no API key found" on an empty stored key) is now covered by a test; the response is still passed through as-is. Verified: new tests `TestHandleSearch_Exa_UpstreamParity`, `TestHandleSearch_Exa_KeylessFailsFast`, `TestBuildExaSearchBody` pass; `make build` passes (1.9.16-exp.8). No provider connections were checked or tested live.

### Upstream sync to c74841e7 (2026-11-09)

- Synced upstream luqman-v1/9router-go through `c74841e7` (59 commits after `5e496d12`: #1009 banned tab + appeal URL + provider risk notice, #1010 SenseNova models port + Gemini remote-image MIME, #1006 minimax-code dead-grant park, #1003 custom-model service kind, #1002 TTS audio output modality, #1001 Gemini Omni thinking level, #999 TTFT + ResponseBuf stamping, #998 OpenAI output floor, #968 SSO OIDC/SAML login flows, #979 built-in headroom engine, #982 custom JSON API endpoints, #984 live availability + cooldown, #971 Mistral catalog + live fetcher, #990 Groq strict-message parity, #987 Groq reasoning_format pin, #986 deepseek-web thinking envelope, #985 cloudflare-ai image models to /ai/run, #992 minimax-code Claude reply translate, #993 release binaries off AVX, #996/#995 Antigravity TOS ban + appeal URL, #950 Codex reset-credit scope, plus version bumps to 1.9.16-exp.7 then 1.9.16-exp.8). Trial merge had 4 conflicts, all resolved: `fallback.go` keeps the fork's request-line helpers alongside upstream's new `!customJSON` guard; `registry_models.go` and `web/src/lib/models.ts` keep the fork's SenseNova 4-model list — upstream's decolua port re-added the China Token Plan models (`deepseek-v4-flash`, `deepseek-flash`, `glm-5.2`, `kimi-k3`) plus `glm-5.3`, none of which the global docs name; `background.go` takes upstream's real-status version of the OAuth park call — upstream's `IsRefreshGrantDead` is now byte-identical to the fork's wider check (same four dead-grant markers), so the fork's Nous parking fix is fully preserved by upstream's own code. Upstream's port also added a capability row for `sensenova-6.8-flash` (ContextWindow 1000000, MaxOutput 131072) — kept, it arrived with upstream's port. Upstream's cloudflare-ai image /ai/run routing (#985) coexists with the fork's clef systemone lane (different lanes). Fork behaviour preserved: CodeBuddy fresh request/conversation ids, fork-only `make build` (only `9router-go`), model context tags, Nous catalogue (10 models incl. the 3 serenhope seeds), Jev tilde alias, SenseNova `.ai` pointing + 4-model list, clef systemone routing, Exa probe + fetch translation, media-page Add Connection `apiKey` fix (fragment 946). VERSION is now 1.9.16-exp.8. Verified: `make build` passes (only `9router-go`, version 1.9.16-exp.8 embedded), `go test ./internal/providers/` and `./internal/handlers/dashboard/` pass, web unit 840 pass / 0 fail (the old timezone-only failure is gone upstream), svelte-check 74 errors == baseline 74. Known pre-existing failures, unchanged: executor `TestForwardQoder_UnknownModelNamesTheModel` (VM egress proxy TLS); media `TestHandleAudioVoices_elevenlabs` and `TestTargetPrivate_RespectsAllowPrivateNetworking`; chat `TestLocalGatewayConnectionPermittedWithOptIn` (fails on pristine upstream too — the VM egress proxy resolves api.deepseek.com to 198.18.25.163, which the #840 SSRF guard blocks). Note: upstream's own `CHANGELOG.md` still carries committed conflict markers from an earlier upstream merge (`1b5688b`, Nace/NIM fragments) — identical on both sides, left untouched. No provider connections were checked or tested.

### Fix: connections added from the media provider page keep their API key (2026-11-07)

- The Add Connection modal on the media (web search/fetch) provider page sent the pasted key as `key:`, but `POST /api/connections` only reads `apiKey` (or `data.apiKey`), so the key was silently dropped at creation. The connection was stored active but key-less, and every search request went out with an empty `x-api-key` — which is why Exa answered `{"code": "invalid_api_key", "message": "no API key found"}` for a valid official key. Every other add-connection flow already sends `apiKey:`; only this page was wrong. One-line fix: the modal now sends `apiKey: newConnKey.trim()`. Fragment `.changes/946-media-add-connection-api-key.md`.
- Note: connections created before this fix stay key-less; re-add them from the media page, or paste the key via the Providers tab Edit modal (which has the API key field with a Check button). The media page's inline Edit modal only edits Display Name + Active. No provider connections were checked or tested live.

### Upstream sync to 5e496d12 (2026-11-07)

- Synced upstream luqman-v1/9router-go through `5e496d12` (48 commits after `a239df34`: #939 honor SearchURL on direct search, #906 Nace Drex lane, #901 combo call budget, #969 Vertex / Perplexity Web / gemini-cli executors, #940 oh-my-pi wire contracts, #967 Gemini /v1beta surface, #840 per-connection SSRF opt-out, #945 Codex/OpenAI OAuth quarantine, #916 Cloudflare FLUX.2 reference image, #918 cloudflare-ai vision parts, #947 stream tok/s metrics, #925 Headroom autostart, #761 visual model/client editor, #958/#960 Copilot routing, #977 Kimi reasoning enum, #966 thought-signature normalization, #935 Copilot embeddings wrap, #964 Antigravity fingerprint, #952 minimax-code agent id, #954 Volcengine GLM-5 cap, #956 token-saver bypass header, #913 minimax-code broker, #924 HTML error sanitize, #933 Gemini schema normalization, #928 errorclassify, #931 Gemini TTS models, #920 Gemini tool preservation, #961 CommandCode max_tokens cap, #911 Claude identity block, #899 cursor think-tag decode, #897 minimax output_config strip, #943 quota table mobile, #902 coverage gate, #908 wiki search, plus version bumps to 1.9.16-exp.5). Trial merge had 7 conflicts, all resolved: VERSION and updater take 1.9.16-exp.5; `fallback.go` keeps the fork's request-line helpers alongside upstream's new `streamNoteWriter`; `oauth_lock_test.go` takes upstream's newer table (adds `refresh_token_reused`) — the fork's `IsRefreshGrantDead` covers every case, and the previously failing `TestRejectedOAuthRefreshParksAccount/expiry_refresh/spent_token_is_parked` subtest now passes; `background.go` keeps the fork's wider `IsRefreshGrantDead` check on upstream's error-typed `RecordConnectionOAuthFailure` call; `validate.go` takes upstream's `validateSystemOneLane` refactor while the fork's Exa web probe stays; `media.go` takes upstream #939's SearchURL hunk including its `!hasCustomBaseURL` guard, which supersedes the fork's identical-except-guard fragment 945 search hunk, while the fork's fetch `urls` translation stays. Upstream also deleted 44 old `.changes/` fragments; fork fragments 813/818/939-945 survive. Fork behaviour preserved: CodeBuddy fresh request/conversation ids, fork-only `make build` (only `9router-go`), model context tags, Nous catalogue (10 models incl. the 3 serenhope seeds), Jev tilde alias, SenseNova `.ai` pointing + 4-model list, clef systemone routing, Exa probe + fetch translation. VERSION is now 1.9.16-exp.5. Verified: `make build` passes (only `9router-go`, version 1.9.16-exp.5 embedded), `go test ./internal/providers/` and `./internal/handlers/dashboard/` pass, web unit 795 pass. Known pre-existing failures, unchanged: executor `TestForwardQoder_UnknownModelNamesTheModel` (VM egress proxy TLS); media `TestHandleAudioVoices_elevenlabs` and `TestTargetPrivate_RespectsAllowPrivateNetworking`; chat `TestLocalGatewayConnectionPermittedWithOptIn` (fails on pristine upstream too — the VM egress proxy resolves api.deepseek.com to 198.18.25.163, which the new #840 SSRF guard blocks); web unit 1 timezone-only failure (`consoleLog` expects UTC, VM is Asia/Jakarta). Note: upstream's own `CHANGELOG.md` carries committed conflict markers from an earlier upstream merge (`1b5688b`, Nace/NIM fragments) — identical on both sides, left untouched. No provider connections were checked or tested.

### Fix: Exa (and other web search/fetch providers) actually work (2026-11-07)

- The dashboard Test/Validate probe misclassified Exa as an Anthropic provider (`isAnthropicProbe` keys on AuthHeader `x-api-key`) and POSTed an Anthropic messages body to `https://api.exa.ai/search`, which Exa cannot accept; a valid official key was reported as "Invalid API key". `validateProviderKey` now probes web-only providers the way upstream's `probeWebProvider` does: a minimal POST of `{"query": "ping"}` to the provider's SearchURL (or FetchURL) with its own auth header; only 401/403 mean a bad key. Fragment `.changes/945-exa-web-provider-probe-search-fetch.md`.
- The direct `/v1/search` lane ignored SearchURL and dialed `https://api.exa.ai/search/v1/search`, a route Exa does not serve; it now uses SearchURL like the combo lane already did. Affects every SearchURL provider, not just Exa.
- `/v1/web/fetch` forwarded the client's singular `{"url": ...}` body as-is; Exa `/contents` and Tavily `/extract` require `{"urls": [...]}`. Both lanes now translate it; bodies already carrying `urls` pass through. No provider connections were checked or tested live.

### Fix: SenseNova model list matches the official docs (2026-11-05)

- Removed the four China Token Plan models (`deepseek-v4-flash`, `deepseek-flash`, `glm-5.2`, `kimi-k3`) from the `sensenova` list (backend + dashboard); the list is now exactly the four models the global docs name (`sensenova-6.8-flash`, `sensenova-6.8-flash-lite`, `sensenova-u1.5-lite`, `sensenova-u1.5-fast`). Follows the endpoint move in the previous entry. Fragment `.changes/944-sensenova-model-list-matches-docs.md`. Untouched: global capability rows for those four ids (other providers still serve them) and the `hcnsec` relay lists. No provider connections were checked or tested live.

### Fix: SenseNova points at the global platform (token.sensenova.ai), not the China endpoint (2026-11-05)

- Provider `sensenova` pointed at the China Token Plan endpoint (`token.sensenova.cn`). Per the current docs at `https://platform.sensenova.ai/docs`, the global platform base URL is `https://token.sensenova.ai/v1` (chat `POST /v1/chat/completions`, images `POST /v1/images/generations`, keys `sk-` from `https://platform.sensenova.ai/console/keys`); BaseURL, ImageURL, Website, and API-key URL switched to `.ai`. Added missing documented model `sensenova-6.8-flash` to the catalogue (backend + dashboard); no capability row, since the docs give pricing in points, not a context window. Fragment `.changes/943-sensenova-global-platform-ai.md`. Unverified: `deepseek-v4-flash`, `deepseek-flash`, `glm-5.2`, `kimi-k3` were China Token Plan models and the `.ai` docs do not name them; if they 404 on `.ai`, remove them. Note: a `.cn` API key does not work on `.ai`, a new key is needed; a custom node still pointing at `token.sensenova.cn` keeps working (the reasoning strip already accepts both hosts). No provider connections were checked or tested live.

### Fix: OpenRouter Jev latest id is ~typesafe/jev-latest, not typesafe/jev-latest (2026-11-04)

- Corrected the experiment from the previous entry: OpenRouter rejects `typesafe/jev-latest` with HTTP 400 "Model typesafe/jev-latest does not exist". Per OpenRouter's Jev Latest page, the rolling alias is `~typesafe/jev-latest` (tilde prefix) and it always redirects to the latest model in the Jev family (currently Jev 1.14). Catalogue, kind tag, dashboard entry, and test updated to `~typesafe/jev-latest`; `typesafe/jev-1.14` unchanged. Fragment `.changes/942-openrouter-jev-latest-tilde-alias.md`. No provider connections were checked or tested live.

### Experiment: OpenRouter systemone models typesafe/jev-1.14 and typesafe/jev-latest (2026-11-04)

- Added `typesafe/jev-1.14` and `typesafe/jev-latest` to the OpenRouter catalogue (backend `registry_models.go` list + `systemone` kind tag, dashboard `web/src/lib/models.ts`) next to the existing `typesafe/jev-1.13`; routing is unchanged since OpenRouter already has a `SystemoneURL`. Note: OpenRouter documents the rolling alias as `~typesafe/jev-latest` (tilde prefix); the id added here is exactly `typesafe/jev-latest` as requested, and the tilde form is the fallback if a live call rejects it. Fragment `.changes/941-openrouter-jev-114-latest.md`. No provider connections were checked or tested live.

### Fix: cloudflare-ai systemone routing to the Workers AI run endpoint (2026-11-04)

- Fixed the 400 code 7000 ("No route for that URI") on cloudflare-ai decision models (`@cf/cloudflare/clef-flash`, `@cf/cloudflare/clef`): the `/v1/systemone` lane used to send them to `.../ai/v1/systemone`, which Workers AI does not route. The lane (and the smart-router evaluation lane) now sends them to `.../ai/run/@cf/cloudflare/<model>` and puts the short name (`clef-flash`/`clef`) in the body model field, the only form Workers AI accepts there. Other providers keep the generic systemone rules. New helper `providers.CloudflareSystemoneURL`/`CloudflareSystemoneModel` plus unit, mock-upstream handler, and smart-router tests; fragment `.changes/940-cloudflare-systemone-run-route.md`. Known limitation: the Cloudflare `{"success": true, "result": {...}}` envelope is forwarded as-is, so callers expecting top-level `answers` must read `result.answers`. No provider connections were checked or tested live.

### Docs: update GOAL.md and HANDOFF.md state to the a239df34 sync and Nous port (2026-11-04)

- Updated the state sections of `GOAL.md` and `HANDOFF.md` from the superseded picture (fork up to date at `d6dd9481`, VERSION 1.9.16-exp.1) to the actual state: upstream synced through `a239df34` (12 commits after `d6dd9481`, snapshot push `68372337`, VERSION 1.9.16-exp.2), serenhope Nous seed models ported (push `d1cdeb7c`, Nous catalogue now 10 models), upstream head is `a239df34`, and the next sync covers only commits after `a239df34`. Also recorded that the Nov 04 snapshot pushes brought back the merge-base counting trap: the git merge-base is `d6dd9481` again, so a naive rev-list reports 12 already-synced commits; future syncs must count from `a239df34`. No code changed in this commit.

### Port: Nous seed models from serenhope/9router (2026-11-04)

- Restored 3 serenhope seed models to the fork's Nous static catalogue (`registry_models.go`) and capabilities (`capabilities.go`), verbatim from serenhope `open-sse/providers/registry/nous.js` and `capabilities.js`: `meituan/longcat-2.5-preview:free`, `stepfun/step-5-preview:free`, and `stealth/missingno`, all with tools, reasoning and the `nous` thinking format, so the Nous thinking dialect now applies to them too. The newer beam/ornith rows are kept. Tests updated in the same change (`nous_test.go` capabilities table, `aggregators_test.go` catalogue list); the thinking-levels fixture (v0.5.109) already listed all 8 serenhope models with base levels, so the now-obsolete pending-divergence entry for the two previously dropped models was removed from `thinking_levels_fixture_test.go`. Verified: `go test ./internal/providers/` passes; `go test ./internal/handlers/chat/` has exactly one failure, the pre-existing `TestRejectedOAuthRefreshParksAccount/expiry_refresh/spent_token_is_parked` (refresh-slot test state not reset between subtests, already noted in the `a7d36edf` sync baseline, unrelated to Nous). No provider connections were checked or tested.

### Upstream sync to a239df34 (2026-11-04)

- Synced upstream luqman-v1/9router-go through `a239df34` (12 commits after `d6dd9481`: #863/#879 dashboard mobile wrap, #865 cloudflare-ai accountId, #894 codex namespace tools, #898 chat default stream false, #900/#904 claude thinking/whitespace, #905 commandcode conversation-stable session identity, #896 cloudflare FLUX size, #921 loop-detection false positive, #922 gemini timeout failover, #938 SSE delta.reasoning test, plus the 1.9.16-exp.2 version bump). Trial merge was clean with 0 git conflicts, 41 files, no deletions and no binary files. Fork behaviour preserved: CodeBuddy fresh request/conversation ids, fork-only `make build` (only `9router-go`), model context tags, and the nous `invalid_grant` parking fix. VERSION is now 1.9.16-exp.2. Verified: `make build` passes (only `9router-go`, version embedded), `go test ./internal/proxy/` passes, web unit 768 pass (previous tree; #938 adds only a stream-peeker test case, re-verified in `./internal/proxy/`). Known pre-existing failure, unchanged: executor `TestForwardQoder_UnknownModelNamesTheModel` (VM egress proxy TLS). No provider connections were checked or tested.

### SenseNova strict reasoning history (2026-11-04)

- Strip replayed `reasoning_content`, `reasoning`, and `reasoning_details` fields from assistant history sent to SenseNova's strict Chat Completions endpoints, including custom nodes targeting `token.sensenova.ai` or `token.sensenova.cn`. This fixes the generic `400 inference request is invalid` on affected multi-turn requests while preserving content and leaving other providers unchanged. Verified with `go test ./internal/proxy/executor`.

### Docs: update GOAL.md and HANDOFF.md state to the a7d36edf / d6dd9481 sync (2026-11-03)

- Updated the state sections of `GOAL.md` and `HANDOFF.md` from the superseded picture (fork `ab5bb7ba`, upstream `4249283b` pending) to the actual state: merge `a7d36edf` synced upstream through `d6dd9481` (#895, 24 commits after `ae4b4149`), upstream head is `d6dd9481`, so the fork is up to date and the next sync covers only commits after `d6dd9481`; VERSION stays `1.9.16-exp.1`. Also recorded that the real merge resolved the snapshot-era divergences: the 45 old `.changes/` fragments are deleted (count 65 down to 51) and the git merge-base is now `d6dd9481` itself, so the old merge-base counting trap no longer applies. No code changed in this commit.

### Docs: HANDOFF.md converted to English only (2026-11-03)

- Removed all Javanese/Indonesian wording from `HANDOFF.md` (owner line, provider-check standing instruction quote, dropped-thread quotes, language-mirroring line), replacing each with an English equivalent with the same meaning. `GOAL.md` and `AGENTS.md` were checked and were already English only, so they were not changed in this commit. Note: this commit was based on main `a7d36edf`; the state sections of `HANDOFF.md` (fork `ab5bb7ba`, upstream `4249283b` pending) predate that sync and are superseded by it, see the trial-sync entry below.

### Trial local sync: luqman-v1/9router-go to d6dd9481 (2026-11-03)

- Synced the 24 upstream commits after the previous fork sync at `ae4b4149`, through `d6dd9481` (#895). The local trial merge includes the upstream provider-registry/catalog refactor, resilience admission control, service-tier routing, media/MCP proxy-pool routing, HCNSec, OAuth context propagation, tool-argument healing, success-rate analytics, and Gemini embedding dimensions.
- Resolved 11 conflict areas while preserving fork behavior: CodeBuddy fresh per-request conversation IDs plus upstream CK API-key handling, the fork's `invalid_grant` OAuth parking behavior, model context tags, and fork-only build history. Upstream's dynamic provider catalog, translator sanitization, GLM context propagation, anomaly panel, and provider admission control were retained.
- Verification so far: executor, translator, providers, and OAuth package tests pass. Chat tests run through the existing pre-sync failure `TestRejectedOAuthRefreshParksAccount/expiry_refresh/spent_token_is_parked`, caused by the pre-existing refresh-slot test state not being reset between subtests. No provider connections were checked. This is a local trial merge only; nothing was pushed.

### Docs: publish GOAL.md and HANDOFF.md to the repo root (2026-11-03)

- Added `GOAL.md` (fork maintenance goal, updated from the stale 1.9.15-exp.15 / Tailscale-route wording to the current state: last sync 1.9.16-exp.1 at `ab5bb7ba`, upstream ahead again at `4249283b`, next sync pending) and `HANDOFF.md` (full handover for another agent: job definition, hard rules, fork changes to preserve, sync procedure, snapshot-push limits, pending VPS items) at the repo root, so an agent that only clones the repo can find them. Previously both files lived only in Muse's goal workspace, not in the repo. The goal-folder copies remain the working copies for Muse agents; the repo copies are published from them.

### Upstream sync: luqman-v1/9router-go to 1.9.16-exp.1 (ae4b4149, trial merge 2026-11-01)

- Synced 15 genuinely new upstream commits since v1.9.15 (`1b73cafe`): `ae4b4149` ZCode Start/Weekend plan lane (#861), `a832c429` Tier 1 gap ports (#848), `36e4e32a` claude-cloaking fix (#850), `302f4e16` stream empty tool-call deltas (#847), `dd17af20` stall watchdog (#821), `636acdf7` combo reasoning streams (#845), `1e63` cache TTL (#839), `0ea93715` bump to 1.9.16-exp.1 (#843), `2e8a217e` tokenmaxxing combos (#786), `332bbbdb` bedrock reasoning (#831), `8fbd94f0` translator json_schema (#824), `7208370b` tool-arg auto-heal (#837), `0e02` strip n param (#836), `70530a3f` API key search/filter (#829), `2019d085` commandcode reasoning inject (#833). VERSION and `CurrentVersion` move 1.9.15 to 1.9.16-exp.1. Note: the git trial merge is measured against merge-base `f1eb955f` (159 files) because the previous v1.9.15 sync landed as a snapshot commit, not a merge, so that count re-applies the 16 already-synced commits.
- Conflict resolutions (6 files, rest auto-merged): `VERSION` + `internal/updater/updater.go` taken from upstream (1.9.16-exp.1); `internal/proxy/executor/codebuddy.go` kept the fork's fresh-per-request conversation ids and took only upstream's new stall-timeout line (upstream's session-pinned derivation dropped again, same as the v1.9.15 sync); `internal/proxy/executor/providers.go` took upstream's `InjectReasoningContent` call before the session-aware commandcode body builder (fork's session threadId preserved); `internal/proxy/executor/commandcode_request_test.go` kept upstream's new reasoning-block test; `web/src/components/keys/ApiKeysTable.svelte` kept upstream's rewritten mask comment (the #829 search/filter/cards rewrite merged around it).
- Known divergence, same as the v1.9.15 sync: the 45 upstream-deleted `.changes/` fragments are NOT deleted by a snapshot push (snapshot push cannot delete), so they stay on fork main and the dashboard changelog keeps showing them; no code files are deleted by this merge. Fork files untouched: Makefile fork-only build line (fragment 818), model context tags, and this changelog.
- Verified before push: `make build` pass (only `9router-go`, no `opencode-go`), `go vet` executor pass, codebuddy/commandcode executor tests pass, providers tests pass, web build pass, web unit (src/) 688 pass with 1 fail that is timezone-only (`consoleLog` test expects UTC, VM runs WIB), svelte-check 75 errors (below the previous 77 baseline, none in the conflict files). Two failures are pre-existing and unrelated to this merge: executor `TestForwardQoder_UnknownModelNamesTheModel` fails identically on fork HEAD `da3cc653` (VM proxy TLS handshake to api3.qoder.sh), and web e2e tests fail in this VM because the Playwright dashboard harness cannot start. No provider connections checked.

### Fork-only: `make build` stops compiling the `opencode-go` copy (fragment 818)

- Removed the second `go build` line from the `Makefile` `build` target, which compiled the same gateway package twice (as `9router-go` and as `opencode-go`). Verified in a clean worktree: `make build` produces only `9router-go` (version 1.9.15 embedded), no `opencode-go` file. No behaviour change: upstream CI/Docker/release build only `9router-go`, nothing invokes the `opencode-go` binary, and the `opencode-go` provider (ocg) is routed inside the gateway regardless of binary filename. Note: a stale `opencode-go` file from an older build is not removed by `make clean`; delete it once by hand (`del opencode-go.exe` on Windows).

### Upstream sync: luqman-v1/9router-go to v1.9.15 (1b73cafe, direct push 2523972b, PR #4 closed unmerged)

- Synced 16 upstream commits since merge-base `f1eb955f` to release `v1.9.15` (head `1b73cafe`, 293 files upstream, VERSION exp.16 to 1.9.15). Trial merge conflicted in 2 files only; the other 88 merged cleanly. This sync did NOT go through a PR merge: cross-fork PR #4 was created via the browser but could not be resolved in GitHub's web editor, because that editor commits to the head branch (`luqman-v1:main`) which the fork owner cannot push to, so PR #4 was closed unmerged with an explanatory comment and the resolved tree was pushed directly to main (commit `2523972b`, 41 files).
- Conflict resolutions: `CHANGELOG.md` taken wholesale from upstream (fork history lives in this file, not upstream's); `codebuddy.go` combined the fork's fresh-per-request conversation ids (3fd925a port) with upstream's CN CLI identity (2.133.1, `#3257`), dropping upstream's session-pinned derivation whose helpers were already removed as dead code by the fork fix.
- Known divergences from a full merge, documented for the next sync: 45 upstream-deleted `.changes/` fragments were NOT deleted here (snapshot push cannot delete), so the dashboard changelog will show those entries twice until the next sync deletes them; 4 upstream screenshot assets (1 PNG + 3 WebP, docs only, not referenced by the dashboard) were NOT pushed (snapshot push is UTF-8 only). Fork fragment `.changes/813-model-context-tags.md` survives.
- Verified before push: `go build ./...` pass, `go vet` on executor pass, codebuddy tests pass (fork fresh-per-request + session-not-pinned, upstream CN version + identity + telemetry), providers tests pass, web build pass, web unit 689 pass. No provider connections checked.

### Upstream sync: luqman-v1/9router-go to f1eb955f (PR #3, merge 00f20381)

- Merged 12 upstream commits since merge-base `38e29580` via PR #3, mergeable_state clean, 0 git conflicts, 70 files changed (5790 insertions, 716 deletions). Upstream head `f1eb955f` (#817 key-name validation + Windows test fixes, #815 arewedaks gap analysis, #813 oauth device slowdown backoff, #816 freebuff gzip recovery, #811 dashboard sidebar search/pins/icon rail, #814 tunnel Global tab). VERSION stays 1.9.15-exp.16. Fork-local files untouched (context tags, codebuddy header fix, provider filter, CHANGELOG.fork.md).
- Note: `.changes/813` now has 2 fragments (fork model-context-tags, upstream oauth-device-slowdown-backoff) with different filenames, no git conflict, same pattern as the `.changes/790` triple from PR #1.
- Verified before merge: `go build ./...` pass, web build pass, providers tests pass, web unit 682 pass (up from 658 with upstream's new sidebar/command-palette tests). No provider connections checked.

### Model context tags: show each model's max context where it is picked (port from serenhope/9router v0.5.117-Custom)

- Ported "Every model shows its real window" from serenhope/9router
  (v0.5.117-Custom, `ModelSelectModal.js` ContextTag +
  `src/shared/utils/contextWindow.js`): a model now carries its context
  window next to its name (1.0M, 256k, 200, full token count in the tooltip),
  so a 1M model is never picked as if it were a small one. Literal git
  cherry-pick was impossible: serenhope is Next.js/JS, this fork is Go+Svelte,
  so the tag was rewritten in Svelte against data this fork already serves.
- Shown in three places: model picker chips, combo member rows in the combo
  editor, and the Available Models grid on each provider page (caps endpoint
  wins there, static catalog is the fallback). No backend change; a model
  with no known window shows no tag, never a guess.
- Coverage note: only ~77 static catalog entries carry a `contextLength`,
  so picker/combo tags appear only where the catalog knows the window; the
  provider grid covers more via the caps endpoint (synced catalogs, custom
  models). Custom/alias models in the picker show no tag yet.
- Verified: web unit tests 658 pass, 0 fail; `bun run build` pass;
  svelte-check 77 errors, same as baseline, none in the changed lines.
  No provider connections checked. Fragment: `.changes/813-model-context-tags.md`
  (813 is fork-assigned, not an upstream PR number).
- Files: `web/src/lib/formatContextWindow.ts` (+ test),
  `web/src/lib/models.ts` (+ test), `web/src/components/combos/pickerData.ts`,
  `ModelPill.svelte`, `ModelPickerModal.svelte`, `CreateComboModal.svelte`,
  `web/src/components/connections/types.ts`, `ProviderDetailView.svelte`.

### Upstream sync: luqman-v1/9router-go to 1.9.15-exp.16 (PR #2, merge 3a5d7e52)
- Merged 11 upstream commits since merge-base `b7816cf6` via PR #2, mergeable_state clean, 0 git conflicts, 70 files changed (7849 insertions, 215 deletions). Upstream head `38e29580` (#807 combo-to-combo fallback, #812 Cloudflare Tunnel, #802 opencode alias fix, bump to exp.16, #798 wire-contract table, #797 anthropic-compatible Claude wire, #795 client context limits, #774 combo no-same-connection resend, #782 live-catalog negative cache, #755 opencode probe test message, #793 route response headers). Fork-local files untouched (provider filter, connected filter, nous fix, codebuddy header fix, CHANGELOG.fork.md).
- Verified before merge: `go build ./...` pass, web build pass, providers tests pass. Executor, tunnel and chat packages each showed only sandbox live-network failures (Qoder/DeepSeek TLS via egress proxy, .invalid DNS sinkhole), verified identical on pure upstream main `38e29580`. No provider connections checked.

### CodeBuddy: plugin request/conversation id headers (upstream decolua/9router@3fd925a)
- Chat requests to `codebuddy-cn` and `codebuddy-intl` now send fresh ids per request with `X-Conversation-Request-ID` and `X-Conversation-Message-ID` both equal to `X-Request-ID`, and a fresh `X-Conversation-ID` per request. Previously Go minted an independent conversation-request id and pinned the conversation id to the session (#581), diverging from the plugin contract the upstream channel check verifies (400 / 11128 "request illegal" without the plugin ids). `stableConversationUUID` and `firstUserText` removed as dead code. Backoff/combo half of the upstream commit was already ported (#5024) and is unchanged. Verified: codebuddy executor tests, provider classify tests, `go build ./...`, `go vet` pass; no provider connections checked. Fragment: `.changes/797-codebuddy-plugin-headers.md`.
- Files: `internal/proxy/executor/codebuddy.go`, `codebuddy_conversation_headers_test.go`, `codebuddy_test.go`.

### Upstream sync: luqman-v1/9router-go to 1.9.15-exp.15 (PR #1, merge 5e5e9b69)
- Merged 25 upstream commits since merge-base `acddcdd` via PR #1, mergeable_state clean, 0 git conflicts, 175 files changed. Trial merge locally covered the first 23 (to 1.9.15-exp.15); GitHub added 2 more before merge (#794 stream usage, #796 batched key-list usage read). Fork-local files (provider filter, connected filter, nous invalid_grant fix, CHANGELOG.fork.md) untouched by the merge.
- Verified before merge: `go build ./...` pass, web build pass, `bun test` 666 pass in UTC (1 timezone-only failure in WIB), svelte-check 77 errors, same as baseline, 0 unresolved identifiers. Go test note: 2 packages failed only on live-network tests in this sandbox (DeepSeek/ElevenLabs TLS via egress proxy, example.com SSRF check), not on merged code.
- Fragment files upstream are carried as-is; note `.changes/790` now has 3 fragments (fork provider-filter, upstream token-normalizer, upstream qoder-pool) with different filenames, no git conflict.

### Provider filter: three modes only (`all` / `connected` / `not_connected`)
- Upstream shipped a 5-option status filter (`all`, `connected`, `error`,
  `disabled`, `not_connected`) in the providers overview. Fork simplifies to
  the 3 requested modes: all providers, only providers with connections, only
  providers without connections. "Connected" means holding any connection
  (oauth or api key) regardless of state: errored or disabled rows still count
  as having a connection, and legacy `error`/`disabled` values map to
  `connected` for the same reason. Covered by `matchesFilter` tests in
  `types.test.ts` (`bun test` 12 pass). Fragment:
  `.changes/790-provider-filter-three-modes.md`.
  Files: `web/src/components/connections/ProvidersOverviewGrid.svelte`,
  `web/src/components/connections/types.ts`,
  `web/src/components/connections/types.test.ts`.

### Fork rule: every fork change logged here
- `AGENTS.md` step 6 now requires every fork change to also append an entry
  here under `[Unreleased]`. `CHANGELOG.fork.md` is append-only per change;
  the release process files rows under a version heading the same way
  `make changelog-merge` does for `CHANGELOG.md`.

### Nous refresh: park 400 `invalid_grant` as dead grant
- `3c22c5c` fix(nous): `refreshNous` returned unwrapped text on non-2xx, invisible to
  `IsRefreshGrantDead`, so a spent single-use token was retried on every
  request, background, and usage refresh. Now returns a typed
  `OAuthRefreshError`; chat/background parkers use the shared grant-dead gate
  (was 401-only), the account shows re-authorize instead of hammering the token
  endpoint. Verified: `go test ./internal/proxy/oauth/ ./internal/providers/
  ./internal/handlers/chat/` green, incl. new `nous_grant_test.go` fixture
  (spent grant parks, transient 4xx/5xx stay retryable). Fragment:
  `.changes/774-nous-refresh-invalid-grant-parking.md`.

### Backfill: fork-local changes since the `acddcdd` fork point
- `3c7dd0c` / `e5a6fb3` fix(nous): free catalog refreshed against the live
  endpoint. Dropped `meituan/longcat-2.5-preview:free`,
  `stepfun/step-5-preview:free`, `stealth/missingno` (gone from
  `GET /v1/models`, probed 2026-10-28, upstream 404s); added
  `reflection/beam-1.0-501b-a23b:free` (reasoning mandatory) and
  `ornith/ornith-1.5-122b:free` (vision, 262144 max output). Fragment
  `.changes/755-nous-catalog-live-refresh.md` (folded into `CHANGELOG.md`
  by `c7df158`, no longer on disk).
- `0e41829` / `ec41096` feat(log): request/RTK/usage log lines match the
  Next.js gateway shape. One correlated dispatch line with stable
  per-connection emoji tag; RTK line mirrors `formatRtkLog`; usage DONE
  line mirrors `formatDoneLine`; `[HH:MM:SS]` stamps, JSON shape unchanged.
  Fragment `.changes/756-js-style-log-lines.md` (folded by `c7df158`).
- `c7df158` docs: folded 206 pending `.changes/` fragments into
  `CHANGELOG.md` (no version bump, no tag). Folder hygiene only, fragments
  archived in `CHANGELOG.md`, not fork-new behavior.
