# Story 004: Scaffold Vite + Svelte frontend

**Status:** not-started  
**Type:** ui  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

---

## Goal
Initialize a Vite + Svelte SPA project, wire it into the Go backend via `embed.FS`, and serve the built assets from the single binary. Establish the design token system and global layout shell.

## Verification
Run `cd frontend && npm run build`, then `go build -o maestro ./cmd/maestro` and execute `./maestro`. Open `http://localhost:8080` and see the app shell load with the correct dark theme, fonts (DM Mono, Fraunces), and sidebar navigation.

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
- API business logic (already built in story-003)
- Import parser

## Dependencies
- story-003

---

## Checklist
- [ ] Initialize Vite project with Svelte and TypeScript (`npm create vite@latest frontend -- --template svelte-ts`)
- [ ] Configure build output to `frontend/dist`
- [ ] Add DM Mono and Fraunces font imports in `index.html`
- [ ] Create CSS custom properties in `app.css` matching mockup tokens (bg, bg2, accent, text, etc.)
- [ ] Build `App.svelte` with global layout shell: sidebar, main area, status bar slots
- [ ] Create `api.ts` fetch wrapper with base URL and JSON helpers
- [ ] Set up basic Svelte stores for `project` and `view`
- [ ] Update `cmd/maestro/main.go` to serve `frontend/dist` via `embed.FS` on `/`
- [ ] Add dev proxy config so `npm run dev` can reach the Go API on `:8080`
- [ ] Verify `go build` produces a binary that serves the frontend correctly

---

## Issues

---

## Completion Summary
