# Maestro — Plan

**Created:** 2026-05-06  
**Last updated:** 2026-05-09

---

## What we're building
Maestro is a self-contained roadmap and execution health tool for PMs using Azure DevOps. A single Go binary with an embedded Svelte frontend and SQLite database. PMs run a prescribed "Tree of Work Items" query in DevOps, export to CSV, import into Maestro, and get an auto-generated Epic → Feature → Story roadmap on a sprint timeline with execution health metrics.

## What we're not building (v1 scope)
- Real-time sync with Azure DevOps
- Kanban or Calendar views
- Multi-user collaboration, auth, or RBAC
- Mobile-optimized UI
- Write-back to Azure DevOps
- Business-day calculations (calendar days only)
- Notifications and alerts
- JSON import (CSV via prescribed Tree query only)

## Features
### Feature 1: Import & Onboarding
Prescribed Azure DevOps Tree query CSV import, header normalization, three-level hierarchy parsing, preview with import report, post-import PM date assignment, sprint configuration, confirm to generate roadmap.

### Feature 2: List View
Three-level grouped table with inline date editing, date source badges, click-to-assign for undated items, sort/filter/group, orphaned feature/story reassignment, CSV export.

### Feature 3: Gantt Timeline
SVAR Svelte Gantt (MIT edition) with Epic → Feature hierarchy, drag-to-adjust dates, snap-to-sprint toggle, today line overlay, sprint boundary overlays, status color coding, synthetic "Unassigned" Epic muted styling, and detail panel on task click. Custom virtualized renderer replaced by SVAR component.

### Feature 4: Health Dashboard
7 execution metrics with charts: deadline hit rate (feature + story level), sprint load, scope creep (story-level), slip analysis with inherited exclusion, velocity, orphaned features, orphaned stories.

### Feature 5: Settings
Sprint management, metrics thresholds, project config, re-import.

## Story queue
| Story | Title | Status | Blocks |
|---|---|---|---|
| story-001 | Bootstrap Go backend with SQLite schema | **archived** | — |
| story-002 | Update SQLite schema for Story entity and date_source | **archived** | story-001 |
| story-003 | Build prescribed query CSV parser with three-level hierarchy | **archived** | story-002 |
| story-004 | Update REST API for Story entity and date_source | **complete** | story-002, story-003 |
| story-005 | Scaffold Vite + Svelte frontend | **complete** | story-004 |
| story-006 | Build post-import date assignment screen | **complete** | story-003, story-005 |
| story-007 | Onboarding flow | **complete** | story-004, story-005 |
| story-008 | List view with inline editing and CSV export | **complete** | story-004, story-005 |
| story-009 | Update List view for three-level hierarchy and date assignment | **complete** | story-004, story-008 |
| story-010 | SVAR M0 spike — install, bundle baseline, API verification | not-started | story-005 |
| story-011 | SVAR data bridge + basic Gantt screen | not-started | story-004, story-005, story-010 |
| story-012 | SVAR styling + today line + sprint boundary overlays | not-started | story-011 |
| story-013 | SVAR drag handler + snap-to-sprint | not-started | story-011, story-012 |
| story-014 | SVAR detail panel | not-started | story-011, story-012 |
| story-015 | Health dashboard with 6 metrics | not-started | story-004, story-005 |
| story-016 | Update Health dashboard for story-level metrics | not-started | story-004, story-015 |
| story-017 | Settings screens | not-started | story-004, story-005, story-015 |
| story-018 | Single-binary packaging | not-started | story-005, story-017 |

## Replanning log
- 2026-05-06 — Initial plan created from PRD v0.3 and mockup.
- 2026-05-06 — Replan review confirmed existing story-001 through story-012 queue covers base POC scope; no new stories added.
- 2026-05-06 — Replan for PRD addendum (prescribed query, three-level hierarchy, Story entity, date_source, PM date insertion). Added story-013 through story-019 for delta work.
- 2026-05-06 — Archived story-002 and story-003 (two-level parser and API superseded by addendum). Renumbered all active stories sequentially: story-002 through story-017.
- 2026-05-09 — Replan for SVAR Gantt addendum. Renumbered active stories: story-010 through story-014 now cover SVAR integration (M0 spike, data bridge, styling/overlays, drag/snap, detail panel). Story-015 through story-018 shifted from previous 014–017 (Health dashboard, Settings, Binary packaging). Custom virtualized Gantt renderer superseded by SVAR MIT edition. Epic → Feature hierarchy in Gantt; Stories remain in List view and Health dashboard only.
