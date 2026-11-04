### Experiment: OpenRouter systemone models typesafe/jev-1.14 and typesafe/jev-latest (2026-11-04)

- Added `typesafe/jev-1.14` and `typesafe/jev-latest` to the OpenRouter static catalogue (`registry_models.go`) alongside the existing `typesafe/jev-1.13`, tagged `systemone` in the provider kind map, and mirrored them in the dashboard catalogue (`web/src/lib/models.ts`). Routing is unchanged: OpenRouter already has `SystemoneURL` (`https://openrouter.ai/api/v1/systemone`), so the `/v1/systemone` lane forwards these models there with the model id as given.
- Note for the experiment: OpenRouter documents the rolling alias with a tilde prefix (`~typesafe/jev-latest`); this change adds the id exactly as requested (`typesafe/jev-latest`). If a live call rejects that form, the documented tilde form is the fallback.
- Verified: `GetProviderModelKind` table cases for both new models pass; `gofmt` clean. No provider connections were checked or tested live.
