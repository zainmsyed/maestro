# Story 005: Onboarding flow

**Status:** not-started  
**Type:** ui  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

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
- story-003
- story-004

---

## Checklist
- [ ] Create `Onboarding.svelte` with 3-step state machine (Import → Sprints → Confirm)
- [ ] Build `DropZone.svelte` with drag-and-drop and file picker fallback
- [ ] POST file to `/api/import` and display preview counts (epics, features, sprints, missing dates, orphans)
- [ ] Build `SprintPreview.svelte` table showing imported/generated sprints with editable start dates
- [ ] Add manual config toggle with sprint length, first start date, and count inputs
- [ ] Show confirmation summary card before finalizing
- [ ] Wire "build roadmap" to create project and route to List view
- [ ] Handle import errors with inline messages (malformed rows, unsupported format)

---

## Issues

---

## Completion Summary
