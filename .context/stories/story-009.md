# Story 009: Gantt virtualization for large roadmaps

**Status:** not-started  
**Type:** ui  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

---

## Goal
Ensure the Gantt view renders smoothly with 500+ items by virtualizing both the row labels panel and the gantt grid so only visible rows are in the DOM.

## Verification
Generate a test dataset with 500+ features across 20 epics. Open the Gantt view, scroll vertically, and drag a bar. Use browser devtools Performance tab to confirm no dropped frames during scroll or drag.

## Scope — files this story may touch
- `frontend/src/components/RowLabels.svelte`
- `frontend/src/components/GanttGrid.svelte`
- `frontend/src/lib/virtualList.ts`
- `frontend/src/lib/ganttLayout.ts`
- `cmd/maestro/main.go` (test data generator endpoint, optional)

## Out of scope — do not touch
- New Gantt features or styling changes
- Health dashboard
- List view

## Dependencies
- story-007
- story-008

---

## Checklist
- [ ] Measure baseline performance with 500 items (non-virtualized) to document the problem
- [ ] Implement `virtualList.ts` helper that calculates visible row range from scroll position and row height
- [ ] Update `RowLabels.svelte` to render only visible rows using absolute positioning
- [ ] Update `GanttGrid.svelte` to render only visible gantt rows and bars
- [ ] Maintain scroll sync between virtualized row labels and grid
- [ ] Ensure expand/collapse still works correctly with virtual offsets
- [ ] Ensure drag interactions work on virtualized rows
- [ ] Test with 500-item dataset and confirm smooth 60fps scroll
- [ ] Document row height assumptions and any virtualization limitations

---

## Issues

---

## Completion Summary
