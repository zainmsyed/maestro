# Story 011: Update Gantt for three-level nesting and stub bars

**Status:** in-progress  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-08  
**Completed:** —

---

## Goal
Update the Gantt view to render three-level nesting (Epic → Feature → Story), show stub bars for undated items, make Stories collapsed by default, and handle click-to-assign on stub bars.

## Verification
Open the Gantt view with test data containing epics, features, and stories at various nesting levels. Verify Stories are collapsed by default, expanding a Feature shows its Stories, undated Stories render as stub bars with ⚠ indicator, clicking a stub opens the detail panel with date picker, and dated Stories render as thin bars (8px).

## Scope — files this story may touch
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/components/GanttGrid.svelte`
- `frontend/src/components/RowLabels.svelte`
- `frontend/src/components/GanttBar.svelte`
- `frontend/src/components/FeatureBar.svelte`
- `frontend/src/components/StoryBar.svelte` (new)
- `frontend/src/components/StubBar.svelte` (new)
- `frontend/src/components/DetailPanel.svelte`
- `frontend/src/lib/ganttLayout.ts`

## Out of scope — do not touch
- List view
- Health dashboard
- Settings screens
- Onboarding flow
- Drag engine (story-012 handles drag; this story focuses on rendering)

## Dependencies
- story-004
- story-010

---

## Checklist
- [ ] Update `RowLabels.svelte` to support three levels: Epic rows, Feature rows (indented under epics), Story rows (indented under features)
- [ ] Stories default to collapsed under Features; expand/collapse chevron on Feature rows
- [ ] Update `GanttGrid.svelte` to render Story rows with thin bars (8px height)
- [ ] Create `StoryBar.svelte` for dated Stories
- [ ] Create `StubBar.svelte` for undated Stories: minimal bar at sprint start position with ⚠ indicator and lighter opacity
- [ ] Clicking a stub bar opens `DetailPanel.svelte` with date picker
- [ ] After assigning a date, stub becomes a full `StoryBar`
- [ ] Story bars are not draggable until a date is assigned
- [ ] Update `GanttBar.svelte` for Epics (22px) and `FeatureBar.svelte` for Features (10px)
- [ ] Update `ganttLayout.ts` to handle three-level row heights and offsets
- [ ] Update scroll sync and virtualization for three levels
- [ ] Verify synthetic unassigned epic and feature rows render in muted styling at the bottom

---

## Issues

---

## Completion Summary
