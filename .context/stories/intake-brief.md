# Intake Brief

**Last updated:** 2026-05-06

## Planning brief
Maestro is a self-contained roadmap and execution health tool for PMs using Azure DevOps. The PRD addendum introduces a prescribed "Tree of Work Items" query with a three-level hierarchy (Epic → Feature → Story), CSV-only import, header normalization, post-import PM date assignment, and `date_source` tracking.

## Assumptions made
- One Azure DevOps project per `maestro.db` file for the POC.
- Hierarchy is three levels: Epic → Feature → Story. Story type names vary by process template (User Story / Product Backlog Item / Requirement).
- Story rows in the prescribed query export may have empty `ID` fields; Maestro generates synthetic IDs for these.
- JSON import is removed from POC scope; only CSV via the prescribed Tree query is supported.
- `date_source` field tracks origin: `imported` (from DevOps Target Date), `pm_assigned` (set by PM in Maestro), `inherited` (from sprint end_date).
- Items with `date_source = inherited` are excluded from Deadline Hit Rate and Slip Analysis.
- The first PM-assigned date to an undated item becomes its locked `original_end_date`.
- PM date assignment screen appears post-import before the Gantt; it can be skipped.
- Synthetic unassigned epic and synthetic unassigned feature are auto-created for orphaned items.
- Frontend design follows the reference mockup where possible; three-level nesting is an extension.

## Source files
- `.context/intake/prd/maestro-prd.md` (25793 bytes) — Base POC requirements, data model, API spec, metrics.
- `.context/intake/prd/maestro-prd-addendum-devops-query.md` (22996 bytes) — Prescribed query, three-level hierarchy, Story entity, parser spec, PM date insertion, view updates.
- `.context/intake/references/maestro-mockup.html` (59259 bytes) — Reference screens with design tokens (two-level hierarchy; three-level is an extension).

## Distilled notes

### Objectives
Build a self-contained roadmap and execution-health tool for PMs using Azure DevOps. POC targets: import-to-roadmap < 30s for 500 items, binary < 25 MB, startup < 2s.

### Users
Internal product managers using Azure DevOps. Single-user local tool; no auth or multi-user in POC.

### Stack
Go backend, SQLite via `modernc.org/sqlite`, Vite + Svelte SPA frontend, Go `embed.FS` for single-binary distribution.

### Data model (updated)
Epic, Feature, Sprint, DateAuditLog, **Story**. `original_end_date` is immutable once set (on import or first PM assignment). `committed_end_date` is PM-adjustable. `date_source` tracks origin. Synthetic unassigned epic and feature catch orphans.

### Key product areas
- **Import & Onboarding:** Prescribed Tree query CSV only, header normalization, validation, preview, post-import date assignment, sprint config.
- **List View:** Three-level grouping, inline date editing, date source badges, click-to-assign, CSV export.
- **Gantt Timeline:** Virtualized three-level nesting, draggable bars, stub bars for undated items, detail panel, snap-to-sprint toggle.
- **Health Dashboard:** 6 base metrics + story-level variants, orphaned story rate, inherited date exclusion.
- **Settings:** Sprint management, metrics thresholds, project config, re-import.

### Constraints and non-goals
No real-time sync, no auth/RBAC, no multi-user collaboration, no mobile-first UI, no write-back to Azure DevOps, no non-Azure integrations, no JSON import, no business-day math.

### Notable edge cases
- Empty Story IDs in Tree query exports → synthetic IDs.
- Empty Target Date → post-import date assignment list.
- Date format variance across regional DevOps settings.
- Ambiguous US (MM/DD/YYYY) vs EU (DD/MM/YYYY) dates.
- Missing sprint or date data using fallback chain.
- Orphaned Stories without parent Features → synthetic unassigned feature under unassigned epic.
