# Changelog (fork `shane1794/9router-go`)

Fork-local history. Every change landing in this fork gets an entry under
`## [Unreleased]`, newest first — in addition to the `.changes/` fragment
required by the upstream convention (see `AGENTS.md` step 6). This file never
merges into upstream `CHANGELOG.md`.

## [Unreleased]

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
