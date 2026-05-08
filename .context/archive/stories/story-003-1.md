# Story 003: Build prescribed query CSV parser with three-level hierarchy

**Status:** complete  
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06  
**Completed:** 2026-05-06

---

## Goal
Build a CSV parser that reads Azure DevOps "Tree of Work Items" query exports, validates the prescribed column set, normalizes headers, parses Epic → Feature → Story hierarchy via Parent ID (addendum v2.0), handles Story type variants, detects date formats, generates synthetic IDs for empty Story IDs, and produces an import report with a date-assignment candidate list.

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
- [x] Require `Parent`, `ID`, and `Work Item Type` columns; reject import if any missing
- [x] Detect all `Title N` columns dynamically and use last non-empty value as display title
- [x] Normalize headers with aliases per addendum v2.0 (Parent ID, WorkItemType, Sprint, etc.)
- [x] Parse Work Item Type: Epic→epic, Feature→feature, User Story/Product Backlog Item/Requirement→story
- [x] Normalize float IDs by stripping `.0` from `ID` and `Parent`
- [x] Two-pass hierarchy linking via `Parent` ID (row order irrelevant)
- [x] Handle empty Story IDs by generating synthetic IDs (`story-auto-*`)
- [x] Parse Target Date with format detection (ISO 8601, MM/DD/YYYY, DD/MM/YYYY, ISO slashes, US datetime, long form, abbreviated, ISO with timezone)
- [x] Flag ambiguous dates (day ≤ 12) for confirmation; default to US format per Addendum §6.5
- [x] Empty Target Date → `original_end_date = NULL`, add to date-assignment candidate list
- [x] Parse Iteration Path with sprint pattern detection and archive/backlog/queue unscheduled handling
- [x] Parse Assigned To: extract display name and optional email
- [x] Parse Story Points: empty→nil, numeric→int, float→rounded int, non-numeric→log warning, nil
- [x] Generate synthetic unassigned epic if orphaned features exist
- [x] Generate synthetic unassigned feature under unassigned epic if orphaned stories exist
- [x] Lock `original_end_date` on import; set `committed_end_date` to same value; set `date_source = 'imported'`
- [x] Build import report: counts per type, sprints detected, missing dates count, missing sprint count, orphaned features, orphaned stories, skipped rows, detected date format
- [x] Create test fixtures: `testdata/tree-export.csv` and `testdata/tree-export-v1.csv` with Parent IDs, varied date formats, empty IDs, orphaned items, archive/backlog paths
- [x] Write unit tests for header normalization, ID normalization, dynamic title extraction, date detection, type normalization, iteration path parsing, story points parsing
- [x] Run import against test fixture and verify record counts, field values, and synthetic ID generation
- [x] Verify against real reference CSV: `.context/intake/references/Example_FinancialDashboard_Backlog.csv` → skipped rows dropped from 807 to 3

---

## Issues

### /fix — "Importer still uses title-column hierarchy instead of Parent ID hierarchy from addendum v2.0; real reference CSV skips 807 rows."
- **Reported:** 2026-05-06  
- **Status:** resolved  
- **Agent note:** Reference CSV `.context/intake/references/Example_FinancialDashboard_Backlog.csv` imports with `epics=72`, `features=324`, `stories=387`, `skipped=807`; many skipped rows warn `missing hierarchy title columns`. PRD addendum v2.0 supersedes title-column hierarchy and requires `Parent` ID reconstruction plus dynamic title extraction from all `Title N` columns.  
- **Solution:** Rebuilt importer to follow `.context/intake/prd/maestro-prd-addendum-devops-query.md` §5.1–§5.8. After rebuild, same reference CSV imports with `epics=73`, `features=331`, `stories=1181`, `skipped=3` (only unsupported work item types: Spike, Task, Report). Orphaned features/stories dropped to 0. Date format correctly detected as `US datetime`.

---

## Completion Summary
Added a new `internal/importer` package that parses Azure DevOps Tree of Work Items CSV exports from an `io.Reader` per PRD addendum v2.0. The importer validates required columns (`Parent`, `ID`, `Work Item Type`), normalizes headers with aliases, strips float decimals from IDs, dynamically scans all `Title N` columns for display titles, and reconstructs hierarchy via two-pass `Parent` ID linking so row order and title depth are irrelevant. It detects multiple target date formats including US datetime strings, flags ambiguous slash dates, parses iteration paths with sprint pattern recognition and archive/backlog/queue unscheduled handling, handles float story points by rounding, generates `story-auto-*` IDs for empty story IDs, and records missing-date candidates plus warnings in an import report. The importer creates synthetic `epic-unassigned` and `feature-unassigned` records for orphaned features/stories, locks imported dates by setting `original_end_date` and `committed_end_date` together, and persists imported features/stories with `date_source = 'imported'`. Added fixtures in `testdata/tree-export.csv` and `testdata/tree-export-v1.csv`, plus importer unit/integration tests covering header normalization, ID normalization, dynamic title extraction, work item type normalization, iteration path parsing, date parsing, story points parsing, and a full import that verifies database counts, synthetic ID generation, orphan assignment, missing-date reporting, and mixed date-format detection. Verified against the real reference CSV: skipped rows dropped from 807 to 3.
