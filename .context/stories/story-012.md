# Story 012: SVAR styling + today line + sprint boundary overlays

**Status:** complete  
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** 2026-05-10  
**Completed:** 2026-05-10

---

## Goal
Apply Maestro design system styling to the SVAR Gantt: status color coding on bars, muted styling for the synthetic "Unassigned" Epic, a vertical today line overlay, and thin vertical sprint boundary dividers on the day-based timeline.

## Verification
Open the Gantt view. Confirm Epic and Feature bars show the correct colors: grey (not started), blue (in progress), green (done), red (missed). Confirm the "Unassigned" Epic row and its Features are rendered in muted grey/italic styling. Confirm a vertical today line is visible at the current date. Confirm thin vertical lines mark each sprint's start boundary.

## Scope — files this story may touch
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/lib/svarBridge.ts` (add `css` and detail fields)
- `frontend/src/components/TodayLine.svelte` (new)
- `frontend/src/components/SprintBoundaryOverlay.svelte` (new)
- `frontend/src/app.css` (overlay styles if needed)

## Out of scope — do not touch
- Drag interactions (story-013)
- Detail panel (story-014)
- List view, Health dashboard, Settings
- Backend API changes

## Dependencies
- story-011

---

## Checklist
- [x] Implement `css` callback in `svarBridge.ts` that returns dynamic class strings based on `status` and `is_synthetic`
  - `status = not_started` → grey bar
  - `status = in_progress` → blue bar
  - `status = done` + on time → green bar
  - `status = done` + missed → red bar
  - `is_synthetic = true` → muted (grey, reduced opacity, italic label)
- [x] Add CSS classes for each status variant to Gantt view scoped styles
- [x] Build `TodayLine.svelte`: thin vertical div positioned via date-to-pixel math against SVAR's timeline scale
- [x] Build `SprintBoundaryOverlay.svelte`: thin vertical divs at each sprint's `start_date`, positioned via date-to-pixel math
- [x] Ensure overlays sit above the Gantt grid but do not intercept pointer events
- [x] Verify muted styling on synthetic Unassigned Epic and its child Features
- [x] Verify status colors match PRD §5.2 exactly
- [x] Verify today line and sprint boundaries are visible at all zoom levels

## Issues

## Completion Summary
Implemented Gantt status styling and timeline overlays. `svarBridge.ts` now derives `status_variant`, `is_missed`, synthetic inheritance, and `css` class strings for each task. `GanttView.svelte` applies status colors via generated SVAR bar selectors: grey for not started, blue for in progress, green for completed on time, red for missed, and muted grey/italic styling for synthetic rows. Added `TodayLine.svelte` and `SprintBoundaryOverlay.svelte` with date-to-pixel positioning and pointer-events disabled. Updated the Gantt timeline to day scale and added sprint-start boundary extraction from epic sprint starts. Added bridge tests for CSS/status mapping, missed deadlines, and synthetic child styling. Verification performed: `npm run check` passes with 0 warnings, `npm test` passes 31 tests, and `npm run build` succeeds.
