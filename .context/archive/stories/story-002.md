# Story 002: Build Azure DevOps import parser

**Status:** retired  
**Type:** —  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

**Superseded by:** story-003 (Build prescribed query CSV parser with three-level hierarchy)

**Reason:** Original two-level CSV/JSON parser superseded by the prescribed Tree query CSV parser defined in the PRD addendum.

---

## Goal
Build a parser that reads Azure DevOps CSV and JSON exports, maps fields to the Maestro data model, handles date format variance, infers Epic → Feature hierarchy, and assigns orphaned features to a synthetic "Unassigned" Epic.

## Verification
Import the synthetic test fixtures (`testdata/export-v1.csv`, `testdata/export-v2.json`) and query the database to verify: 4 epics, 10 features, 1 synthetic epic, 9 sprints, all dates parsed correctly regardless of input format.

## Scope — files this story may touch
- `internal/importer/*.go`
- `internal/importer/csv.go`
- `internal/importer/json.go`
- `internal/importer/dates.go`
- `internal/importer/hierarchy.go`
- `testdata/*.csv`, `testdata/*.json`
- `internal/repository/epic.go`
- `internal/repository/feature.go`
- `internal/repository/sprint.go`

## Out of scope — do not touch
- HTTP upload handlers (parser accepts `io.Reader`)
- Preview UI or onboarding flow
- Re-import logic (duplicate handling)
- Frontend code

## Dependencies
- story-001

---

## Checklist
- [ ] Create CSV reader that handles standard Azure DevOps export columns
- [ ] Create JSON reader for Azure DevOps JSON export shape
- [ ] Implement date parser supporting `MM/DD/YYYY`, `YYYY-MM-DD`, `DD/MM/YYYY`, and ISO 8601 with timezone offsets
- [ ] Map Work Item Type to Epic or Feature; map Parent field to `epic_id`
- [ ] Build hierarchy inference: Features without parent → synthetic Epic (`is_synthetic = true`, title "Unassigned")
- [ ] Extract sprint data and create Sprint records with `source = imported`
- [ ] Lock `original_end_date` on import and set `committed_end_date` to the same value
- [ ] Create test fixtures with varied date formats and orphaned features
- [ ] Write unit tests for each parser format and edge case
- [ ] Run import against test fixtures and verify record counts and field values

---

## Issues

---

## Completion Summary
