# Story 012: SVAR styling + today line + sprint boundary overlays

**Status:** not-started
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** —
**Completed:** —

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
- [ ] Implement `css` callback in `svarBridge.ts` that returns dynamic class strings based on `status` and `is_synthetic`
  - `status = not_started` → grey bar
  - `status = in_progress` → blue bar
  - `status = done` + on time → green bar
  - `status = done` + missed → red bar
  - `is_synthetic = true` → muted (grey, reduced opacity, italic label)
- [ ] Add CSS classes for each status variant to Gantt view scoped styles
- [ ] Build `TodayLine.svelte`: thin vertical div positioned via date-to-pixel math against SVAR's timeline scale
- [ ] Build `SprintBoundaryOverlay.svelte`: thin vertical divs at each sprint's `start_date`, positioned via date-to-pixel math
- [ ] Ensure overlays sit above the Gantt grid but do not intercept pointer events
- [ ] Verify muted styling on synthetic Unassigned Epic and its child Features
- [ ] Verify status colors match PRD §5.2 exactly
- [ ] Verify today line and sprint boundaries are visible at all zoom levels

## Issues

## Completion Summary
