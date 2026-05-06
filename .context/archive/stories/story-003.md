# Story 003: Build REST API for data access and date mutations

**Status:** retired  
**Type:** —  
**Created:** 2026-05-06  
**Last accessed:** 2026-05-06  
**Completed:** —

**Superseded by:** story-004 (Update REST API for Story entity and date_source)

**Reason:** Original two-level REST API superseded by the three-level API with Story entity and date_source support defined in the PRD addendum.

---

## Goal
Implement all REST endpoints defined in the PRD §7.1, including CRUD for epics/features/sprints, date mutations with audit logging, epic reassignment, and health metrics queries.

## Verification
Run the API test suite and confirm every endpoint returns correct status codes and payloads. Use `curl` or `httpie` to PATCH a feature date and verify a corresponding `DateAuditLog` row is created with accurate `delta_days`.

## Scope — files this story may touch
- `internal/api/server.go`
- `internal/api/routes.go`
- `internal/api/epics.go`
- `internal/api/features.go`
- `internal/api/sprints.go`
- `internal/api/metrics.go`
- `internal/api/audit.go`
- `internal/api/import.go`
- `internal/repository/*.go`

## Out of scope — do not touch
- Frontend components
- File upload UI (endpoint accepts multipart; UI comes later)
- Authentication or authorization
- Go embed.FS static serving

## Dependencies
- story-001
- story-002

---

## Checklist
- [ ] Set up HTTP router (chi or stdlib mux) with JSON request/response helpers
- [ ] `GET /api/epics` — list all epics with nested features
- [ ] `GET /api/epics/:id` — single epic detail
- [ ] `GET /api/features` — list all features
- [ ] `GET /api/features/:id` — single feature detail
- [ ] `GET /api/sprints` — list all sprints
- [ ] `POST /api/sprints` — create sprint
- [ ] `PATCH /api/sprints/:id` — update sprint dates/capacity/name
- [ ] `PATCH /api/epics/:id/date` — update `committed_end_date`, write audit log
- [ ] `PATCH /api/features/:id/date` — update `committed_end_date`, write audit log
- [ ] `PATCH /api/features/:id/epic` — reassign feature to different epic
- [ ] `GET /api/metrics` — return all health metrics
- [ ] `GET /api/metrics/slip/:id` — slip event history for a feature
- [ ] `GET /api/metrics/recoveries` — features with negative net slip days
- [ ] `GET /api/audit` — date change audit log
- [ ] Implement `delta_days` as calendar days on every date change
- [ ] Write API integration tests covering happy paths and errors
- [ ] Verify audit log entries are created with old date, new date, and delta

---

## Issues

---

## Completion Summary
