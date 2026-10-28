# Changelog (fork `shane1794/9router-go`)

Fork-local history. Every change landing in this fork gets an entry under
`## [Unreleased]`, newest first — in addition to the `.changes/` fragment
required by the upstream convention (see `AGENTS.md` step 6). This file never
merges into upstream `CHANGELOG.md`.

## [Unreleased]

### Fork rule: every fork change logged here
- `AGENTS.md` step 6 now requires every fork change to also append an entry
  here under `[Unreleased]`. `CHANGELOG.fork.md` is append-only per change;
  the release process files rows under a version heading the same way
  `make changelog-merge` does for `CHANGELOG.md`.

### Nous refresh: park 400 `invalid_grant` as dead grant
- `refreshNous` returned unwrapped text on non-2xx, invisible to
  `IsRefreshGrantDead`, so a spent single-use token was retried on every
  request, background, and usage refresh. Now returns a typed
  `OAuthRefreshError`; chat/background parkers use the shared grant-dead gate
  (was 401-only), the account shows re-authorize instead of hammering the token
  endpoint. Verified: `go test ./internal/proxy/oauth/ ./internal/providers/
  ./internal/handlers/chat/` green, incl. new `nous_grant_test.go` fixture
  (spent grant parks, transient 4xx/5xx stay retryable). Fragment:
  `.changes/774-nous-refresh-invalid-grant-parking.md`.
