# Intake Brief

**Last updated:** 2026-05-09

## Planning brief
Replan for SVAR Svelte Gantt integration. Retire custom-built Gantt stories (010–013) and replace with SVAR MIT-edition integration stories (018–022). Gantt renders Epic → Feature hierarchy only; Stories remain in List view and Health dashboard.

## Source files
- .context/intake/prd/maestro-prd.md
- .context/intake/prd/maestro-prd-addendum-devops-query.md
- .context/intake/prd/maestro-prd-addendum-svar-gantt.md
- .context/intake/references/Example_FinancialDashboard_Backlog.csv
- .context/intake/references/maestro-mockup.html

## Distilled notes
### PRD v0.3 + DevOps Query Addendum v2.0
Core product unchanged: single-binary Go + Svelte SPA with SQLite. Prescribed Azure DevOps Tree query CSV import. Three-level hierarchy (Epic → Feature → Story) in data model and List view. Sprint-based timeline. Health dashboard with 6 metrics. Settings screens.

### SVAR Gantt Addendum v1.0
Custom virtualized Gantt renderer (PRD §5.2 M3) replaced by `@svar-ui/svelte-gantt` MIT edition (v2.x for Svelte 5).

**Native SVAR capabilities used:**
- Hierarchical rows via `parent` field (Epic → Feature)
- Drag-to-resize and drag-to-move on task bars
- `onTaskClick` event for detail panel trigger
- Expand/collapse summary rows (`open` flag)
- Zoom via configurable `scales` array
- Custom task styling via `css` callback
- Virtualized rendering claimed by library

**Custom work required:**
- Data bridge: flat SVAR task array from Maestro Epic/Feature entities
- Today line overlay (custom div, PRO-only in SVAR)
- Sprint boundary overlay lines (custom CSS on day-based timeline)
- Status color coding via `css` callback
- Synthetic "Unassigned" Epic muted styling
- Snap-to-sprint handler in `onAfterTaskDrag` (no native snap in MIT)
- Detail panel slide-out on task click

**Items with no valid start date:** omitted from Gantt array; flagged in List view.

**Bundle size:** M0 baseline measurement required before integration. Target < 25 MB binary. M3 ceiling tightened to 20 MB.

**Licensing:** MIT edition only. No PRO features required.

## Assumptions
- Gantt shows Epic → Feature only. Stories exist in data model but are not rendered in Gantt (consistent with PRD §5.2).
- Items with NULL start date (no sprint, no fallback) are excluded from Gantt array.
- Today line and sprint boundaries are custom CSS overlays, not SVAR PRO markers.
- Existing backend API (story-004) already supports PATCH endpoints for date updates.
- `original_end_date` is never mutated by drag; only `committed_end_date` changes.

## Planning rules
- Treat listed source files as user-authored planning inputs unless they are explicitly marked as generated artifacts.
- Vazir-generated files in .context/stories/ are replan context, not primary intake.
- Read all text-based planning sources before asking questions.
- Ask only implementation-blocking delta questions after reviewing this brief and any raw files you actually need.
- State safe default assumptions briefly so the user can correct them.
- Surface contradictions instead of resolving them silently.
