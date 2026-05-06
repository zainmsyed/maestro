# Story 006: Build post-import date assignment screen

**Status:** not-started
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

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
- [ ] Create `DateAssignment.svelte` screen showing items missing Target Date with count header
- [ ] Each row shows: checkbox (select for batch), item title, type (Epic/Feature/Story), sprint name, date picker
- [ ] Date picker snaps to sprint boundaries by default (toggle per row)
- [ ] "Set all selected" button assigns dates to all checked items
- [ ] "Skip — I'll do it later" button continues to confirmation without assigning
- [ ] "Set dates now" button assigns dates and continues
- [ ] Call PATCH endpoints for each assigned item; first assignment sets both `original_end_date` and `committed_end_date`
- [ ] Show progress/success state after batch assignment
- [ ] Wire `Onboarding.svelte` to show `DateAssignment` between sprint review and confirmation when `missing_dates > 0`
- [ ] Handle empty state (all items have dates → skip screen automatically)
- [ ] Handle errors gracefully (network failure, invalid date)

---

## Issues

---

## Completion Summary
