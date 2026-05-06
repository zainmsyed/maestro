# Story 002: Update SQLite schema for Story entity and date_source

**Status:** not-started
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

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
- [ ] Add `stories` table to `schema.sql` matching Addendum §4 (id, feature_id, title, description, status, owner, sprint, story_points, original_end_date, committed_end_date, actual_end_date, date_source, created_at, updated_at)
- [ ] Add `date_source TEXT NOT NULL DEFAULT 'imported'` to `features` table
- [ ] Add `date_source TEXT NOT NULL DEFAULT 'imported'` to `stories` table
- [ ] Add `idx_stories_feature_id` index
- [ ] Add `idx_stories_sprint` index
- [ ] Create `internal/models/story.go` with `Story` struct
- [ ] Update `internal/models/feature.go` with `DateSource` field
- [ ] Create `internal/repository/story.go` with CRUD interface and SQLite implementation
- [ ] Update `internal/repository/feature.go` to scan and write `DateSource`
- [ ] Update `internal/repository/repository.go` to include `Stories`
- [ ] Update `internal/repository/repository_test.go` with Story CRUD tests and date_source round-trip tests
- [ ] Verify `go test ./...` passes
- [ ] Verify schema with `sqlite3 maestro.db ".schema"`

---

## Issues

---

## Completion Summary
