# Fork-only: `make build` no longer compiles the `opencode-go` copy

- `Makefile` target `build` previously compiled the same `./cmd/9router-go/`
  package twice, once as `9router-go` and once as `opencode-go`. The second
  line is removed (fork-only change, 818 is fork-assigned, not an upstream
  PR number). Upstream CI, the Dockerfile and the release artifacts build
  only `9router-go`, nothing invokes the `opencode-go` binary, and the
  `opencode-go` *provider* (alias `ocg`) is routed inside the gateway
  regardless of the binary filename, so behaviour is unchanged.
- Verified: `make build` in a clean worktree produces `9router-go` version
  1.9.15 with the embedded version and no `opencode-go` file.
- Note for upgraders: a stale `opencode-go` file from an older build is
  not removed by `make clean` (it only removes `9router-go-*`); delete it
  once by hand (`del opencode-go.exe` on Windows, `rm opencode-go` elsewhere).
