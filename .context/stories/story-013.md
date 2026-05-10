# Story 013: SVAR drag handler + snap-to-sprint

**Status:** not-started
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** —
**Completed:** —

---

## Goal
Make the SVAR Gantt interactive: dragging a Feature bar updates its `committed_end_date` via the backend API. Implement a snap-to-sprint toggle (default: on) that rounds dragged dates to the nearest sprint boundary. Preserve the rule that `original_end_date` is never mutated.

## Verification
In the Gantt view, drag a Feature bar to the right and release. Verify the bar snaps to the nearest sprint boundary if the toggle is on, or lands at the free-drag position if off. Refresh the page and confirm the new `committed_end_date` persisted. Check the Date Audit Log and confirm a new entry exists with the correct `delta_days`.

## Scope — files this story may touch
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/lib/svarBridge.ts` (update helpers)
- `frontend/src/components/SnapToggle.svelte` (new)
- `frontend/src/lib/api.ts` (PATCH helpers)
- `frontend/src/stores/ui.ts` (snap toggle state)

## Out of scope — do not touch
- Detail panel (story-014)
- Epic-level drag (Epics are summary rows; only Features are draggable tasks)
- Settings screen for default snap preference
- List view, Health dashboard

## Dependencies
- story-011
- story-012

---

## Checklist
- [ ] Add `on:afterTaskDrag` handler to `GanttView.svelte`
- [ ] Extract new `start` and `end` dates from the dragged task object
- [ ] Build `SnapToggle.svelte` component with on/off state (default: on)
- [ ] When snap is enabled: round `end` to the nearest sprint's `end_date`; optionally round `start` to sprint `start_date`
- [ ] Call `PATCH /api/features/:id/date` with the new `committed_end_date`
- [ ] Optimistically update local task array so the bar stays at the dragged position
- [ ] Handle API errors by reverting the bar to its pre-drag position
- [ ] Ensure `original_end_date` is never sent in the PATCH payload
- [ ] Verify audit log entry is created with correct `delta_days`
- [ ] Handle drag cancellation (Esc key resets bar position)

## Issues

## Completion Summary
