# Story 017: Single-binary packaging

**Status:** not-started
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Produce a single Go binary with the frontend embedded, verify it stays under the 25 MB budget, starts in under 2 seconds, and can be cross-compiled for Windows, macOS, and Linux.

## Verification
Run `make build` and confirm `maestro` binary is < 25 MB. Execute it and measure time from process start to browser tab open < 2s. Run `make build-all` and confirm Windows, macOS, and Linux binaries are produced without error.

## Scope — files this story may touch
- `Makefile`
- `cmd/maestro/main.go` (embed directive, startup logic)
- `frontend/vite.config.ts` (build optimizations)
- `go.mod` (build tags if needed)
- `.gitignore`

## Out of scope — do not touch
- Application features or screens
- Database schema changes
- API endpoint changes

## Dependencies
- story-005
- story-016

---

## Checklist
- [ ] Add `//go:embed frontend/dist` directive in `cmd/maestro/main.go`
- [ ] Serve embedded files on `/` with proper MIME types and SPA fallback to `index.html`
- [ ] Create `Makefile` with targets: `build`, `build-all`, `test`, `clean`
- [ ] Configure Vite build for production: tree-shaking, minification, no source maps in embed
- [ ] Measure baseline binary size after `go build`
- [ ] If > 20 MB at M3 checkpoint, apply: lazy-loading for heavy routes, gzip on embed, review modernc.org/sqlite size
- [ ] Measure cold startup time from binary execution to first HTTP response
- [ ] Add cross-compilation for `GOOS=windows GOARCH=amd64`, `GOOS=darwin GOARCH=amd64`, `GOOS=linux GOARCH=amd64`
- [ ] Document build steps in README
- [ ] Final size check: confirm binary < 25 MB

---

## Issues

---

## Completion Summary
