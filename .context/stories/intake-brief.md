# Intake Brief

**Last updated:** 2026-05-06

## Planning brief
Use the PRD and mockup in `.context/intake/` as the source of truth for Maestro’s POC scope.

## Assumptions made
- One Azure DevOps project per `maestro.db` file for the POC.
- Hierarchy stops at Epic → Feature; User Stories are out of scope.
- Audit attribution uses a simple PM name entry on first launch.
- `delta_days` uses calendar days, not business days.
- Frontend implementation should match the reference mockup’s dark theme, typography, and layout patterns.

## Source files
- `.context/intake/prd/maestro-prd.md` (25793 bytes) — Full POC requirements, data model, API surface, metrics, milestones, and constraints.
- `.context/intake/references/maestro-mockup.html` (59259 bytes) — Reference screens for onboarding, roadmap, settings, detail panel, and visual tokens.

## Distilled notes
### Objectives
Build a self-contained roadmap and execution-health tool for PMs using Azure DevOps exports. Ship as a single Go binary with an embedded frontend and local SQLite storage.

### Success targets
- Import-to-roadmap time under 30 seconds for a 500-item export
- Binary size under 25 MB
- Startup under 2 seconds
- Smooth Gantt interaction at 200+ items
- At least 2 PMs complete a sprint cycle in the POC

### Users and journeys
Primary users are internal product managers. Core journey: import Azure DevOps CSV/JSON → review sprint setup → inspect roadmap in List and Gantt views → adjust committed dates → monitor health metrics.

### Stack and deployment
- Backend: Go
- Frontend: Vite + Svelte SPA
- Database: SQLite via `modernc.org/sqlite`
- Distribution: single binary using Go `embed.FS`

### Core data model
- `Epic`
- `Feature`
- `Sprint`
- `DateAuditLog`

Original dates are immutable. Committed dates are PM-adjustable. Orphaned features are grouped under a synthetic `Unassigned` Epic.

### Key product areas
- Import and onboarding flow
- List view with inline date editing and grouping
- Virtualized Gantt timeline with drag interactions
- Health dashboard with 6 metrics and recoveries callout
- Settings for sprint and threshold management

### Constraints and non-goals
No real-time sync, no auth or RBAC, no multi-user collaboration, no mobile-first UI, no write-back to Azure DevOps, and no non-Azure integrations in the POC.

### Notable edge cases
- Date format variance across Azure DevOps exports
- Missing sprint or date data using the fallback chain from the PRD
- Duplicate imports needing update-vs-new-project handling
- Missing sprint capacity must not break unrelated metrics
- Gantt must stay performant with large datasets
