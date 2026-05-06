# Story 001: Bootstrap Go backend with SQLite schema

**Status:** complete  
**Type:** —  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** 2026-05-06

---

## Goal
Set up the Go project with all database tables and a working SQLite connection using `modernc.org/sqlite`. This is the foundation every other story builds on.

## Verification
Run `go test ./...` and see all repository tests pass. Open `maestro.db` with `sqlite3` and confirm all four tables exist with correct columns and indexes.

## Scope — files this story may touch
- `go.mod`, `go.sum`
- `cmd/maestro/main.go`
- `internal/db/db.go`
- `internal/db/schema.sql`
- `internal/models/epic.go`
- `internal/models/feature.go`
- `internal/models/sprint.go`
- `internal/models/audit.go`
- `internal/repository/*.go`
- `internal/config/*.go`

## Out of scope — do not touch
- HTTP handlers or routing
- Import parsing logic
- Frontend code
- Build scripts for cross-compilation

## Dependencies
- None

---

## Checklist
- [x] Initialize Go module (`go mod init`) and add `modernc.org/sqlite` dependency
- [x] Create `cmd/maestro/main.go` with basic CLI flag parsing (port, db path)
- [x] Write `internal/db/schema.sql` with CREATE TABLE for Epic, Feature, Sprint, DateAuditLog
- [x] Implement `internal/db/db.go` to open SQLite, run migrations, and return `*sql.DB`
- [x] Define model structs in `internal/models/` matching the PRD §4.1 schema
- [x] Create repository interfaces and basic query helpers in `internal/repository/`
- [x] Add unit tests for schema creation and basic CRUD on each table
- [x] Verify `go test ./...` passes and `maestro.db` is created on first run

---

## Issues

---

## Completion Summary
Implemented the Go bootstrap for Maestro with a SQLite-backed schema using `modernc.org/sqlite`, CLI config parsing, model structs, and repository helpers for Epic, Feature, Sprint, and DateAuditLog. Added migration-backed database initialization, schema/index coverage tests, and repository CRUD tests. Verified the foundation by running `go test ./...` successfully and confirming `maestro.db` is created on first run via `go run ./cmd/maestro -port 9090 -db maestro.db`.
