# HANDOFF: 9router-go fork maintenance

Published at the fork repo root as HANDOFF.md (2026-11-03). The copy
in the goal folder is the working copy for Muse agents; the repo
copy is published from it.

For the next agent picking up this project. Read this file first, then
`GOAL.md` in this folder, then `CHANGELOG.fork.md` in the fork repo.
Do not reconstruct history from GitHub search alone: search hides forks
by default, and snapshot pushes do not preserve local git history.

Written: 2026-11-03. Owner: Arthur. Previous session: side chat
"Continue 9router-go fork and fixes".

## 1. Project identity

- Fork: `shane1794/9router-go`, public, branch `main`, language Go.
  Created 2026-10-28. GitHub login: `shane1794`, Muse GitHub App
  installed with All repositories.
- Upstream: `luqman-v1/9router-go`, branch `main`. The fork's
  upstream ancestor is decolua/9router, but syncs are against
  luqman-v1 only.
- What it is: a high-performance native Go gateway (dashboard + proxy
  for many LLM providers). The repo's own `AGENTS.md` is the upstream
  engineering guide; section 7.6 is the one that matters most for
  fork work (see section 5 below).

## 2. Current state (verified 2026-11-03)

- Fork `main`: `ab5bb7ba` (snapshot push, 2026-11-01).
- Last upstream sync landed: `1.9.16-exp.1`, upstream head
  `ae4b4149` (15 genuinely new commits since v1.9.15).
- `VERSION` and `internal/updater/updater.go` `CurrentVersion`:
  `1.9.16-exp.1`.
- Upstream `main` NOW: `4249283b` (fix translator regex lookarounds,
  #886, plus at least 9 more commits on top of `ae4b4149`, including
  #885 tray, #884 ImgBot image optimisation, #870 proxy-pool,
  #860 service tiers, #858 retry budget, #854 HCNSec provider).
  The fork is behind upstream again. The next sync has NOT been
  started. Do not claim the fork is up to date.
- Counting trap: a naive `git rev-list HEAD..upstream/main` overcounts.
  Previous syncs landed as snapshot commits, not merges, so the git
  merge-base is still `f1eb955f` (the PR #3 sync point). The 15-commit
  figure for the last sync was counted from `1b73cafe`, not from the
  merge-base. Always count genuinely new commits from the last synced
  upstream head, and confirm with a trial merge.

## 3. The standing job (Arthur-confirmed 2026-10-30)

1. Sync upstream.
2. Check conflicts. If safe, merge to `main` and push to his repo.
3. If he asks to add/remove a feature or fix something, push that
   change and update `CHANGELOG.fork.md`.
4. Anything else depends on his request. Do not invent extra work.

## 4. Hard rules, learned the hard way

- One denial stops the work. If Arthur declines a push or any action
  on an approval card, do NOT retry the same action, do NOT re-send
  approval cards, and do NOT work around it through another route
  (browser, PR, etc.) unless he explicitly asks to try again.
- NEVER check or test a provider connection in 9router unless he
  explicitly asks (standing instruction, 2026-10-30). Build and unit
  tests are fine; live provider checks are not.
- Secrets: never ask for, accept in chat, store, or repeat a token,
  API key, or GitHub password. Tokens are entered on his own machine
  only. Never run `gh auth login`, device auth, or embed credentials
  in git commands. Never persist his API keys anywhere.
- Tailscale: the Muse VM's Tailscale client is client-only. It
  accepts no inbound connections and hosts nothing. Do NOT claim
  Tailscale access to the VM dashboard works until a real TCP
  request to it actually succeeds. Upstream `tailscale serve/funnel`
  is not available through that client and must not be installed
  as a workaround.
- The Muse VM is a workshop, not hosting. It is replaced roughly
  every hour, killing processes; files in `~` persist, processes
  do not. Never propose running the gateway 24/7 from the VM.

## 5. Fork-local changes that must survive every sync

These are fork work, not upstream work. A sync that drops one of
them is a broken sync, even if it builds.

1. CodeBuddy header fix (decolua/9router@3fd925a port, then merged
   with upstream CN identity): `internal/proxy/executor/codebuddy.go`
   sends FRESH request/conversation ids per request
   (`X-Request-ID == X-Conversation-Request-ID`,
   `X-Conversation-ID != X-Conversation-Message-ID`). Upstream's
   session-pinned derivation (`stableConversationUUID`,
   `firstUserText`) was removed as dead code by the fork fix and
   must NOT come back in a conflict resolution. In the last two
   syncs the resolution was: keep the fork's fresh ids, take only
   the upstream change that is genuinely new (CN CLI 2.133.1
   identity in the v1.9.15 sync; stall-timeout line in the
   v1.9.16-exp.1 sync).
2. Model context tags (ported from serenhope/9router
   v0.5.117-Custom): every model shows its real context window
   (1.0M, 256k, full count in tooltip) in picker chips, combo
   member rows, and each provider's Available Models grid.
   Files: `web/src/lib/formatContextWindow.ts` (+test),
   `web/src/lib/models.ts` `getCatalogContextLength` (+test),
   picker/ModelPill/ModelPickerModal, CreateComboModal,
   connections/types + ProviderDetailView,
   `.changes/813-model-context-tags.md`. Only about 77 static
   catalog entries carry a `contextLength`, so no tag where the
   catalog does not know it, never a guess.
3. Makefile fork-only build (fragment 818):
   `make build` compiles ONLY `9router-go`. The second `go build`
   line that also produced `opencode-go` was removed, with a
   fork-only comment explaining why.
4. `CHANGELOG.fork.md`: the fork's own history, newest entry first
   under `## [Unreleased]`. Every fork change gets an entry here.
   This file never merges into upstream `CHANGELOG.md`.
5. Provider filter and connected-filter count fixes (2026-10-29),
   nous `invalid_grant` parking fix, and the fork's
   `.changes/813` fragment above.

## 6. The opencode-go distinction (Arthur corrected this twice)

- `opencode-go` is a PROVIDER: id `opencode-go`, alias `ocg`,
  display name OpenCode Go, in the gateway's provider registry.
- The `opencode-go` BINARY was a second compile of the same gateway
  (upstream commit #556, menu/tray aliasing), not a provider-only
  binary. Upstream CI, Docker, and release build only `9router-go`.
  Removing the second compile does NOT remove the provider.
- `make clean` does NOT delete a stale `opencode-go` file (its
  pattern is only `9router-go-*`). After the Makefile fix, the
  operator must delete it once by hand: `del opencode-go.exe`
  (Windows) or `rm opencode-go` (VPS). Do not tell him the fix
  failed just because the stale file is still there.

## 7. How to do a sync (the procedure that worked)

1. Start from a FRESH clone, never from an old worktree. Old
   worktrees in `~/workspace/` (`9router-go-sync-check`,
   `9router-go-makefile-fix`, `9router-go-fork-install`,
   `9router-go-sync-nov01`) are dirty, stale, or on trial branches.
   The last one is on branch `trial-sync-nov01` with a staged
   merge and deletions; do not reuse it as a base. Fetch
   `origin/main` first: local git history here does not match
   remote history after snapshot pushes, and local clones may be
   missing git identity (set it locally in the clone only for a
   trial merge; snapshot pushes use the GitHub identity).
2. Add `upstream` remote (`luqman-v1/9router-go`), fetch it, and
   run a trial merge on a trial branch. Record: upstream head,
   genuinely new commit count, conflicting files, and the rest
   that auto-merged.
3. Resolve conflicts following the pattern in section 5, item 1:
   keep fork behaviour, take only genuinely new upstream lines.
   Previous conflict files and their resolutions:
   - v1.9.15 sync: `CHANGELOG.md` taken wholesale from upstream
     (fork history lives in `CHANGELOG.fork.md`, not there);
     `codebuddy.go` fork fresh ids + upstream CN identity.
   - v1.9.16-exp.1 sync (6 files): `VERSION` + updater from
     upstream; `codebuddy.go` fork fresh ids + upstream stall
     timeout only; `providers.go` upstream `InjectReasoningContent`
     before the session-aware commandcode builder (fork session
     threadId preserved); `commandcode_request_test.go` upstream
     reasoning-block test kept; `ApiKeysTable.svelte` upstream
     mask comment kept (the #829 rewrite merged around it).
4. Verify BEFORE asking to push, in this order:
   - `make build` (web build first, then Go build). Must produce
     only `9router-go`, no `opencode-go`.
   - `go vet` on the executor package.
   - Executor tests for codebuddy/commandcode, plus provider tests.
   - Web build, web unit tests in `web/src/` only (see section 9
     for why `bun test` without a path is misleading here).
   - `svelte-check` (ratchet, see section 9).
   - No provider connection checks (section 4).
5. Update `CHANGELOG.fork.md` with a new entry under
   `## [Unreleased]`, newest first: what synced, which upstream
   commits, how conflicts were resolved, known divergences, and
   the verification numbers. Also follow the repo `AGENTS.md`
   section 7.6 rule: upstream convention uses a `.changes/`
   fragment per change; for a sync the upstream fragments arrive
   with the sync itself. NEVER hand-edit upstream `CHANGELOG.md`
   (the v1.9.15 sync took it wholesale from upstream as a merge
   resolution, which is the only exception, and it is already
   done).
6. Present the resolutions, the verification numbers, and the
   snapshot limits (section 8) to Arthur and get an explicit go
   for the push. Conflicts are never "safe auto-merge": he decides.
7. Push with the snapshot procedure, not a normal git push:
   read `/opt/hatch/connectors/github/references/snapshot-push.md`
   first. In short: build a NUL-delimited manifest OUTSIDE the
   source dir from `git diff --cached --name-only
   --diff-filter=AM`, excluding binary files and sync records;
   verify the remote head still matches the recorded head
   immediately before pushing; run `github push-files --owner
   shane1794 --repo 9router-go --branch main --message ...`;
   then verify the new head with `list_commits` and by reading
   back `VERSION`, the conflict files, and `CHANGELOG.fork.md`
   from the remote. Record the confirmed SHA in the adjacent
   sync record file (see section 8), never from a later head
   lookup that could reflect another writer.

## 8. Snapshot push limits (disclose these every time)

- Snapshot push creates ONE commit with the submitted files. It
  does NOT preserve local git history, branches, tags, file modes,
  or symlinks, and it CANNOT delete files.
- It is UTF-8 text only. Binary files are rejected and would reject
  the whole snapshot if included, so omit them explicitly with
  Arthur's agreement. Docs-only binary assets are the usual case.
- Limits: at most 2,000 files, 16 MiB of content, manifest under
  1 MiB and NUL-terminated.
- Sync records live NEXT TO the source dir, never inside it:
  `<source-dir>.muse-github-snapshot-sync.json`, e.g.
  `~/workspace/9router-go-sync-nov01.muse-github-snapshot-sync.json`
  for the last sync (commit `ab5bb7ba`, 107 paths). The record
  can reveal other repos' targets, paths, and SHAs; it is excluded
  from every manifest.
- Known divergences carried forward (same in both last syncs):
  45 upstream-deleted `.changes/` fragments were NOT deleted on
  fork `main`, so the dashboard changelog shows those entries
  twice. No code files were deleted by either merge. The next
  sync will show the same 45 deletions in its trial merge again
  until a real (non-snapshot) push deletes them. In the
  v1.9.16-exp.1 sync, 8 binary PNG/WebP docs assets were also not
  pushed (1 request-details screenshot, 4 evidence PNGs for #829,
  3 WebP screenshots); the v1.9.15 sync left 4 similar assets
  behind. The new upstream head also contains an ImgBot image
  optimisation commit (#884), so expect more binary assets in the
  next sync and handle them the same way.

## 9. Verification baseline and pre-existing failures

Do NOT chase these as if the sync caused them. All were confirmed
on fork HEAD before the last sync, or are environment artefacts.

- Go executor: `TestForwardQoder_UnknownModelNamesTheModel`
  fails identically on `da3cc653` and on the merged tree: the VM
  egress proxy breaks the TLS handshake to `api3.qoder.sh`, so
  the error does not name the model. Not a merge regression.
- Web unit (`bun test src/`): 688 pass, 1 fail on the last sync.
  The fail is `consoleLog` `toPlainText`: the test expects UTC
  timestamps, the VM runs Asia/Jakarta (WIB), so 07:08 becomes
  14:08. Timezone only.
- Web e2e (`bun test` with no path also runs `e2e/`): fails in
  this VM because the Playwright dashboard harness cannot start.
  Never report the all-tests total as the unit result.
- `svelte-check`: baseline file pins 76 errors
  (`web/scripts/svelte-check-baseline.json`). The last sync
  measured 75, none in the conflict files. The ratchet may only
  shrink; never widen the baseline to make a failure disappear.
  `tsc` cannot parse `.svelte` files, so `bun run build` passing
  does not prove a component's imports are right; svelte-check
  is the gate for that class of bug.
- Previous baselines for reference: web unit 682 (PR #3 sync),
  689 (v1.9.15 sync, before the timezone fail was isolated),
  688 + 1 timezone fail (v1.9.16-exp.1 sync).

## 10. Toolchain and environment

- Go: `~/workspace/toolchains/go/bin/go`, go1.27.3, installed
  2026-11-01 because a VM replacement had left no Go at all.
  Muse exec `PATH` does NOT include it by default, so always use
  `PATH=$HOME/workspace/toolchains/go/bin:$PATH` for fork builds.
  `go.mod` requires go 1.27, toolchain go1.27.3.
- Bun: available in the image (`bun --version` 1.4.4 at last sync).
  Node is also present. The web build needs `bun install` first
  in a fresh clone (no `node_modules` committed).
- VM replacements also lose apt packages: the Hermes work needed
  `libgtk-3-dev` and `libwebkit2gtk-4.1-dev` reinstalled after a
  replacement. If a build suddenly misses a system package, check
  that first before blaming the code.
- GitHub access check (per the local connector amendment): GitHub
  search excludes forks by default. Never state a repo count or
  "missing access" from a default search. Search with
  `fork:true`, or read the named repo directly (e.g.
  `list_commits` on `shane1794/9router-go`), before concluding.

## 11. Arthur's machine and VPS: pending, not done

These were pending at handoff time. Verify them fresh; do not
assume they happened.

1. On his machine (Windows) and/or VPS: `git pull origin main`
   to `ab5bb7ba`, then `make build`, then stop/start the gateway.
   Confirm `9router-go` (version 1.9.16-exp.1 embedded) runs and
   `opencode-go` does not reappear (delete the stale file once
   by hand, section 6).
2. VPS gateway restart, port 20130: the last known blocker was a
   port occupant OUTSIDE the gateway's own daemon pidfile
   (identified one morning as 9router-go itself, PID 57761, but
   that PID may have changed or he may have cleaned it himself).
   Menu option 4 (stop background server) will NOT clean it,
   because it only knows the daemon's own pidfile. Next step is
   always a FRESH port recheck (`ss -tlnp | grep 20130`, then
   `ps -o pid,ppid,etime,cmd -p <PID>`) before any kill, and
   only kill it if the command is actually 9router-go.
   The gateway on his VPS is `instabox-tes-agent` via Tailscale
   for his other project; the fork gateway itself runs on port
   20130 with dashboard on port 200 locally in the older notes.
3. Windows `make`: install GNU Make with
   `winget install ezwinports.make` (GNU Make 4.4.x), close and
   reopen the terminal so PATH updates, then `make --version`.
   Do NOT use `winget install GnuWin32.Make` (old 3.81, does not
   support `$(file...)` used by the Makefile). Fallbacks:
   `choco install make` or MSYS2. The Makefile is written for
   cmd.exe, no WSL needed. Never tell him to run `go build`
   again after `make build`: `make build` already runs it and
   embeds the version. Go 1.27 and Bun 1.x must be on PATH, plus
   internet on the first run for modules.

## 12. Do not revive these (Arthur-closed, same handling)

- Botkeep 9router-go deploy: failed end to end, closed
  2026-10-09, reopened by Arthur himself 2026-10-13 with
  luqman-v1/9router-go, now blocked on "no server capacity".
  Any retry needs a FRESH API key from him first. Do not rebuild,
  re-offer, or pitch it unprompted.
- Laptop via Tailscale for the gateway: dropped 2026-10-30.
  Do not revive or pitch.
- SSH-to-WSL plan: dropped 2026-10-30.
  Do not revive. Tailscale enrollment itself stays as is.
- Adobe Stock videos, Honeygain, and the original muse.ai Threads
  promo closure: dropped threads, not this project, but the rule
  is the same: never revive, reframe, or pitch a thread Arthur
  closed. (Note: he reopened the Threads promo himself on
  2026-10-31 and posted it once; that reopening was his call,
  not an agent revival.)

## 13. How Arthur likes output (standing, every session)

- ADHD mode always on until he says stop ("stop adhd mode"):
  lead with the next action, number multi-step work, one
  concrete next action at the end, restate state across turns,
  suppress tangents, give specific time estimates where they
  exist, make wins visible, matter-of-fact errors, max 5 items
  per list group, no preamble/recap/closers.
- Antislop core always on, with the family skills when the work
  touches their domain. No em dashes in generated text: he reads
  them as AI slop, including in subtitles and narration.
- SANITY always on since 2026-11-01: verify claims against their
  source, label CONFIRMED / INFERRED / UNCERTAIN / UNKNOWN where
  it matters, do not present a full score dump on a casual reply.
- Mirror the language he is using in the current conversation.
  He crosschecks rather than agrees: when he challenges a claim,
  open the source and answer from it, and concede immediately
  when his evidence wins.
- Never generate images or video unless he explicitly asks.
  Quota cost is uncertain and posts late; never quote the old
  2 to 3 percent per clip figure.
- When he accepts something as good enough even while noting it
  is imperfect, close the loop and offer no further passes.

## 14. Next actions for the incoming agent

1. Confirm with Arthur whether the next job is the upstream sync
   to `4249283b` (section 2), or a specific fix he will name.
   Do not start the sync unprompted: the last two syncs both
   needed his explicit go, and the push needs it again.
2. If it is the sync, follow section 7 exactly, starting from a
   fresh clone, and expect the same 45 `.changes/` deletions to
   reappear in the trial merge plus any new binary assets from
   #884 and friends.
3. Ask Arthur for the section 11 confirmations (pull to
   `ab5bb7ba`, `make build`, gateway restart) only if he has
   not already confirmed them in another chat; check that chat
   first before asking him to repeat himself.
