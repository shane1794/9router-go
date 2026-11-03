# 9router-go fork maintenance

Goal ID: goal_92d8adc4d99f
Goal slug: 9router-go-fork-maintenance

## Description
Arthur maintains a public fork of luqman-v1/9router-go, and Muse's standing job is to sync upstream, check conflicts, merge to main and push when safe, and push requested feature or fix changes with CHANGELOG.fork.md updates. The fork was last synced to upstream d6dd9481 via merge a7d36edf (2026-11-03, 24 commits after ae4b4149); upstream head is d6dd9481, so the fork is up to date as of 2026-11-03 evening, and the next sync is only for upstream commits after d6dd9481. VERSION is 1.9.16-exp.1. The Muse VM is a workshop, not 24/7 hosting, and the laptop-via-Tailscale route was dropped 2026-10-30, so neither is an always-on route. See HANDOFF.md in this folder (also published at the repo root, English only since 2026-11-03) for the full handover: job definition, hard rules, fork changes to preserve, sync procedure, and snapshot-push limits.
