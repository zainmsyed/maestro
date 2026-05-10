# Story 016: Update Health dashboard for story-level metrics

**Status:** not-started
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Update the Health dashboard to support story-level metric calculations, add the Orphaned Story Rate metric, exclude inherited dates from slip analysis, and add a toggle between feature-level and story-level Deadline Hit Rate.

## Verification
Open the Health dashboard with test data containing features and stories with mixed `date_source` values (imported, pm_assigned, inherited). Verify story-level Deadline Hit Rate differs from feature-level when stories have independent dates, inherited items are excluded from slip calculations, and Orphaned Story Rate appears as a new metric card.

## Scope — files this story may touch
- `frontend/src/screens/HealthView.svelte`
- `frontend/src/components/MetricCard.svelte`
- `frontend/src/components/DeadlineHitRateToggle.svelte` (new)
- `frontend/src/components/OrphanedStoryCard.svelte` (new)
- `frontend/src/lib/api.ts`

## Out of scope — do not touch
- Gantt view
- List view
- Settings screens
- Onboarding flow
- Import parser

## Dependencies
- story-004
- story-015

---

## Checklist
- [ ] Update `GET /api/metrics` backend to return both feature-level and story-level Deadline Hit Rate
- [ ] Update `GET /api/metrics` to return story-level Scope Creep Rate (stories added after sprint start)
- [ ] Update `GET /api/metrics` to return Orphaned Story Rate
- [ ] Update slip analysis backend to exclude items with `date_source = 'inherited'` from all slip calculations
- [ ] Update `HealthView.svelte` to show a toggle between Feature-level and Story-level Deadline Hit Rate
- [ ] Add Orphaned Story Rate metric card with reassignment shortcut
- [ ] Update Scope Creep Rate card to note it now includes Stories
- [ ] Update Slip Analysis section to show inherited-excluded count
- [ ] Update Sprint Load Index: unchanged (already uses story points at Story level)
- [ ] Update Velocity Trend: unchanged
- [ ] Update filter bar to include `date_source` filter option
- [ ] Write backend tests verifying inherited exclusion and story-level calculations
- [ ] Verify metrics match hand-calculated expectations for fixture data

---

## Issues

---

## Completion Summary
