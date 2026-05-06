# Story 010: Gantt timeline core rendering

**Status:** not-started
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Render the primary Gantt view: sprint headers on the X-axis, epic rows with expandable feature rows on the Y-axis, bars positioned by sprint dates with status colors, today line, and the synthetic "Unassigned" Epic at the bottom.

## Verification
Navigate to the Timeline view. Verify all epics and features render as bars aligned to the correct sprint columns. Expand/collapse an epic. Confirm the "Unassigned" epic appears at the bottom in muted styling. Confirm the today line is visible.

## Scope — files this story may touch
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/components/GanttGrid.svelte`
- `frontend/src/components/RowLabels.svelte`
- `frontend/src/components/SprintHeader.svelte`
- `frontend/src/components/GanttBar.svelte`
- `frontend/src/components/TodayLine.svelte`
- `frontend/src/lib/ganttLayout.ts`

## Out of scope — do not touch
- Drag interactions (click only for now)
- Detail panel slide-out
- Virtualization (render all rows)
- Zoom controls

## Dependencies
- story-004
- story-005

---

## Checklist
- [ ] Create `GanttView.svelte` screen shell with topbar view-tabs integration
- [ ] Build `SprintHeader.svelte` row showing sprint names and date ranges
- [ ] Build `RowLabels.svelte` panel with epic rows and nested feature rows, expand/collapse chevrons
- [ ] Implement `GanttGrid.svelte` with vertical grid lines per sprint
- [ ] Create `GanttBar.svelte` for epics and `FeatureBar.svelte` for features
- [ ] Position bars using sprint start date fallback chain (sprint start → sprint end → import date)
- [ ] Apply status colors: grey (not started), blue (in progress), green (completed on time), red (missed)
- [ ] Render synthetic "Unassigned" Epic at bottom with muted styling (italic, no color pill)
- [ ] Add `TodayLine.svelte` vertical marker with "today" pip
- [ ] Sync scroll between `RowLabels` and `GanttGrid` containers
- [ ] Add zoom level state (sprint / month / quarter) with basic CSS scaling

---

## Issues

---

## Completion Summary
