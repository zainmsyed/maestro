# Story 010: Gantt timeline core rendering

**Status:** complete  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-08  
**Completed:** 2026-05-08

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

## Approved scope exceptions
- `frontend/src/App.svelte` — user-approved integration so the existing Roadmap/Gantt tab renders the new Gantt screen and Health/Settings keep explicit placeholders.
- `frontend/src/screens/ListView.svelte` — user-requested full-width cleanup to match the Gantt layout refinement.
- `frontend/src/lib/ganttLayout.test.ts` — added during review remediation to cover core Gantt layout behavior.

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
- [x] Create `GanttView.svelte` screen shell with topbar view-tabs integration
- [x] Build `SprintHeader.svelte` row showing sprint names and date ranges
- [x] Build `RowLabels.svelte` panel with epic rows and nested feature rows, expand/collapse chevrons
- [x] Implement `GanttGrid.svelte` with vertical grid lines per sprint
- [x] Create `GanttBar.svelte` variants for epic and feature bars
- [x] Position bars using sprint start date fallback chain (sprint start → sprint end → import date)
- [x] Apply status colors: grey (not started), blue (in progress), green (completed on time), red (missed)
- [x] Render synthetic "Unassigned" Epic at bottom with muted styling (italic, no color pill)
- [x] Add `TodayLine.svelte` vertical marker with "today" pip
- [x] Sync scroll between `RowLabels` and `GanttGrid` containers
- [x] Add zoom level state (sprint / month / quarter) with basic CSS scaling

---

## Issues

No open implementation blockers. Notes:
- `GanttBar.svelte` implements both epic and feature bar variants rather than adding a separate `FeatureBar.svelte` file.
- `frontend/src/App.svelte` was touched with user approval so the Gantt screen is reachable from the existing Roadmap/Gantt tab.
- Vertical scrolling was adjusted after review feedback to use the main page scrollbar; the grid keeps horizontal scrolling for wide timelines.
- `frontend/src/screens/ListView.svelte` was also adjusted to full width by explicit user request, outside the original story scope.

---

## Completion Summary
Implemented the primary Timeline/Gantt rendering path. The Roadmap/Gantt tab now loads `GanttView.svelte`, fetches epics from `/api/epics`, builds sprint columns, and renders epic rows with expandable feature rows. `SprintHeader.svelte`, `RowLabels.svelte`, `GanttGrid.svelte`, `GanttBar.svelte`, and `TodayLine.svelte` provide the timeline header, hierarchy labels, grid lines, status-colored bars, and today marker.

`frontend/src/lib/ganttLayout.ts` computes sprint columns, row order, visible rows, bar positions, synthetic epic ordering, status tones, and today-line placement. The synthetic "Unassigned" epic is sorted to the bottom and styled as muted/italic without the type color pill. The Gantt UI was refined based on review feedback to use full available width, avoid a nested vertical scrollbar, increase row spacing, and mirror the List view's hierarchy cues with gold epic bands and colored type pills.

Verification performed: `cd frontend && npm run build` passes.
