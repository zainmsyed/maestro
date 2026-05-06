# Maestro — Plan

**Created:** 2026-05-06  
**Last updated:** 2026-05-06

---

## What we're building
Maestro is a self-contained roadmap and execution health tool for PMs using Azure DevOps. A single Go binary with an embedded Svelte frontend and SQLite database. PMs import DevOps exports, get an auto-generated Epic → Feature roadmap on a sprint timeline, and can monitor execution health via 6 metrics.

## What we're not building (v1 scope)
- Real-time sync with Azure DevOps
- Kanban or Calendar views
- User Story level hierarchy (stop at Feature)
- Multi-user collaboration, auth, or RBAC
- Mobile-optimized UI
- Write-back to Azure DevOps
- Business-day calculations (calendar days only)
- Notifications and alerts

## Features
### Feature 1: Import & Onboarding
Import Azure DevOps CSV/JSON, preview parsed data, configure sprints, confirm to generate roadmap.

### Feature 2: List View
Flat/grouped table with inline date editing, sort/filter/group, orphaned feature reassignment, CSV export.

### Feature 3: Gantt Timeline
Virtualized timeline with draggable bars, detail panel, snap-to-sprint toggle, today line.

### Feature 4: Health Dashboard
6 execution metrics with charts: deadline hit rate, sprint load, scope creep, slip analysis, velocity, orphaned features.

### Feature 5: Settings
Sprint management, metrics thresholds, project config, re-import.

## Story queue
| Story | Title | Status | Blocks |
|---|---|---|---|
| story-001 | Bootstrap Go backend with SQLite schema | not-started | — |
| story-002 | Build Azure DevOps import parser | not-started | story-001 |
| story-003 | Build REST API for data access and date mutations | not-started | story-001, story-002 |
| story-004 | Scaffold Vite + Svelte frontend | not-started | story-003 |
| story-005 | Onboarding flow | not-started | story-003, story-004 |
| story-006 | List view with inline editing and CSV export | not-started | story-003, story-004 |
| story-007 | Gantt timeline core rendering | not-started | story-003, story-004 |
| story-008 | Gantt drag interactions and detail panel | not-started | story-007 |
| story-009 | Gantt virtualization for large roadmaps | not-started | story-007, story-008 |
| story-010 | Health dashboard with 6 metrics | not-started | story-003, story-004 |
| story-011 | Settings screens | not-started | story-003, story-004, story-010 |
| story-012 | Single-binary packaging | not-started | story-004, story-011 |

## Replanning log
- 2026-05-06 — Initial plan created from PRD v0.3 and mockup.
- 2026-05-06 — Replan review against the same PRD and mockup confirmed the existing story-001 through story-012 queue still covers the current POC scope, so no new stories were added.
