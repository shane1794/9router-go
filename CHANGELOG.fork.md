# Changelog (fork `shane1794/9router-go`)

Fork-local history. Every change landing in this fork gets an entry under
`## [Unreleased]`, newest first — in addition to the `.changes/` fragment
required by the upstream convention (see `AGENTS.md` step 6). This file never
merges into upstream `CHANGELOG.md`.

## [Unreleased]

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
