# Story 002: Update SQLite schema for Story entity and date_source

**Status:** complete  
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06  
**Completed:** 2026-05-06

---

## Goal
Extend the existing SQLite schema to support the three-level hierarchy by adding a Story table, adding `date_source` to Feature and Story, and updating the repository layer and tests.

## Verification
Run `go test ./...` and see all tests pass. Inspect `maestro.db` schema and confirm: `stories` table exists, `features` has `date_source`, `stories` has `date_source`, synthetic unassigned feature support is in place.

## Scope — files this story may touch
- `internal/db/schema.sql`
- `internal/db/db.go` (migration)
- `internal/models/story.go` (new)
- `internal/models/feature.go`
- `internal/models/epic.go` (verify no date_source needed)
- `internal/repository/story.go` (new)
- `internal/repository/feature.go`
- `internal/repository/repository.go`
- `internal/repository/repository_test.go`

## Out of scope — do not touch
- HTTP handlers or routing
- Import parser logic
- Frontend code
- Existing Epic repository (no schema change)

## Dependencies
- story-001

---

## Checklist
- [x] Add `stories` table to `schema.sql` matching Addendum §4 (id, feature_id, title, description, status, owner, sprint, story_points, original_end_date, committed_end_date, actual_end_date, date_source, created_at, updated_at)
- [x] Add `date_source TEXT NOT NULL DEFAULT 'imported'` to `features` table
- [x] Add `date_source TEXT NOT NULL DEFAULT 'imported'` to `stories` table
- [x] Add `idx_stories_feature_id` index
- [x] Add `idx_stories_sprint` index
- [x] Create `internal/models/story.go` with `Story` struct
- [x] Update `internal/models/feature.go` with `DateSource` field
- [x] Create `internal/repository/story.go` with CRUD interface and SQLite implementation
- [x] Update `internal/repository/feature.go` to scan and write `DateSource`
- [x] Update `internal/repository/repository.go` to include `Stories`
- [x] Update `internal/repository/repository_test.go` with Story CRUD tests and date_source round-trip tests
- [x] Verify `go test ./...` passes
- [x] Verify schema with programmatic inspection (sqlite3 CLI unavailable)

---

## Issues

---

## Completion Summary
Extended the SQLite schema and repository layer to support the three-level hierarchy (Epic → Feature → Story). Added a `stories` table with a foreign key to `features`, along with `date_source` columns on both `features` and `stories` (defaulting to `'imported'`). Created the `Story` model and corresponding `StoryRepository` with CRUD operations, updated `FeatureRepository` to read and write `DateSource`, and wired `Stories` into the central `Repositories` struct. All repository methods enforce the `'imported'` default at the code layer when an empty `DateSource` is passed. Added comprehensive tests for Story CRUD, Feature date_source round-trip, and Story date_source round-trip. All tests pass.
