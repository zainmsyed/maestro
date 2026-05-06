# Story 006: List view with inline editing and CSV export

**Status:** not-started  
**Type:** ui  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

---

## Goal
Build the List view as a rich data table: display epics and features, inline edit committed dates, sort and filter by any column, group by Epic/Sprint/Owner, reassign orphaned features, and export the current view to CSV.

## Verification
Open the List view with test data. Edit a feature's committed date inline and refresh to confirm it persisted. Sort by Owner, filter to Status = Active, group by Sprint, reassign an orphaned feature to a real epic, and export the filtered view to CSV.

## Scope — files this story may touch
- `frontend/src/screens/ListView.svelte`
- `frontend/src/components/DataTable.svelte`
- `frontend/src/components/DatePicker.svelte`
- `frontend/src/components/FilterBar.svelte`
- `frontend/src/components/GroupBySelect.svelte`
- `frontend/src/lib/csvExport.ts`
- `frontend/src/lib/api.ts`

## Out of scope — do not touch
- Gantt view
- Health dashboard
- Detail panel (used in Gantt; List uses inline cells)
- Settings screens

## Dependencies
- story-003
- story-004

---

## Checklist
- [ ] Create `ListView.svelte` screen shell with topbar integration
- [ ] Build `DataTable.svelte` with fixed header and scrollable body
- [ ] Columns: Title, Type, Owner, Sprint, Original Date, Committed Date, Actual Date, Slip Events, Status, Health indicator
- [ ] Implement inline date editing in Committed Date cells via `DatePicker.svelte`
- [ ] Add sort toggles on every column header
- [ ] Build `FilterBar.svelte` with dropdowns for Epic, Owner, Sprint, Status
- [ ] Implement group-by: Epic (default), Sprint, Owner
- [ ] Show orphaned features under "Unassigned" group with reassignment dropdown
- [ ] Add "Export CSV" button that exports the current filtered/sorted/grouped view
- [ ] Handle empty states and loading states

---

## Issues

---

## Completion Summary
