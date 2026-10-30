### Model context tags: show each model's max context where it is picked

- Ported from serenhope/9router v0.5.117-Custom ("Every model shows its real
  window"): a model now carries its context window next to its name, so a 1M
  model is never picked as if it were a small one. The short form matches the
  source: 1.0M, 256k, 200, with the full token count in the tooltip.
- Shown in three places: the model picker chips (`ModelPill`, via the static
  catalog window), the combo member rows in the combo editor (catalog lookup
  by `provider/model`, first slash splits provider from a namespaced model
  id), and the Available Models grid on each provider page (caps endpoint
  wins, static catalog is the fallback).
- No backend change: the picker and combo rows read the static catalog, the
  provider grid reads the already-served `GET /api/models/caps` window first.
  A model with no known window shows no tag, never a guess.
- Coverage note: only ~77 static catalog entries carry a `contextLength`,
  so picker/combo tags appear only where the catalog knows the window; the
  provider grid covers more models because the caps endpoint resolves synced
  catalogs and custom models too. Custom and alias models in the picker show
  no tag until a caps source is wired into the picker data flow.
- Verified: web unit tests 658 pass, 0 fail (`src/`); `bun run build` pass;
  svelte-check 77 errors, same as baseline, none in the changed lines. The
  32 e2e failures in a full `bun test` run are the harness needing a built
  gateway binary, not this change. No provider connections checked.
- Files: `web/src/lib/formatContextWindow.ts` (+ test),
  `web/src/lib/models.ts` (+ test, `getCatalogContextLength`),
  `web/src/components/combos/pickerData.ts`, `ModelPill.svelte`,
  `ModelPickerModal.svelte`, `CreateComboModal.svelte`,
  `web/src/components/connections/types.ts`, `ProviderDetailView.svelte`.
