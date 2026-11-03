### Upstream sync: luqman-v1/9router-go to d6dd9481 (trial)

- Synced the 24 first-parent commits after `ae4b4149`, through `d6dd9481` (#895), including the provider-registry/catalog refactor, resilience admission control, service tiers, proxy-pool routing, HCNSec, OAuth context propagation, translator healing, success-rate analytics, and Gemini embedding dimensions.
- Resolved conflicts while preserving the fork's fresh CodeBuddy conversation IDs, `invalid_grant` OAuth parking, model context tags, and fork-only build behavior.
- Verified executor, translator, providers, and OAuth tests. Chat verification reaches the pre-existing refresh-slot test isolation failure. No provider connections were checked.
