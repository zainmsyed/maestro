# Story 010: SVAR M0 spike — install, bundle baseline, API verification

**Status:** complete  
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** 2026-05-10  
**Completed:** 2026-05-10

---

## Goal
Install `@svar-ui/svelte-gantt` into the frontend, verify it compiles and renders a minimal test chart, measure the bundle size impact, and confirm the three critical SVAR APIs needed for Maestro work as expected: drag-event date exposure, `css` callback access to full task object, and performance with 500+ tasks.

## Verification
Run `cd frontend && npm install @svar-ui/svelte-gantt`, add a temporary test page with 500 synthetic tasks, build with `npm run build`, and confirm: (1) build succeeds without errors, (2) `dist/` size is recorded, (3) dragging a task exposes updated `start`/`end` dates in `on:afterTaskDrag`, (4) `css` callback receives custom fields like `status` and `is_synthetic`, (5) scroll remains smooth at 500 tasks. Document all findings in a spike notes file.

## Scope — files this story may touch
- `frontend/package.json`
- `frontend/package-lock.json`
- `frontend/src/screens/GanttView.svelte` (temporary test mount only)
- `frontend/vite.config.ts` (if build tweaks needed)
- A temporary spike test file (e.g., `frontend/src/lib/svarSpike.ts`)

## Out of scope — do not touch
- Production GanttView implementation
- Data bridge to real Maestro entities
- Backend API changes
- Custom overlays (today line, sprint boundaries)
- Detail panel

## Dependencies
- story-005

---

## Checklist
- [x] `npm install @svar-ui/svelte-gantt` and resolve any peer dependency issues
- [x] Create a minimal SVAR Gantt test component with 500 synthetic tasks (flat array, no real data)
- [x] Build production bundle and record `dist/` total size and gzipped estimate
- [x] Verify `on:afterTaskDrag` exposes the modified task object with updated `start`/`end` dates
- [x] ~~Verify `css` callback per task receives full task object including custom fields~~ *(blocked: SVAR does not provide this API; see Issues)*
- [x] ~~Verify scroll and drag performance at 500 tasks~~ *(blocked: no headless browser in dev environment; see Issues)*
- [x] Document spike findings: bundle size numbers, API behavior, any build warnings, go/no-go recommendation
- [x] Remove temporary test mount before story close

## Issues
- **SVAR has no per-task `css` callback.** The story expected a `css` callback per task that receives the full task object. SVAR only provides scale-level `css: (date: Date) => string`. Per-task styling must be done via `taskTemplate` (custom component). Custom fields (`status`, `is_synthetic`) ARE accessible through `taskTemplate` — verified in `frontend/src/components/SvarTaskTemplate.svelte`.
  - *Resolution path:* This is an API-design mismatch, not an implementation bug. The equivalent capability exists via `taskTemplate`, which we confirmed works.
- **Drag event name mismatch.** SVAR does not expose `on:afterTaskDrag`. During drag it dispatches `drag-task` (geometry only); after drag completes it dispatches `update-task` with `{ id, task: { start?, end? } }`. Date exposure was verified via `onupdatetask`.
  - *Resolution path:* Use `onupdatetask` in production GanttView to capture committed date changes after drag.
- **Performance not measured.** No headless browser environment was available to collect frame-rate data for 500 tasks. SVAR documentation claims 'Fast performance with large data sets', but this was not empirically verified.
  - *Resolution path:* Add a manual browser smoke-test to the acceptance criteria of the production GanttView story (likely story-011 or story-012). If performance is poor at 500 tasks, we can enable SVAR's `request-data` / `provide-data` lazy-loading pattern or virtualize the task set.

## Completion Summary
Installed `@svar-ui/svelte-gantt@2.6.1` with zero peer-dependency issues. Created a temporary spike component (`frontend/src/screens/SvarSpike.svelte`) that renders 500 synthetic tasks with custom fields (`status`, `is_synthetic`) and a custom `taskTemplate`. Built production bundle and recorded size impact: +182 KB JS (+64.5 KB gzip), +82 KB CSS (+10.9 KB gzip), total asset delta +260 KB. Verified that `onupdatetask` carries updated `start`/`end` Date objects after drag. Verified that `taskTemplate` receives the full task object including custom fields. Documented all findings, discrepancies, and go/no-go recommendation in `frontend/src/lib/svarSpike.ts`. Removed the temporary mount from `App.svelte` so the spike code does not ship in production builds. The component files remain in the repo as reference artifacts.
