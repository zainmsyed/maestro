# Maestro — Product Requirements Document

**Version:** 0.3 (POC)
**Status:** Draft
**Last Updated:** 2026-05-05
**Author:** Zain

---

## Changelog

| Version | Changes |
|---|---|
| 0.3 | Clarified orphaned feature placement in Gantt (synthetic "Unassigned" Epic row). Defined explicit sprint start date fallback chain. Specified delta_days as calendar days (business days deferred to v1.1). Scoped M0 parser spike to internal teams only. Added "Recoveries" callout to health dashboard. Added binary size M0 checkpoint for modernc.org/sqlite baseline. |
| 0.2 | Updated frontend stack (SvelteKit → Vite + Svelte SPA). Redesigned slip rate metric to track slip events and net slip days. Added orphaned features health metric. Changed Gantt snap-to-sprint from constraint to toggle. Added Gantt virtualization requirement to M3. Tightened sprint capacity UX. Added M0 parser spike milestone. Added M6 strategic research questions. Added execution risk register. |
| 0.1 | Initial draft |

---

## 1. Overview

### 1.1 Product Summary

Maestro is a self-contained roadmap and execution health tool for product managers. It ingests Azure DevOps exports, automatically constructs an Epic → Feature hierarchy mapped to sprint timelines, and gives PMs a clear view of what is planned, what is in flight, and how well the team is executing against commitments.

The long-term vision is a unified PM + engineering ticketing system. This document scopes the POC only.

### 1.2 Problem Statement

Product managers spend significant time manually transcribing work from engineering tools (Azure DevOps, Jira) into roadmap tools (Productboard, spreadsheets, slide decks). These artifacts immediately fall out of sync, contain no execution health data, and give leadership a false sense of control. There is no single place where a PM can answer: *Are we hitting our deadlines? Are we bringing too much work into a sprint? Is this quarter's plan realistic?*

### 1.3 Solution

Maestro eliminates the transcription step entirely. A PM uploads an Azure DevOps export and Maestro builds the roadmap automatically — epics, features, sprint assignments, and dates. The PM's job shifts from data entry to date adjustment and health monitoring.

### 1.4 Target Users (POC)

- Product Managers at companies using Azure DevOps
- Internal use; not a public SaaS product at this stage

---

## 2. Goals & Non-Goals

### 2.1 POC Goals

- Ingest Azure DevOps CSV/JSON exports and parse Epic → Feature hierarchy with sprint data
- Render two roadmap views: **Timeline (Gantt)** and **List**
- Allow PMs to adjust dates via drag (Gantt) or inline date entry (both views)
- Track and display execution health metrics (deadline hit rate, sprint load, scope creep, slip analysis, velocity, orphaned features, schedule recoveries)
- Ship as a **single Go binary** with embedded frontend — no install steps for end users

### 2.2 Future Scope (v1.1+)

- Kanban board view
- Calendar view
- Full ticketing system for PMs and engineers
- Two-way sync with Azure DevOps (write back)
- Multi-user collaboration
- Notifications and alerts
- Business-day calculation for Net Slip Days

### 2.3 Non-Goals (POC)

- Real-time sync with Azure DevOps
- Role-based access control
- Authentication / multi-user
- Mobile-optimized UI
- Integration with tools other than Azure DevOps

---

## 3. Architecture

### 3.1 Stack

| Layer | Technology | Rationale |
|---|---|---|
| Backend | Go | Single binary, fast startup, cross-platform |
| Frontend | Vite + Svelte SPA | Lightweight SPA; maps cleanly to `embed.FS`; avoids SvelteKit SSR overhead irrelevant in a single-binary distribution |
| Database | SQLite | Single file, portable, no infrastructure |
| Distribution | Go `embed` package | Frontend assets baked into binary |

> **Note on frontend choice:** SvelteKit is a server-rendered framework. Its routing and layout conventions add build complexity with no benefit in a statically embedded context. Vite + Svelte SPA produces the same Svelte authoring experience with a simpler, flatter build that embeds cleanly into Go's `embed.FS`.

### 3.2 Distribution Model

The end user receives one executable file. On launch, it starts a local HTTP server and opens a browser tab. All data is stored in a `.db` file created in the same directory. No Docker, no Node, no pip.

```
maestro(.exe)          ← single binary
maestro.db             ← auto-created on first run
```

### 3.3 High-Level Data Flow

```
Azure DevOps Export (CSV/JSON)
        ↓
   Import Parser (Go)
        ↓
   SQLite (epics, features, sprints, metrics)
        ↓
   REST API (Go)
        ↓
   Vite + Svelte SPA
   ├── Gantt View (virtualized)
   ├── List View
   └── Health Dashboard
```

---

## 4. Data Model

### 4.1 Core Entities

**Epic**
```
id                  TEXT (Azure DevOps ID)
title               TEXT
description         TEXT
status              TEXT (New | Active | Resolved | Closed)
owner               TEXT
sprint_start        TEXT (sprint name)
sprint_end          TEXT (sprint name)
original_end_date   DATE  ← locked on import, never overwritten
committed_end_date  DATE  ← PM-adjustable
actual_end_date     DATE  ← set when status = Closed/Resolved
is_synthetic        BOOLEAN  ← true for auto-generated "Unassigned" Epic
created_at          DATETIME
updated_at          DATETIME
```

**Feature**
```
id                  TEXT (Azure DevOps ID)
epic_id             TEXT (FK → Epic; NULL if orphaned pre-assignment)
title               TEXT
description         TEXT
status              TEXT
owner               TEXT
sprint              TEXT
original_end_date   DATE  ← locked on import
committed_end_date  DATE  ← PM-adjustable
actual_end_date     DATE
story_points        INTEGER
created_at          DATETIME
updated_at          DATETIME
```

**Sprint**
```
id                  TEXT
name                TEXT
start_date          DATE
end_date            DATE
team                TEXT
capacity            INTEGER (story points; nullable)
source              TEXT (imported | generated)
```

**DateAuditLog**
```
id                  INTEGER PK
entity_type         TEXT (epic | feature)
entity_id           TEXT
changed_by          TEXT
old_date            DATE
new_date            DATE
delta_days          INTEGER  ← calendar days; positive = pushed out, negative = pulled in
reason              TEXT (optional)
changed_at          DATETIME
```

> **Critical design note:** `original_end_date` is written once on import and is never updated. `committed_end_date` is what the PM adjusts. Health metrics always compare `original_end_date` against `actual_end_date`. This prevents metric gaming through date manipulation.

> **`delta_days` is calendar days** for the POC. Business-day calculation (which better reflects sprint reality) is a v1.1 consideration if PMs flag it during M6 feedback.

> **`is_synthetic` flag:** The auto-generated "Unassigned" Epic is a real Epic record, not a UI fiction. This makes it queryable, sortable, and consistent with all other Epics. The flag lets the UI render it distinctly without special-casing the query layer.

### 4.2 Import Schema Mapping

Azure DevOps exports map to Maestro as follows:

| Azure DevOps Field | Maestro Field |
|---|---|
| `ID` | `id` |
| `Work Item Type` | determines Epic vs Feature |
| `Title` | `title` |
| `State` | `status` |
| `Assigned To` | `owner` |
| `Iteration Path` (last segment) | `sprint` |
| `Target Date` | `original_end_date` / `committed_end_date` |
| `Story Points` | `story_points` |
| `Parent` | `epic_id` (for Features) |
| `Capacity` (if present) | `sprint.capacity` |

**Hierarchy inference rules:**
- Work Item Type = `Epic` → creates an Epic record
- Work Item Type = `Feature` → creates a Feature record, linked to parent Epic via `Parent` field
- Features without a parent Epic are assigned to the synthetic "Unassigned" Epic (`is_synthetic = true`), which is created automatically on first import if needed. The Orphaned Feature Rate metric tracks this population.
- Features reassigned by the PM within Maestro are moved out of "Unassigned" and linked to a real Epic; their `epic_id` is updated accordingly.

**Sprint start date fallback chain:**

When determining a bar's start position in the Gantt, the following chain applies in order:

1. Sprint's `start_date` (preferred — sprint is the unit of time)
2. Sprint's `end_date` (if start is missing but end is known)
3. Import date (today at time of import)
4. `NULL` → bar rendered with a warning indicator; PM prompted to assign a sprint or date in the detail panel

Items in state 4 appear in the List view with a ⚠ flag but are not excluded from the Gantt — they render as a minimal stub bar at the leftmost position with a distinct visual treatment.

> **Export format variance — first technical spike:** Azure DevOps date formats vary by regional settings and export method. The parser must handle at minimum: `MM/DD/YYYY`, `YYYY-MM-DD`, `DD/MM/YYYY`, and ISO 8601 with timezone offsets. **This is M0** — collect exports from 3 internal teams before writing production parser code.

---

## 5. Feature Specifications

### 5.1 Import

**Entry point:** Drag-and-drop or file picker on the onboarding screen and via Settings → Import.

**Accepted formats:**
- Azure DevOps CSV export
- Azure DevOps JSON export

**Import flow:**
1. PM uploads file
2. Parser validates format and reports any unrecognizable rows
3. Preview screen shows: `X Epics`, `Y Features`, `Z Sprints detected`, `N items missing dates`, `M orphaned features`
4. PM confirms or cancels
5. On confirm: data written to SQLite, `original_end_date` locked, `committed_end_date` set to same value
6. PM proceeds to sprint configuration (first import) or returns to Gantt view (re-import)

**Error handling:**
- Missing required fields (ID, Title, Work Item Type) → row skipped, logged in import report
- Missing dates → apply sprint start date fallback chain (see §4.2); if all fallbacks exhausted, flagged for PM review
- Duplicate import (same IDs) → prompt: "Update existing items or create new project?"

---

### 5.2 Gantt View (Timeline)

The primary roadmap view. Rows are Epics with Features nested beneath them. The X-axis is time, organized by sprint.

**Display:**
- X-axis: sprints (labeled with sprint name and date range)
- Y-axis: Epic rows, expandable to show child Features. The synthetic "Unassigned" Epic appears at the bottom of the list, visually muted (no color pill, italic label), collapsed by default. It is always present if orphaned features exist.
- Bars represent `committed_end_date` range; start position derived from the sprint start date fallback chain (§4.2)
- Bar color encodes status: grey = not started, blue = in progress, green = completed on time, red = missed deadline
- Synthetic "Unassigned" Epic bar and its feature bars rendered in a distinct muted style
- Today line shown as a vertical marker

**Interactions:**
- Drag bar right edge to extend `committed_end_date`
- Drag bar body to shift start + end together
- Click bar to open detail panel: shows original date, committed date, net slip days, slip event count, status, owner
- Expand/collapse Epic rows
- Zoom: sprint-level, month-level, quarter-level

**Date adjustment rules:**
- Dragging updates `committed_end_date` only
- `original_end_date` is displayed in the detail panel but cannot be edited
- Every date change written to `DateAuditLog` with `delta_days` computed and stored as calendar days
- Snap-to-sprint is a **user toggle** (default: on), not an enforced constraint — PMs must never be forced to misrepresent a committed date

**Performance requirement:**
- The Gantt renderer must be **virtualized from M3** — only rows in the visible viewport are rendered to the DOM. A 500-item roadmap must not degrade scroll or drag performance. This is a build requirement, not a later optimization.

---

### 5.3 List View

A flat or grouped table of all work items, including orphaned features under the "Unassigned" group.

**Columns (default):**
- Title
- Type (Epic / Feature)
- Owner
- Sprint
- Original Date
- Committed Date
- Actual Date
- Slip Events (count of pushes)
- Status
- Health indicator (icon: on track / at risk / missed)

**Interactions:**
- Inline date editing: click `Committed Date` cell → date picker opens
- Sort by any column
- Filter by: Epic, Owner, Sprint, Status
- Group by: Epic (default), Sprint, Owner
- Export current view to CSV

---

### 5.4 Health Dashboard

A dedicated screen surfacing execution metrics. Refreshes on every import and date change.

---

**Metric 1: Deadline Hit Rate**
```
Deadline Hit Rate = (Features delivered on or before original_end_date)
                  / (Total completed Features) × 100
```
Displayed as a percentage with a trend line across the last 4 sprints.

---

**Metric 2: Sprint Load Index**
```
Sprint Load Index = (Committed story points in sprint) / (Sprint capacity)
```
- < 0.8 → underloaded
- 0.8–1.0 → healthy
- > 1.0 → overloaded

Shown per sprint as a bar chart. Sprints with no capacity set display "—" with an inline prompt: *"Add capacity in Settings → Sprints to enable this metric."* Missing capacity for one sprint never affects other sprints or metrics.

---

**Metric 3: Scope Creep Rate**
```
Scope Creep Rate = (Features added after sprint started)
                 / (Total features in sprint) × 100
```

---

**Metric 4: Slip Analysis**

Two complementary measures distinguish a feature that slipped once and stabilized from one being continuously deferred.

```
Slip Event Count (per feature) = Count of DateAuditLog entries where delta_days > 0

Net Slip Days (per feature)    = SUM of all delta_days for that feature
                                  positive = net delayed
                                  negative = net recovered (schedule recovery)
```

Dashboard-level aggregates:

```
Slip Event Rate    = (Features with 1+ slip events) / (Total features) × 100
Avg Net Slip Days  = Mean of Net Slip Days across features with any slip events
```

Distribution breakdown shown: features that slipped once / twice / 3+ times.

**Schedule Recoveries** — features with negative Net Slip Days are surfaced in a dedicated "Recoveries" callout, visually separated from the slip analysis section. A PM who pulls features in early should see that recognized, not buried in a table of slip events. The callout shows: count of features recovered, total days recovered, and the top 3 features by days recovered. This is a morale signal as much as a data point.

---

**Metric 5: Velocity Trend**
```
Story points completed per sprint over time (rolling 4-sprint average)
```

---

**Metric 6: Orphaned Feature Rate**
```
Orphaned Feature Rate = (Features assigned to synthetic "Unassigned" Epic)
                      / (Total features) × 100
```
A rate above 10% is flagged with a note: *"This reflects hierarchy hygiene in your Azure DevOps project. Assign orphaned features to Epics in the List view."* PM can reassign directly from the dashboard callout.

---

All metrics are filterable by Epic and by date range.

---

### 5.5 Sprint Management

**Detection priority:**
1. Sprint data present in Azure DevOps export → imported automatically, shown for PM confirmation
2. No sprint data in export → PM configures manually during onboarding

**Onboarding sprint step (first import):**
- If sprint data detected: confirmation table with imported sprints; PM can edit any row before confirming
- If no sprint data: manual config form — sprint length (1/2/3/4 weeks), first sprint start date, number of sprints to generate
- Both paths show a sprint preview table before proceeding
- Sprint source tagged: `imported` or `generated`

**Sprint fields:**
- Name
- Start date
- End date
- Capacity (story points) — optional; sourced from import if available
- Source (imported | generated)

**Capacity handling:**
- Imported from DevOps `Capacity` field if present; otherwise `NULL`
- Missing capacity shows "—" with inline prompt; never breaks other metrics
- Editable at any time in Settings → Sprints with no side effects

**Settings → Sprints:**
- Full sprint table with inline capacity editing
- Add sprint manually
- Toggle: snap-to-sprint boundaries (default: on)
- Reset to import defaults (requires confirmation)

---

## 6. UX / Design Principles

- **Minimal chrome.** The roadmap is the UI. Navigation, controls, and menus should not compete with the content.
- **Data density without clutter.** PMs need to see many items at once. Typography and spacing must support scanning, not reading.
- **Sprint as the default unit of time, not a constraint.** Snap-to-sprint aids workflow but is always overridable. PMs must never be forced to misrepresent a committed date.
- **Metrics are always visible.** Health indicators appear inline (Gantt bar color, list health icon, slip event count) without requiring navigation to the dashboard.
- **No data loss.** Original dates are sacred. The system preserves them forever and makes them visible alongside committed dates.
- **Unavailable metrics are explained, not hidden.** When a metric cannot be computed, show "—" and tell the PM exactly what to do to enable it.
- **Recoveries are celebrated, not buried.** Schedule improvements are a positive signal and should be surfaced as such — not hidden inside slip tables.

---

## 7. Technical Specifications

### 7.1 API Endpoints (REST)

```
POST   /api/import                   Upload and parse DevOps export
GET    /api/epics                    List all epics (with features nested)
GET    /api/epics/:id                Single epic detail
PATCH  /api/epics/:id/date           Update committed_end_date
GET    /api/features                 List all features
GET    /api/features/:id             Single feature detail
PATCH  /api/features/:id/date        Update committed_end_date
PATCH  /api/features/:id/epic        Reassign feature to a different epic
GET    /api/sprints                  List all sprints
POST   /api/sprints                  Create sprint
PATCH  /api/sprints/:id              Update sprint (dates, capacity, name)
GET    /api/metrics                  Return all health metrics
GET    /api/metrics/slip/:id         Slip event history for a specific feature
GET    /api/metrics/recoveries       Features with negative net slip days
GET    /api/audit                    Date change audit log
```

### 7.2 Binary Build

```bash
# Build frontend first
cd frontend && npm run build

# Build for current platform
go build -o maestro .

# Cross-compile
GOOS=windows GOARCH=amd64 go build -o maestro.exe .
GOOS=darwin  GOARCH=amd64 go build -o maestro-mac .
GOOS=linux   GOARCH=amd64 go build -o maestro-linux .
```

Frontend compiled assets embedded via:

```go
//go:embed frontend/dist
var frontendAssets embed.FS
```

> Vite outputs to `dist/` by default. Update embed path if changed.

### 7.3 Storage

SQLite via `modernc.org/sqlite` (pure Go, no CGO dependency). Database file created at `./maestro.db` on first run.

### 7.4 Binary Size Budget and Checkpoints

Target: < 25 MB total binary.

**M0 checkpoint:** Measure `modernc.org/sqlite` baseline before writing any application code. A minimal Go binary with only the SQLite driver linked should be measured and recorded. If this baseline alone exceeds 10 MB, the 25 MB target needs to be reassessed before M1 begins — not after M5.

**M1 checkpoint:** Backend only (parser + schema + REST API, no embedded frontend). If > 10 MB, assess mitigation before adding frontend bundle.

**M3 checkpoint:** Full binary with frontend embedded. If > 20 MB, apply: tree-shaking on frontend bundle, lazy-loading for Gantt renderer, gzip compression on embedded assets.

Track on each milestone build:

```bash
go build -o maestro . && ls -lh maestro
```

---

## 8. Success Metrics (POC)

| Metric | Target | Definition |
|---|---|---|
| Import-to-roadmap time | < 30 seconds for a 500-item export | Timed from file drop to Gantt render |
| Binary size | < 25 MB | Measured at M0, M1, M3, M5 checkpoints |
| Startup time | < 2 seconds to first browser tab | Cold start on reference machine |
| PM adoption | At least 2 PMs complete a full sprint cycle | Import at least once, adjust at least 3 dates, return in a subsequent session |
| Date adjustment UX | Net positive | At least 2 of 2 pilot PMs rate it easier than current tool |
| Gantt performance | No visible jank at 200+ items | Smooth scroll and drag; no dropped frames |

---

## 9. Milestones

| Milestone | Scope | Target |
|---|---|---|
| M0 — Spike | Collect exports from 3 internal teams; validate date format variance; measure SQLite baseline binary size; finalize parser strategy | Week 1 |
| M1 — Core | Import parser, SQLite schema, REST API, backend-only binary size checkpoint | Week 1–2 |
| M2 — List View | List view with inline date editing, slip event column, orphaned feature reassignment, filters | Week 2–3 |
| M3 — Gantt View | Virtualized timeline, drag-to-adjust, snap toggle, "Unassigned" Epic row, NULL bar stubs | Week 3–5 |
| M4 — Health Dashboard | All 6 metrics, slip analysis, recoveries callout, per-sprint charts | Week 5–6 |
| M5 — Binary | Embedded frontend, single-binary distribution, final size check | Week 6 |
| M6 — POC Review | Internal user testing, structured feedback, strategic research interviews | Week 7–8 |

---

## 10. Open Questions

1. **Multiple projects** — Should the POC support multiple Azure DevOps project imports, or one project per database file?
2. **Hierarchy depth** — Azure DevOps can have User Stories under Features under Epics. Should the POC display User Story level, or stop at Feature?
3. **PM identity** — Since there is no auth in the POC, how should we attribute date changes in the audit log? Assumption: simple name entry on first launch.

---

## 11. M6 Strategic Research Questions

The M6 review must answer the question that determines Maestro's product direction:

> **Do PMs want a separate tool, or do they want better views inside Azure DevOps?**

"Separate tool" → Maestro is a new product category. The roadmap toward unified PM + engineering ticketing is the right path.

"Better DevOps views" → Maestro should become an Azure DevOps extension or plugin. The architecture and go-to-market change entirely.

The POC will not resolve this. The M6 interviews must ask directly:

- Did you find yourself switching between Maestro and Azure DevOps during the sprint? How often?
- Did engineers ever ask to see your Maestro roadmap, or did you always translate it back for them?
- If Maestro's Gantt and health metrics were available natively inside Azure DevOps, would you still want a separate tool? Why?
- What would make you pay for Maestro as a product versus asking your DevOps admin to install a plugin?

These answers set the direction for v1.0.

---

## 12. Appendix

### A. Glossary

| Term | Definition |
|---|---|
| Epic | A large body of work spanning multiple sprints, containing multiple Features |
| Feature | A deliverable within an Epic, assignable to a single sprint |
| Orphaned Feature | A Feature assigned to the synthetic "Unassigned" Epic due to missing parent in the import |
| Synthetic Epic | The auto-generated "Unassigned" Epic (`is_synthetic = true`); a real DB record, not a UI fiction |
| Original End Date | The date captured on import; never modified; used for all health calculations |
| Committed End Date | The PM-adjusted delivery target; shown on the Gantt |
| Actual End Date | The date the item moved to Closed/Resolved status |
| Sprint Load Index | Ratio of committed story points to sprint capacity |
| Deadline Hit Rate | Percentage of Features delivered by their original end date |
| Slip Event | A single instance of committed_end_date being pushed out (delta_days > 0) |
| Net Slip Days | Sum of all delta_days for a feature; positive = net delayed, negative = net recovered |
| Slip Event Rate | Percentage of features with at least one slip event |
| Schedule Recovery | A feature with negative Net Slip Days; surfaced in the Recoveries callout |

### B. Competitive Landscape

| Tool | Strength | Gap Maestro Fills |
|---|---|---|
| Productboard | Great prioritization | No execution health metrics; no DevOps sync |
| Jira Roadmaps | Native to Jira | Not available to Azure DevOps users |
| Azure DevOps Boards | Source of truth | Not a PM-facing roadmap tool |
| Linear | Beautiful UX | Not designed for enterprise PM workflows |
| Spreadsheets | Flexible | Manual, no metrics, always stale |

### C. Execution Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Azure DevOps export format variance | High | High | M0 spike; collect from 3 internal teams before writing production parser |
| SQLite binary overhead exceeds budget headroom | Medium | High | M0 baseline measurement; reassess 25 MB target before M1 if > 10 MB |
| Gantt performance at 500+ items | Medium | High | Virtualized renderer required from M3; not a later optimization |
| Sprint capacity missing from exports | High | Medium | Capacity nullable; missing shows "—" with inline prompt; never breaks other metrics |
| Binary size > 25 MB | Medium | Low | M1 and M3 checkpoints; tree-shake and lazy-load if at risk |
| PMs adjust dates off sprint boundaries | Medium | Low | Snap-to-sprint is a toggle; PMs never forced to misrepresent committed dates |
| Pilot PMs don't return after first session | Medium | High | Define adoption precisely (§8); instrument session return from day one |
| Tool vs. plugin question goes unanswered at M6 | Low | High | Embed strategic questions explicitly in M6 interview guide |
