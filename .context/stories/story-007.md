# Story 007: Onboarding flow

**Status:** complete  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-07  
**Completed:** 2026-05-07

---

## Goal
Build the multi-step onboarding screen: file drop/picker, import preview, sprint review table, manual sprint config fallback, and confirmation before entering the app.

## Verification
Launch the app with an empty database. Drop a test CSV onto the onboarding screen, proceed through sprint review, click "build roadmap," and land on the List view with imported data visible.

## Scope — files this story may touch
- `frontend/src/screens/Onboarding.svelte`
- `frontend/src/components/DropZone.svelte`
- `frontend/src/components/SprintPreview.svelte`
- `frontend/src/components/StepBar.svelte`
- `frontend/src/lib/api.ts`
- `frontend/src/stores/project.ts`

## Out of scope — do not touch
- List view implementation (just navigate to it after onboarding)
- Re-import logic from inside the app (Settings → Import comes later)
- Gantt or Health screens

## Dependencies
- story-004
- story-005

---

## Checklist
- [x] Create `Onboarding.svelte` with 3-step state machine (Import → Sprints → Confirm)
- [x] Build `DropZone.svelte` with drag-and-drop and file picker fallback
- [x] POST file to `/api/import` and display preview counts (epics, features, sprints, missing dates, orphans)
- [x] Build `SprintPreview.svelte` table showing imported/generated sprints with editable start dates
- [x] Add manual config toggle with sprint length, first start date, and count inputs
- [x] Show confirmation summary card before finalizing
- [x] Wire "build roadmap" to create project and route to List view
- [x] Handle import errors with inline messages (malformed rows, unsupported format)

---

## Issues

---

## Completion Summary
Implemented the onboarding flow: a three-step Import → Sprints → Confirm screen with CSV drag/drop and picker fallback, import preview counts and inline import warnings, sprint review with editable starts, manual sprint generation controls, confirmation summary, project-store creation, and List-view routing via the view store. Added scoped components `DropZone.svelte`, `SprintPreview.svelte`, and `StepBar.svelte`, plus API/project-store types and helpers for CSV import and onboarding project state. Wired `App.svelte` to mount onboarding for the empty-roadmap state and from the topbar Import CSV action, replacing the temporary DateAssignment test mount.
