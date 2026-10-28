### Fixed — Nous free catalog matches the live endpoint

- Drop `meituan/longcat-2.5-preview:free`, `stepfun/step-5-preview:free`, and `stealth/missingno`: gone from `GET /v1/models` (probed 2026-10-28), so offering them only produces upstream 404s.
- Add `reflection/beam-1.0-501b-a23b:free` (reasoning mandatory, cannot disable) and `ornith/ornith-1.5-122b:free` (vision input, 262144 max output): both zero-priced with >= 200k context on the live catalogue.
- Verified with `go build ./...` and `go test ./internal/providers/ ./internal/pricing/ ./internal/handlers/chat/`.
