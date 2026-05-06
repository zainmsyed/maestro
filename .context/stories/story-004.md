# Story 004: Update REST API for Story entity and date_source

**Status:** not-started
**Type:** —
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

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
- [ ] Add Story endpoints: `GET /api/stories`, `GET /api/stories/:id`, `POST /api/stories` (for synthetic/manual), `PATCH /api/stories/:id/date`, `PATCH /api/stories/:id/feature` (reassign)
- [ ] Update Feature `GET /api/features/:id` to include `date_source` and nested stories
- [ ] Update Feature `PATCH /api/features/:id/date` to write `date_source` and audit log
- [ ] Update Epic `GET /api/epics/:id` to include nested features with their nested stories
- [ ] Update Epic `GET /api/epics` to include full three-level nested structure
- [ ] Update `PATCH /api/features/:id/date`: if `original_end_date` is NULL (undated item), first assignment sets both `original_end_date` and `committed_end_date` to the new date, `date_source = 'pm_assigned'`, `delta_days = 0`
- [ ] Update `PATCH /api/stories/:id/date`: same first-assignment logic as Feature
- [ ] Update `GET /api/metrics` to include story-level Deadline Hit Rate and story-level Scope Creep Rate
- [ ] Update `GET /api/metrics/slip/:id` to exclude items with `date_source = 'inherited'` from slip calculations
- [ ] Add `GET /api/metrics/orphaned-stories` endpoint
- [ ] Update `GET /api/audit` to include `date_source` in relevant responses
- [ ] Add `GET /api/import/report` endpoint returning the last import report
- [ ] Write API integration tests covering Story CRUD, first-assignment logic, inherited exclusion, and three-level nesting
- [ ] Verify three-level nesting in API responses

---

## Issues

---

## Completion Summary
