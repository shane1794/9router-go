# Changelog (fork `shane1794/9router-go`)

Fork-local history. Every change landing in this fork gets an entry under
`## [Unreleased]`, newest first — in addition to the `.changes/` fragment
required by the upstream convention (see `AGENTS.md` step 6). This file never
merges into upstream `CHANGELOG.md`.

## [Unreleased]

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
