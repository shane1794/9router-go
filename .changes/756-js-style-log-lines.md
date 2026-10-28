### Changed — Request log lines match the Next.js gateway shape

- Request dispatch now emits one correlated line: `<tag> ▶ POST <client-model> → <provider>/<model> · STREAM|JSON · ACC:<conn>` with a stable per-connection emoji tag (`log.TagForSession`, same 8-color rotation as `open-sse` `tagForSession`).
- RTK line now reads `[RTK] saved <n>B / <n>B (<pct>%) via [<techniques>] hits=<k>` (`open-sse/rtk/index.js formatRtkLog` shape) with the distinct `TechniquesUsed` from `CompressTextDetailedWithCommand`; token estimates stay behind as kv.
- Usage line now reads `📊 DONE <lat>ms · [TTFT <n>ms ·] IN <n> [(CACHE ↻r +c)] · OUT <n>` (`requestDetail.js formatDoneLine` shape) with the full kv pairs preserved for parsers.
- Text log lines carry `[HH:MM:SS]` and the stdlib date/time prefix is disabled so stamps do not duplicate. JSON log shape unchanged.
