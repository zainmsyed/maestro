# Story 008: Gantt drag interactions and detail panel

**Status:** not-started  
**Type:** ui  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

---

## Goal
Make the Gantt interactive: drag bars to shift dates, drag edges to extend end dates, click bars to open a detail panel showing original/committed dates, slip events, audit log, and a snap-to-sprint toggle.

## Verification
In the Gantt view, drag an epic bar to the right and verify the committed date updates via API. Drag the right edge of a feature bar to extend it. Click a bar and verify the detail panel slides in showing original date (strikethrough), committed date, status, and audit history. Toggle snap-to-sprint and drag again to verify snapping behavior.

## Scope — files this story may touch
- `frontend/src/components/GanttBar.svelte`
- `frontend/src/components/FeatureBar.svelte`
- `frontend/src/components/DetailPanel.svelte`
- `frontend/src/components/SnapToggle.svelte`
- `frontend/src/lib/dragEngine.ts`
- `frontend/src/lib/api.ts`
- `frontend/src/stores/ui.ts`

## Out of scope — do not touch
- Virtualization (still render all rows)
- Settings screen for snap toggle default (just local state for now)
- Health dashboard

## Dependencies
- story-007

---

## Checklist
- [ ] Implement drag engine in `dragEngine.ts`: pointer down → track delta → pointer up → compute new date
- [ ] Add drag-to-shift on bar body (updates `committed_end_date` only)
- [ ] Add drag-to-resize on bar right edge (extends `committed_end_date`)
- [ ] Build `DetailPanel.svelte` slide-out overlay with fixed positioning
- [ ] Detail panel fields: type, title, sprint, original date (strikethrough), committed date, status, slip delta
- [ ] Show date audit log entries in detail panel (newest first)
- [ ] Build `SnapToggle.svelte` component (default: on)
- [ ] When snap is enabled, round dragged dates to nearest sprint boundary
- [ ] Call PATCH endpoints on drag end and update local store optimistically
- [ ] Handle drag cancellation (Esc key or mouse leave)

---

## Issues

---

## Completion Summary
