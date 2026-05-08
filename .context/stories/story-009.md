# Story 009: Update List view for three-level hierarchy and date assignment

**Status:** complete  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-08  
**Completed:** 2026-05-08

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
- [x] Update `DataTable.svelte` to render three levels: Epic rows (expandable), Feature rows (expandable under epics), Story rows (under features)
- [x] Add `DateSourceBadge.svelte` component showing `imported` / `pm assigned` / `inherited` as small colored badges
- [x] Story rows show: Title, Sprint, Original Date, Committed Date, Actual Date, Status, Date Source, Story Points
- [x] Undated items (NULL `committed_end_date`) show "— set date" in Committed Date column; clicking opens inline date picker
- [x] After assigning a date via inline picker, the row updates to show the new date and `date_source = 'pm_assigned'`
- [x] Group by: Epic (default), Sprint, Owner — all three levels visible in each group
- [x] Filter by: Epic, Owner, Sprint, Status, Type (Epic/Feature/Story)
- [x] Update CSV export to include all three levels and date_source column
- [x] Handle empty states and loading states
- [x] Verify orphaned stories appear under synthetic unassigned feature under synthetic unassigned epic

---

## Issues

---

## Completion Summary

### Data layer (`frontend/src/lib/listView.ts`)
- Extended `ListRow` type with `story` variant, plus `dateSource`, `storyPoints`, `featureId`, and `featureTitle` fields.
- Updated `flattenRows` to recursively flatten stories under their parent features.
- Added `type` to `FilterState` and updated `filterRows` predicate.
- Rewrote `groupByEpic` to maintain three-level ordering: epic → sorted features → sorted stories under each feature.
- Updated `buildExportRows` to emit `date_source` and `story_points` columns.

### API client (`frontend/src/lib/api.ts`)
- Added `patchStoryFeature` helper for reassigning stories to features.

### New component (`frontend/src/components/DateSourceBadge.svelte`)
- Small colored badge component mapping `imported` → blue, `pm_assigned` → green, `inherited` → purple, `manual` → amber.

### DataTable (`frontend/src/components/DataTable.svelte`)
- Three-level expandable hierarchy with ▼/▶ toggles on epic and feature rows.
- Indentation classes: feature (28px), story (52px).
- Added Story Points and Date Source columns.
- Inline date picker now works for both features and stories; undated items show "— set date".
- Type pills get distinct color tokens (epic=blue, feature=purple, story=green).
- Reassignment dropdowns appear for synthetic/orphaned features and stories.

### FilterBar (`frontend/src/components/FilterBar.svelte`)
- Added Type filter dropdown (All / Epic / Feature / Story).
- Updated grid to 6 columns; adjusted responsive breakpoint.

### ListView (`frontend/src/screens/ListView.svelte`)
- Header summary now shows 4 metrics: Epics, Features, Stories, Groups.
- `saveDate` handler dispatches to the correct entity type via `api.patchDate(row.type, …)`.
- Added `reassignStory` handler wired to `api.patchStoryFeature`.
- Cascading filter options include the new `type` dimension.

### Tests
- Updated `listView.test.ts` with story fixtures, type-filter test, three-level grouping test, and export assertion for `date_source`.
- All 21 frontend tests pass; Go backend tests pass; production build succeeds.
