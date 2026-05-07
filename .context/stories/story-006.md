# Story 006: Build post-import date assignment screen

**Status:** complete  
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-07  
**Completed:** 2026-05-07

---

## Goal
Build the post-import date assignment screen that appears after CSV import when items have missing Target Dates. The PM can assign dates inline with sprint-boundary snap, skip for later, or set all at once. The first assigned date becomes the locked `original_end_date`.

## Verification
Run an import with a fixture containing items missing Target Dates. After the sprint review step, verify the date assignment screen appears with the correct items listed. Assign dates to two items, skip one, click continue, and verify in the database that assigned items have `original_end_date` and `committed_end_date` set with `date_source = 'pm_assigned'`, and skipped items remain with NULL dates.

## Scope — files this story may touch
- `frontend/src/screens/DateAssignment.svelte` (new)
- `frontend/src/components/DateAssignmentRow.svelte` (new)
- `frontend/src/components/DatePicker.svelte`
- `frontend/src/screens/Onboarding.svelte`
- `frontend/src/lib/api.ts`
- `frontend/src/stores/project.ts`

## Out of scope — do not touch
- Import parser (already built)
- Sprint configuration screen (story-007)
- Gantt or List views
- Settings screens

## Dependencies
- story-003
- story-005

---

## Checklist
- [x] Create `DateAssignment.svelte` screen showing items missing Target Date with count header
- [x] Each row shows: checkbox (select for batch), item title, type (Epic/Feature/Story), sprint name, date picker
- [x] Date picker snaps to sprint boundaries by default (toggle per row)
- [x] "Set all selected" button assigns dates to all checked items
- [x] "Skip — I'll do it later" button continues to confirmation without assigning
- [x] "Set dates now" button assigns dates and continues
- [x] Call PATCH endpoints for each assigned item; first assignment sets both `original_end_date` and `committed_end_date`
- [x] Show progress/success state after batch assignment
- [x] Wire `Onboarding.svelte` to show `DateAssignment` between sprint review and confirmation when `missing_dates > 0`
- [x] Handle empty state (all items have dates → skip screen automatically)
- [x] Handle errors gracefully (network failure, invalid date)

---

## Issues
- Current backend has PATCH date endpoints for stories and features, but no `PATCH /api/epics/{id}/date` route. The frontend calls the expected endpoint for epic candidates, but epic date assignment will return 404 until backend epic date support is added.
- Sprint boundary snapping is implemented when a sprint label contains ISO date boundaries (for example `2026-05-01…2026-05-15`). Imported sprint names currently look like labels (`FY26 Q3`, `Sprint 1`) and do not include boundary dates, so the UI defaults snap mode on but displays that no boundary date was detected.
- `Onboarding.svelte` is implemented and wired internally to show `DateAssignment`, but story-006 scope does not include `App.svelte`, so the onboarding screen is not mounted in the app shell yet.

---

## Completion Summary
Implemented the post-import Date Assignment UI: `DateAssignment.svelte`, `DateAssignmentRow.svelte`, and `DatePicker.svelte` show missing-date candidates, sprint names, selectable rows, per-row snap toggles, inline date input/picker, batch assignment, skip/continue actions, progress/success/error states, and empty-state handling. `Onboarding.svelte` loads the latest import report and routes to DateAssignment when `missing_dates_count > 0`, then confirmation. `api.ts` now includes import-report/date-assignment types, endpoint mapping, and a `patchDate` helper.
