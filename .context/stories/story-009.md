# Story 009: Update List view for three-level hierarchy and date assignment

**Status:** in-progress  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-08  
**Completed:** —

---

## Goal
Update the List view to support three-level grouping (Epic → Feature → Story), display date source badges, show "— set date" for undated items with click-to-assign, and add Story-specific columns.

## Verification
Open the List view with test data containing epics, features, and stories (including undated items). Verify three-level grouping works, date source badges appear correctly, undated items show "— set date" and open a date picker when clicked, and CSV export includes all three levels.

## Scope — files this story may touch
- `frontend/src/screens/ListView.svelte`
- `frontend/src/components/DataTable.svelte`
- `frontend/src/components/DatePicker.svelte`
- `frontend/src/components/DateSourceBadge.svelte` (new)
- `frontend/src/lib/csvExport.ts`
- `frontend/src/lib/api.ts`

## Out of scope — do not touch
- Onboarding flow
- Gantt view
- Health dashboard
- Settings screens
- Import parser

## Dependencies
- story-004
- story-008

---

## Checklist
- [ ] Update `DataTable.svelte` to render three levels: Epic rows (expandable), Feature rows (expandable under epics), Story rows (under features)
- [ ] Add `DateSourceBadge.svelte` component showing `imported` / `pm assigned` / `inherited` as small colored badges
- [ ] Story rows show: Title, Sprint, Original Date, Committed Date, Actual Date, Status, Date Source, Story Points
- [ ] Undated items (NULL `committed_end_date`) show "— set date" in Committed Date column; clicking opens inline date picker
- [ ] After assigning a date via inline picker, the row updates to show the new date and `date_source = 'pm_assigned'`
- [ ] Group by: Epic (default), Sprint, Owner — all three levels visible in each group
- [ ] Filter by: Epic, Owner, Sprint, Status, Type (Epic/Feature/Story)
- [ ] Update CSV export to include all three levels and date_source column
- [ ] Handle empty states and loading states
- [ ] Verify orphaned stories appear under synthetic unassigned feature under synthetic unassigned epic

---

## Issues

---

## Completion Summary
