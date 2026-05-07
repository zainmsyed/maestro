# Story 005: Scaffold Vite + Svelte frontend

**Status:** in-progress  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-07  
**Completed:** —

---

## Goal
Initialize a Vite + Svelte SPA project, wire it into the Go backend via `embed.FS`, and serve the built assets from the single binary. Establish the design token system and global layout shell.

## Verification
Run `cd frontend && npm run build`, then `go build -o maestro ./cmd/maestro` and execute `./maestro`. Open `http://localhost:9000` and see the app shell load with the correct dark theme, fonts (DM Mono, Fraunces), and sidebar navigation.

## Scope — files this story may touch
- `frontend/package.json`
- `frontend/vite.config.ts`
- `frontend/tsconfig.json`
- `frontend/index.html`
- `frontend/src/main.ts`
- `frontend/src/App.svelte`
- `frontend/src/app.css`
- `frontend/src/lib/api.ts`
- `frontend/src/stores/*.ts`
- `cmd/maestro/main.go` (static file serving)

## Out of scope — do not touch
- Actual screen components (onboarding, list, gantt, health, settings)
- API business logic (already built in story-004)
- Import parser

## Dependencies
- story-004

---

## Checklist
- [x] Initialize Vite project with Svelte and TypeScript (`npm create vite@latest frontend -- --template svelte-ts`)
- [x] Configure build output to `frontend/dist`
- [x] Add DM Mono and Fraunces font imports in `index.html`
- [x] Create CSS custom properties in `app.css` matching mockup tokens (bg, bg2, accent, text, etc.)
- [x] Build `App.svelte` with global layout shell: sidebar, main area, status bar slots
- [x] Create `api.ts` fetch wrapper with base URL and JSON helpers
- [x] Set up basic Svelte stores for `project` and `view`
- [ ] Update `cmd/maestro/main.go` to serve `frontend/dist` via `embed.FS` on `/`
- [x] Add dev proxy config so `npm run dev` can reach the Go API on `:9000`
- [ ] Verify `go build` produces a binary that serves the frontend correctly

---

## Issues
- The intake references folder currently contains `maestro-mockup.html` and `Example_FinancialDashboard_Backlog.csv`; the referenced PNG was not present after reverting to `pynytysw`.
- `cmd/maestro/main.go` now serves `frontend/dist` from disk with SPA fallback, but not via `embed.FS`. Go embed patterns cannot reference `../../frontend/dist` from `cmd/maestro`; satisfying true single-binary embedding needs an out-of-scope helper package/file (for example `frontend/embed.go`) or moving build output under `cmd/maestro`.
- Verification on default `:9000` is blocked because another process is listening on that port. The same binary was verified on `:18080` and served the frontend shell successfully.

---

## Completion Summary
Partial story-005 implementation is in place: a Vite + Svelte + TypeScript frontend shell was scaffolded, design tokens were established from the HTML mockup with roomy spacing, global app layout/sidebar/topbar/statusbar were added, API and store helpers were created, and Vite dev proxy targets the Go API on `:9000`. The Go entrypoint serves built frontend files from `frontend/dist` at `/` with API routes preserved under `/api/`.`

Remaining before `/complete-story`: approve and implement an out-of-scope embed helper or other structure for true `embed.FS` single-binary serving, and verify on default `:8080` after freeing the occupied port.
