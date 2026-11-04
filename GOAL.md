# 9router-go fork maintenance

Goal ID: goal_92d8adc4d99f
Goal slug: 9router-go-fork-maintenance

## Description
Arthur maintains a public fork of luqman-v1/9router-go, and Muse's standing job is to sync upstream, check conflicts, merge to main and push when safe, and push requested feature or fix changes with CHANGELOG.fork.md updates. The fork was last synced to upstream a239df34 (2026-11-04, 12 commits after d6dd9481, snapshot-pushed as main 68372337), then the serenhope Nous seed models were ported and pushed as main d1cdeb7c; upstream head is a239df34, so the fork is up to date as of 2026-11-04, and the next sync is only for upstream commits after a239df34. VERSION is 1.9.16-exp.2. Note: the Nov 04 pushes were snapshots, not real merges, so the git merge-base is back to d6dd9481 and a naive rev-list overcounts; always count from the last synced head a239df34. The Muse VM is a workshop, not 24/7 hosting, and the laptop-via-Tailscale route was dropped 2026-10-30, so neither is an always-on route. See HANDOFF.md in this folder (also published at the repo root, English only since 2026-11-03) for the full handover: job definition, hard rules, fork changes to preserve, sync procedure, and snapshot-push limits.
