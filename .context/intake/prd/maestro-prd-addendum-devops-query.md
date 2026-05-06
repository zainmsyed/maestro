# Maestro — PRD Addendum: Prescribed Azure DevOps Query + Stories Support

**Addendum to:** maestro-prd.md v0.3
**Version:** 1.1
**Status:** Draft
**Last Updated:** 2026-05-06

---

## Changelog

| Version | Changes |
|---|---|
| 1.1 | Added Story as a first-class entity (third hierarchy level). Added `Title 3` column to prescribed query. Added Story data model. Defined Story type name variants across process templates. Added PM date insertion for items missing Target Date — Epics, Features, and Stories can all have dates set directly in Maestro without going back to DevOps. Updated parser for three-level hierarchy. Updated Gantt and List view specs for three-level nesting. |
| 1.0 | Initial addendum — prescribed query, parser spec, common problems |

---

## Purpose

This addendum defines the prescribed Azure DevOps query that PMs must create before using Maestro, and specifies how Stories are handled as a third hierarchy level beneath Features. It also defines the PM date insertion workflow for work items that have no Target Date in DevOps.

By standardizing the query structure and column set, Maestro can parse a consistent, predictable CSV every time — eliminating the column variance risk identified in the PRD risk register.

The PM sets this up once. After that, running and exporting the query takes under a minute.

---

## Background: Why a Prescribed Query

Azure DevOps CSV exports only include columns that were explicitly added to the query before exporting. There is no default export format that includes all the fields Maestro needs. Without a prescribed query, every PM would export a different set of columns, making a reliable parser impossible.

Additionally, the parent-child hierarchy between Epics, Features, and Stories is only preserved in a specific query type. A flat list query flattens hierarchy entirely. Only a **Tree of Work Items** query encodes parent-child relationships in the exported CSV via indented title columns (`Title 1`, `Title 2`, `Title 3`), which is the structure Maestro's parser is built to read.

---

## Section 1: Hierarchy Model

Maestro supports three levels of work item hierarchy:

```
Epic
└── Feature
    └── Story
```

| Level | Azure DevOps Type (Agile) | Azure DevOps Type (Scrum) | Azure DevOps Type (CMMI) |
|---|---|---|---|
| Epic | Epic | Epic | Epic |
| Feature | Feature | Feature | Feature |
| Story | User Story | Product Backlog Item | Requirement |

> **Key difference from Epic and Feature:** Epic and Feature type names are consistent across all process templates. Story type names vary. The prescribed query and parser must account for all three variants. When Maestro encounters any of `User Story`, `Product Backlog Item`, or `Requirement` in the `Work Item Type` column, it treats that row as a Story.

---

## Section 2: One-Time Setup Guide

### Step 1 — Navigate to Queries

In your Azure DevOps project, go to:

**Boards → Queries → New Query**

---

### Step 2 — Set Query Type

In the Query Editor, change the **Type of query** dropdown from the default ("Flat list of work items") to:

> **Tree of work items**

This is essential. It generates the `Title 1` / `Title 2` / `Title 3` column structure that encodes Epic → Feature → Story hierarchy in the CSV export.

---

### Step 3 — Set Filters

Configure the query clauses as follows. Use **Or** between the three type rows:

| And/Or | Field | Operator | Value |
|---|---|---|---|
| — | Work Item Type | = | Epic |
| Or | Work Item Type | = | Feature |
| Or | Work Item Type | = | User Story |

> **Scrum teams:** Replace `User Story` with `Product Backlog Item`
> **CMMI teams:** Replace `User Story` with `Requirement`
> **Custom process templates:** Use whatever type sits below Feature in your hierarchy

Set the link type to **Child Of** between each level.

---

### Step 4 — Add Required Columns

Click **Column Options**. Add the following columns in this order:

| # | Column name | Why it's needed |
|---|---|---|
| 1 | ID | Unique identifier for each work item |
| 2 | Work Item Type | Distinguishes Epics, Features, and Stories |
| 3 | Title 1 | Epic title |
| 4 | Title 2 | Feature title |
| 5 | Title 3 | Story title |
| 6 | State | Work item status |
| 7 | Assigned To | Owner |
| 8 | Iteration Path | Sprint assignment |
| 9 | Story Points | Effort estimate |
| 10 | Target Date | Planned completion date — may be empty; PMs can set this in Maestro |
| 11 | Area Path | Team/area context |

> **On Target Date:** It is expected and normal that many Stories — and some Features — will not have a Target Date set in Azure DevOps. This is not an error. Maestro surfaces these items after import and lets the PM assign dates directly, without going back to DevOps.

---

### Step 5 — Save the Query

Save the query as:

> **Maestro Export**

Save to **Shared Queries** so any PM on the team can run it.

---

### Step 6 — Export to CSV

Every time you want to update your Maestro roadmap:

1. Open **Boards → Queries → Shared Queries → Maestro Export**
2. Run the query
3. Click the **⋯ actions menu** (top right of results)
4. Select **Export to CSV**
5. Save the file
6. Open Maestro and drag the file onto the import screen

Steps 1–5 are one-time setup. Only step 6 repeats.

---

## Section 3: What the CSV Looks Like

A correctly configured export produces a CSV with this structure:

```
ID,Work Item Type,Title 1,Title 2,Title 3,State,Assigned To,Iteration Path,Story Points,Target Date,Area Path
"1042","Epic","Auth & Identity","","","Active","Sara Chen <sara@co.com>","Project\Team\Sprint 12","","2026-04-25","Project\Platform"
"1055","Feature","","SSO Integration","","Active","Marcus Webb <marcus@co.com>","Project\Team\Sprint 12","","2026-04-18","Project\Platform"
"","User Story","","","Set up SAML config","Active","Marcus Webb <marcus@co.com>","Project\Team\Sprint 12","3","2026-04-11","Project\Platform"
"","User Story","","","Test IdP integration","Active","Priya Nair <priya@co.com>","Project\Team\Sprint 12","2","","Project\Platform"
"1056","Feature","","MFA Enforcement","","Active","Sara Chen <sara@co.com>","Project\Team\Sprint 13","","2026-05-02","Project\Platform"
"","User Story","","","Add TOTP support","Active","Sara Chen <sara@co.com>","Project\Team\Sprint 13","5","","Project\Platform"
```

Key structural rules:

| Row type | Title 1 | Title 2 | Title 3 |
|---|---|---|---|
| Epic | populated | empty | empty |
| Feature | empty | populated | empty |
| Story | empty | empty | populated |

- Stories belong to the last Feature row above them
- Features belong to the last Epic row above them
- Empty `Target Date` is valid and expected — Maestro handles it
- Empty `Story Points` is valid — Sprint Load Index shows "—" for that sprint

---

## Section 4: Data Model — Story Entity

The following table is added to the Maestro SQLite schema alongside Epic, Feature, Sprint, and DateAuditLog.

**Story**
```
id                  TEXT (Azure DevOps ID; may be empty for new items)
feature_id          TEXT (FK → Feature)
title               TEXT
description         TEXT
status              TEXT (New | Active | Resolved | Closed)
owner               TEXT
sprint              TEXT
story_points        INTEGER (nullable)
original_end_date   DATE  ← locked on import if Target Date present; NULL if not
committed_end_date  DATE  ← PM-settable in Maestro; NULL until PM assigns
actual_end_date     DATE  ← set when status = Closed/Resolved
date_source         TEXT  ← 'imported' | 'pm_assigned' | 'inherited'
created_at          DATETIME
updated_at          DATETIME
```

> **`date_source` field:** Tracks where the committed date came from. `imported` = came from Target Date in DevOps. `pm_assigned` = PM set it directly in Maestro. `inherited` = inferred from parent Feature or sprint end date. This field lets the health dashboard distinguish "PM made a deliberate date decision" from "this date is a guess."

> **`original_end_date` for Stories:** If `Target Date` is empty in the export, `original_end_date` is `NULL` on import. It is set the first time the PM assigns a date in Maestro — that first PM-assigned date becomes the `original_end_date` and is then locked, exactly like imported dates. Subsequent changes update `committed_end_date` only.

---

## Section 5: PM Date Insertion

Items that arrive from DevOps without a Target Date are not broken — they are surfaced to the PM for date assignment directly within Maestro. The PM never needs to go back to Azure DevOps to add dates.

### 5.1 Post-Import Date Assignment Screen

After import, if any items are missing dates, Maestro shows a **Date Assignment** screen before the PM lands on the Gantt:

```
┌─────────────────────────────────────────────────────────┐
│  3 items need dates                                     │
│  You can set them now or come back to them later.       │
│                                                         │
│  ○ Test IdP integration       Sprint 12   [pick date]  │
│  ○ Add TOTP support           Sprint 13   [pick date]  │
│  ○ Webhook payload signing    Sprint 15   [pick date]  │
│                                                         │
│  [Set dates now]              [Skip — I'll do it later] │
└─────────────────────────────────────────────────────────┘
```

- Each item shows its sprint (inferred from `Iteration Path`) as context for picking a date
- Date picker snaps to sprint boundaries by default (toggle off if needed)
- PM can skip and assign dates later — undated items appear with a ⚠ indicator in all views

### 5.2 Date Entry Points

PMs can assign or adjust dates from three places:

**1. Gantt view** — undated Stories/Features render as a minimal stub bar at their sprint's position. Clicking the stub opens the detail panel with a date picker. Once a date is set the stub becomes a full bar.

**2. List view** — the `Committed Date` column shows "— set date" for undated items. Clicking it opens an inline date picker.

**3. Detail panel** — available from both Gantt and List. Shows `original_end_date` (locked, or "not set" if never assigned), `committed_end_date` (editable), sprint context, and audit log.

### 5.3 First Date Sets Original

The first time a PM assigns a date to an item that had no Target Date in DevOps:

- `original_end_date` = the date they entered (now locked)
- `committed_end_date` = the same date
- `date_source` = `pm_assigned`
- A `DateAuditLog` entry is created: `old_date = NULL`, `new_date = assigned date`, `delta_days = 0`

Subsequent changes follow the standard audit flow — only `committed_end_date` updates, `original_end_date` stays locked.

### 5.4 Inherited Dates

If a PM does not assign a date and the item has a valid sprint:

- `committed_end_date` = sprint `end_date`
- `date_source` = `inherited`
- These items are visually distinguished in the Gantt (lighter bar opacity, no lock icon)
- They are excluded from Deadline Hit Rate calculations — inherited dates are not commitments

---

## Section 6: Maestro Parser Specification (Updated for Three Levels)

### 6.1 File Validation

On import, before parsing any rows, the parser validates:

1. File is UTF-8 encoded
2. First row is a header row
3. Header contains all 11 required columns (any order; normalised after detection)
4. At least one data row exists

If validation fails, import is rejected with a specific error message identifying which column is missing.

### 6.2 Header Normalisation

Column names matched case-insensitively with whitespace trimmed:

| Canonical name | Also accepted |
|---|---|
| `Title 1` | `Title1`, `title 1` |
| `Title 2` | `Title2`, `title 2` |
| `Title 3` | `Title3`, `title 3` |
| `Iteration Path` | `Iteration`, `Sprint` |
| `Story Points` | `StoryPoints`, `Effort`, `Size` |
| `Target Date` | `TargetDate`, `Due Date`, `Finish Date` |
| `Assigned To` | `AssignedTo`, `Owner` |

### 6.3 Hierarchy Parsing

The parser tracks current Epic and current Feature using two pointers:

```
current_epic    = nil
current_feature = nil

for each row:

  type = Work Item Type value (normalised — see §6.4)

  if type == "epic":
    → create Epic record
    → current_epic = this epic
    → current_feature = nil  ← reset; next features belong to this epic

  elif type == "feature":
    → create Feature record
    → if current_epic != nil: feature.epic_id = current_epic.id
    → else: feature.epic_id = synthetic_unassigned_epic.id
    → current_feature = this feature

  elif type == "story":
    → create Story record
    → if current_feature != nil: story.feature_id = current_feature.id
    → else: story.feature_id = synthetic_unassigned_feature.id
    → (synthetic unassigned feature is auto-created under unassigned epic if needed)

  else:
    → unknown type — log warning, skip row
```

### 6.4 Work Item Type Normalisation

The parser maps all known Story type variants to the internal `story` type:

| Raw value in CSV | Internal type |
|---|---|
| `Epic` | `epic` |
| `Feature` | `feature` |
| `User Story` | `story` |
| `Product Backlog Item` | `story` |
| `Requirement` | `story` |

Any other value is logged as unknown and the row is skipped.

### 6.5 Date Parsing

`Target Date` field parsed with format detection — not assumed format — because Azure DevOps exports dates in the regional format of the exporting user.

Formats attempted in order:

1. `2006-01-02` (ISO 8601)
2. `01/02/2006` (US — MM/DD/YYYY)
3. `02/01/2006` (UK/EU — DD/MM/YYYY)
4. `2006/01/02` (ISO with slashes)
5. `January 2, 2006` (long form)
6. `2-Jan-2006` (abbreviated month)
7. `2006-01-02T15:04:05Z07:00` (ISO 8601 with timezone)

Empty `Target Date` → `original_end_date = NULL`; item added to the post-import date assignment list.

Unparseable non-empty value → log warning, treat as empty, add to date assignment list.

> **Ambiguity:** Formats 2 and 3 are ambiguous for days ≤ 12. Parser defaults to US (MM/DD/YYYY). Confirm regional format during M0.

### 6.6 Iteration Path Parsing

Extracts last path segment as sprint name:

```
"Project\Team\Sprint 14"  → "Sprint 14"
"Project\Sprint 14"       → "Sprint 14"
"Sprint 14"               → "Sprint 14"
""                         → nil → apply sprint fallback chain (PRD §4.2)
```

### 6.7 Assigned To Parsing

```
"Sara Chen <sara@company.com>"  → display_name = "Sara Chen", email = "sara@company.com"
"Sara Chen"                     → display_name = "Sara Chen", email = nil
""                              → owner = nil
```

### 6.8 Story Points Parsing

- Empty → `story_points = nil`
- Numeric → parse as integer; decimals rounded to nearest integer
- Non-numeric → log warning, set nil

### 6.9 Import Report

```
✓  Epics parsed:                4
✓  Features parsed:             10
✓  Stories parsed:              24
✓  Sprints detected:            6
⚠  Items missing Target Date:   8  (will appear in date assignment screen)
⚠  Items missing sprint:        1  (flagged for review)
⚠  Orphaned features:           0
⚠  Orphaned stories:            0
⚠  Rows skipped (unknown type): 0
⚠  Date format detected:        US (MM/DD/YYYY)
```

---

## Section 7: View Updates for Stories

### 7.1 Gantt View

Three levels of nesting:

```
▾ Epic — Auth & Identity                    [════════════════]
  ▾ Feature — SSO Integration               [════════]
      Story — Set up SAML config            [════]
      Story — Test IdP integration          [stub — no date]
  ▾ Feature — MFA Enforcement                        [════════]
      Story — Add TOTP support              [stub — no date]
```

- Epics expand/collapse to show/hide Features
- Features expand/collapse to show/hide Stories
- Stories default to **collapsed** at the Feature level — PM expands if needed
- Stories without dates render as a stub bar with ⚠ indicator at the sprint's start position
- Story bars are thinner than Feature bars (8px height vs 10px for Features, 22px for Epics) to support density
- Story bars are not draggable unless a date has been assigned — clicking a stub opens the date picker instead

### 7.2 List View

Three grouping levels: Epic → Feature → Story.

Story rows indented further than Feature rows. Additional column at Story level:

| Column | Notes |
|---|---|
| Date Source | `imported` / `pm assigned` / `inherited` — shown as a small badge |

Story rows show "— set date" in the Committed Date column if undated, with click-to-assign behaviour.

### 7.3 Health Dashboard

**Metrics updated for Stories:**

- **Deadline Hit Rate** — calculated at Story level as well as Feature level. Story-level rate gives more granular signal since Stories complete within a sprint, not across sprints. Both rates shown; PM can toggle between them.
- **Sprint Load Index** — unchanged; already uses story points which are set at Story level
- **Scope Creep Rate** — Stories added to a sprint after it started now count, not just Features. This is more accurate since Story-level scope changes are the most common form of mid-sprint creep.
- **Slip Analysis** — applies to Stories with PM-assigned or imported dates. Stories with `date_source = inherited` are excluded from slip calculations — inherited dates are not commitments.
- **Orphaned Story Rate** — new metric alongside Orphaned Feature Rate:

```
Orphaned Story Rate = (Stories with no parent Feature) / (Total stories) × 100
```

---

## Section 8: Common Export Problems and Fixes

| Problem | Symptom in Maestro | Fix |
|---|---|---|
| Wrong query type (flat list) | All items import as Epics; no Features or Stories | Re-run as "Tree of work items" |
| `Title 3` column missing | Stories not detected | Add Title 3 via Column Options |
| Wrong Story type name | Stories skipped as unknown type | Match type name to process template (User Story / Product Backlog Item / Requirement) |
| `Target Date` empty on most Stories | 20+ items in date assignment screen | Expected — use date assignment screen or set dates in DevOps first |
| `Iteration Path` missing | All sprints show as unassigned | Add Iteration Path via Column Options |
| `Story Points` empty | Sprint Load Index shows "—" | Set estimates in DevOps or accept metric unavailability for now |
| Features appear at Story level | Title 2 and Title 3 both populated on Feature rows | Query not correctly set to tree type |
| Regional date format mismatch | Dates import as wrong month/day | Confirm regional format; re-export from ISO-format machine if needed |

---

## Section 9: PRD Sections Updated by This Addendum

| PRD Section | Change |
|---|---|
| §4.1 Core Entities | Story entity added (Addendum §4) |
| §4.2 Import Schema Mapping | Hierarchy via `Title 1`/`Title 2`/`Title 3`; full parser spec in Addendum §6 |
| §5.1 Import | JSON format removed from POC scope; post-import date assignment screen added |
| §5.2 Gantt View | Three-level nesting; stub bars for undated items; Stories collapsed by default |
| §5.3 List View | Three-level grouping; date source badge; click-to-assign for undated items |
| §5.4 Health Dashboard | Story-level Deadline Hit Rate; Story-level scope creep; Orphaned Story Rate |
| §9 Milestones — M0 | Validate three-level tree query export; confirm Story type name for internal teams |
