### CodeBuddy: send the plugin request/conversation id headers (upstream decolua/9router@3fd925a)

- **Chat requests to `codebuddy-cn` and `codebuddy-intl` now carry the plugin
  id family exactly as the official plugin sends it.** `codebuddyHeaders`
  (`internal/proxy/executor/codebuddy.go`) issues fresh ids per request:
  `X-Request-ID`, `X-Conversation-Request-ID` and `X-Conversation-Message-ID`
  all reuse the same request id, while `X-Conversation-ID` is a fresh UUID
  per request. Previously the Go executor minted an independent
  `X-Conversation-Request-ID` and pinned `X-Conversation-ID` to the session
  via a stable v5 derivation (#581), both diverging from the contract the
  upstream channel check verifies; a request without the plugin ids is
  rejected upstream with 400 / 11128 "request illegal".
- **Deliberate Go-format note:** request, conversation-request and message
  ids stay 32-hex (no dashes) and the conversation id stays a dashed UUID,
  matching the previous Go wire format and its tests. The upstream fix
  asserts presence, freshness and the request-id equality, not the dash
  format, so the format is preserved rather than churned.
- **Removed dead derivation:** `stableConversationUUID` and `firstUserText`
  are deleted with the session pinning they served; no other caller used
  them. Tests updated to lock the new contract: fresh ids per request,
  request-id equality, conversation-id differing from message-id, and the
  session id no longer pinning the conversation id.
- **Already present, not re-done:** the backoff/combo half of 3fd925a was
  already ported in Go (#5024): `codebuddy-intl` "request illegal" backs
  off, the same text stays request-scoped for other providers, malformed
  CodeBuddy bodies classify request-scoped, and the combo lock path skips
  request-scoped 4xx without locking the account, so a combo advances past
  a malformed member instead of stopping on it.
- Verified: `go test ./internal/proxy/executor/ -run TestCodebuddy` pass,
  `go test ./internal/providers/ -run TestClassifyErrorForProvider` pass,
  `go build ./...` pass, `go vet` on both packages pass. No provider
  connections were checked.
