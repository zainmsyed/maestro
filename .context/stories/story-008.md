# Story 008: List view with inline editing and CSV export

**Status:** completed  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-07  
**Completed:** 2026-05-07

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
- story-004
- story-005

---

## Checklist
- [x] Create `ListView.svelte` screen shell with topbar integration
- [x] Build `DataTable.svelte` with fixed header and scrollable body
- [x] Columns: Title, Type, Owner, Sprint, Original Date, Committed Date, Actual Date, Slip Events, Status, Health indicator
- [x] Implement inline date editing in Committed Date cells via `DatePicker.svelte`
- [x] Add sort toggles on every column header
- [x] Build `FilterBar.svelte` with dropdowns for Epic, Owner, Sprint, Status
- [x] Implement group-by: Epic (default), Sprint, Owner
- [x] Show orphaned features under "Unassigned" group with reassignment dropdown
- [x] Add "Export CSV" button that exports the current filtered/sorted/grouped view
- [x] Handle empty states and loading states

---

## Issues

---

## Completion Summary
Implemented the scoped List view: `ListView.svelte`, `DataTable.svelte`, `FilterBar.svelte`, `GroupBySelect.svelte`, and `csvExport.ts` now load epics plus features, flatten and group them by Epic/Sprint/Owner, sort on every column, filter by Epic/Owner/Sprint/Status, export the current grouped view to CSV, and show loading/empty/error states. Feature rows support inline committed-date editing through the shared `DatePicker.svelte`, and orphaned features under the synthetic unassigned epic can be reassigned inline through the existing feature epic PATCH endpoint. With explicit approval for a scope exception, `App.svelte` was also wired to render `ListView` when the List route/tab is active so story-008 can be visually verified end-to-end.
