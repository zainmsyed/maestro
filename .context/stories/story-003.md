# Story 003: Build prescribed query CSV parser with three-level hierarchy

**Status:** not-started
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Build a CSV parser that reads Azure DevOps "Tree of Work Items" query exports, validates the prescribed column set, normalizes headers, parses Epic → Feature → Story hierarchy via Title 1/2/3, handles Story type variants, detects date formats, generates synthetic IDs for empty Story IDs, and produces an import report with a date-assignment candidate list.

## Verification
Import the test fixture `testdata/tree-export.csv` (containing Epic, Feature, and Story rows with empty IDs and missing Target Dates) and query the database to verify: correct epic count, correct feature count, correct story count, synthetic IDs generated for stories with empty IDs, items with missing Target Date flagged in the import report, orphaned features/stories assigned to synthetic unassigned epic/feature.

## Scope — files this story may touch
- `internal/importer/*.go`
- `internal/importer/csv.go`
- `internal/importer/validation.go`
- `internal/importer/hierarchy.go`
- `internal/importer/dates.go`
- `internal/importer/types.go`
- `testdata/*.csv`
- `internal/repository/*.go` (write paths)

## Out of scope — do not touch
- JSON import (removed from POC scope)
- HTTP upload handlers (parser accepts `io.Reader`)
- Preview UI or onboarding flow
- Re-import logic (duplicate handling)
- Frontend code

## Dependencies
- story-002

---

## Checklist
- [ ] Create CSV header normalizer mapping canonical names to internal fields (Title 1→epic_title, Title 2→feature_title, Title 3→story_title, etc.)
- [ ] Validate required columns present: ID, Work Item Type, Title 1, Title 2, Title 3, State, Assigned To, Iteration Path, Story Points, Target Date, Area Path
- [ ] Reject non-Tree query exports (detect by requiring Title 1/2/3)
- [ ] Parse Work Item Type: Epic→epic, Feature→feature, User Story/Product Backlog Item/Requirement→story
- [ ] Parse hierarchy: Title 1 populated = epic row, Title 2 populated = feature row, Title 3 populated = story row; track current epic/feature pointers
- [ ] Handle empty Story IDs by generating synthetic UUIDs (`story-{hash}` or `story-{auto}`)
- [ ] Parse Target Date with format detection (ISO 8601, MM/DD/YYYY, DD/MM/YYYY, ISO slashes, long form, abbreviated, ISO with timezone)
- [ ] Flag ambiguous dates (day ≤ 12) for confirmation; default to US format per Addendum §6.5
- [ ] Empty Target Date → `original_end_date = NULL`, add to date-assignment candidate list
- [ ] Parse Iteration Path: extract last segment as sprint name
- [ ] Parse Assigned To: extract display name and optional email
- [ ] Parse Story Points: empty→nil, numeric→int, non-numeric→log warning, nil
- [ ] Generate synthetic unassigned epic if orphaned features exist
- [ ] Generate synthetic unassigned feature under unassigned epic if orphaned stories exist
- [ ] Lock `original_end_date` on import; set `committed_end_date` to same value; set `date_source = 'imported'`
- [ ] Build import report: counts per type, sprints detected, missing dates count, missing sprint count, orphaned features, orphaned stories, skipped rows, detected date format
- [ ] Create test fixtures: `testdata/tree-export-v1.csv` with varied date formats, empty IDs, orphaned items
- [ ] Write unit tests for header normalization, hierarchy parsing, date detection, type normalization
- [ ] Run import against test fixture and verify record counts, field values, and synthetic ID generation

---

## Issues

---

## Completion Summary
