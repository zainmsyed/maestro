# Story 004: Update REST API for Story entity and date_source

**Status:** complete  
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-07  
**Completed:** 2026-05-07

---

## Goal
Extend the REST API to support Story CRUD, update Feature endpoints with `date_source` and epic reassignment, update Epic endpoints, add metrics endpoints for story-level calculations, and update date mutation logic to handle PM-assigned original dates and inherited dates.

## Verification
Run the API test suite and confirm every endpoint returns correct status codes and payloads. Use `curl` to POST a Story, PATCH a Feature date (verifying `date_source` and audit log), and GET metrics to confirm story-level calculations exclude inherited dates.

## Scope — files this story may touch
- `internal/api/server.go`
- `internal/api/routes.go`
- `internal/api/stories.go` (new)
- `internal/api/epics.go`
- `internal/api/features.go`
- `internal/api/metrics.go`
- `internal/api/audit.go`
- `internal/repository/*.go` (read paths)
- `internal/models/*.go`

## Out of scope — do not touch
- Import parser (already built in story-003)
- Frontend components
- File upload UI
- Authentication or authorization

## Dependencies
- story-002
- story-003

---

## Checklist
- [x] Add Story endpoints: `GET /api/stories`, `GET /api/stories/:id`, `POST /api/stories` (for synthetic/manual), `PATCH /api/stories/:id/date`, `PATCH /api/stories/:id/feature` (reassign)
- [x] Update Feature `GET /api/features/:id` to include `date_source` and nested stories
- [x] Update Feature `PATCH /api/features/:id/date` to write `date_source` and audit log
- [x] Update Epic `GET /api/epics/:id` to include nested features with their nested stories
- [x] Update Epic `GET /api/epics` to include full three-level nested structure
- [x] Update `PATCH /api/features/:id/date`: if `original_end_date` is NULL (undated item), first assignment sets both `original_end_date` and `committed_end_date` to the new date, `date_source = 'pm_assigned'`, `delta_days = 0`
- [x] Update `PATCH /api/stories/:id/date`: same first-assignment logic as Feature
- [x] Update `GET /api/metrics` to include story-level Deadline Hit Rate and story-level Scope Creep Rate
- [x] Update `GET /api/metrics/slip/:id` to exclude items with `date_source = 'inherited'` from slip calculations
- [x] Add `GET /api/metrics/orphaned-stories` endpoint
- [x] Update `GET /api/audit` to include `date_source` in relevant responses
- [x] Add `GET /api/import/report` endpoint returning the last import report
- [x] Write API integration tests covering Story CRUD, first-assignment logic, inherited exclusion, and three-level nesting
- [x] Verify three-level nesting in API responses

---

## Issues
- **Out-of-scope modification:** `internal/importer/types.go` was modified to alias `ImportReport` from `internal/models` in order to break an import cycle between `internal/repository/import_report.go` and `internal/importer`. This was necessary because the new `ImportReportRepository` lives in `repository` and needs the `ImportReport` type, but `ImportReport` was originally defined in `importer`. The type was extracted to `internal/models/import_report.go` and aliased in `importer/types.go` to preserve backward compatibility.

---

## Completion Summary

Created a complete REST API layer (`internal/api`) and extended repositories to support Story CRUD, Feature/Epic nested responses, date mutation with audit logging, metrics calculations, and import report retrieval.

**New files:**
- `internal/api/server.go` — HTTP server setup with JSON helpers
- `internal/api/routes.go` — route registration and path dispatch
- `internal/api/stories.go` — Story endpoints (list, get, create, patch date, patch feature)
- `internal/api/features.go` — Feature endpoints (list, get, create, patch date, patch epic)
- `internal/api/epics.go` — Epic endpoints (list, get, create) with three-level nesting
- `internal/api/metrics.go` — Metrics endpoints (overview, slip by id, orphaned stories)
- `internal/api/audit.go` — Audit list with enriched `date_source`
- `internal/api/import.go` — Last import report retrieval
- `internal/api/api_test.go` — Comprehensive integration tests
- `internal/repository/import_report.go` — Import report persistence
- `internal/repository/metrics.go` — Metrics queries (deadline hit, scope creep, slip, orphans)
- `internal/models/import_report.go` — `ImportReport` model (breaks import cycle)

**Modified files:**
- `cmd/maestro/main.go` — starts the HTTP server
- `internal/db/schema.sql` — added `import_reports` table
- `internal/repository/story.go` — added `ListByFeatureID`, `UpdateDate`, `UpdateFeatureID`
- `internal/repository/feature.go` — added `ListByEpicID`, `UpdateDate`, `UpdateEpicID`
- `internal/repository/repository.go` — added `ImportReports` and `Metrics` repositories
- `internal/importer/types.go` — reuses `models.ImportReport` to break import cycle

**Key behaviors implemented:**
- First date assignment on undated items sets both `original_end_date` and `committed_end_date`, marks `date_source = 'pm_assigned'`, and logs audit with `delta_days = 0`
- Subsequent assignments compute `delta_days` from previous committed date and keep `date_source = 'pm_assigned'`
- Slip calculations exclude items where `date_source = 'inherited'`
- Audit responses include the current `date_source` of the affected entity
- Three-level nesting (epic → features → stories) verified in `GET /api/epics` and `GET /api/epics/:id`
- All endpoints covered by integration tests passing with in-memory SQLite
