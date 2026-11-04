### Fix: OpenRouter Jev latest id is ~typesafe/jev-latest, not typesafe/jev-latest (2026-11-04)

- The experiment in fragment 941 added `typesafe/jev-latest`, but OpenRouter answers that id with HTTP 400 "Model typesafe/jev-latest does not exist". Per the OpenRouter model page for Jev Latest, the rolling alias is `~typesafe/jev-latest` (tilde prefix), and it "always redirects to the latest model in the Jev family". The catalogue entry, kind tag, dashboard entry, and test case are corrected to `~typesafe/jev-latest`; `typesafe/jev-1.14` is unchanged.
- Verified: `GetProviderModelKind` table case for `~typesafe/jev-latest` passes; `gofmt` clean. No provider connections were checked or tested live; the live call is left to the operator.
