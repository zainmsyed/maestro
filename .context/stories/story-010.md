# Story 010: SVAR M0 spike — install, bundle baseline, API verification

**Status:** in-progress  
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** 2026-05-10  
**Completed:** —

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
- [ ] `npm install @svar-ui/svelte-gantt` and resolve any peer dependency issues
- [ ] Create a minimal SVAR Gantt test component with 500 synthetic tasks (flat array, no real data)
- [ ] Build production bundle and record `dist/` total size and gzipped estimate
- [ ] Verify `on:afterTaskDrag` exposes the modified task object with updated `start`/`end` dates
- [ ] Verify `css` callback per task receives full task object including custom fields
- [ ] Verify scroll and drag performance at 500 tasks (devtools Performance tab, no dropped frames)
- [ ] Document spike findings: bundle size numbers, API behavior, any build warnings, go/no-go recommendation
- [ ] Remove temporary test mount before story close

## Issues

## Completion Summary
