### SenseNova strict reasoning history

- Strip replayed assistant reasoning carriers before forwarding to SenseNova's strict Chat Completions schema, including custom nodes pointed at `token.sensenova.ai` or `token.sensenova.cn`.
- This prevents the generic `400 inference request is invalid` seen on multi-turn requests while preserving the assistant content and leaving other providers unchanged.
- Verified with `go test ./internal/proxy/executor`.
