# Maestro — PRD Addendum: Prescribed Azure DevOps Query + Parser Specification

**Addendum to:** maestro-prd.md v0.3
**Version:** 2.0
**Status:** Draft
**Last Updated:** 2026-05-06

---

## Changelog

| Version | Changes |
|---|---|
| 2.0 | Complete rewrite based on analysis of three real Azure DevOps CSV exports. Parent ID established as the sole hierarchy signal. Title column inference removed entirely. Parser rebuilt around real-world variance: float Parent IDs, 2–5 title columns, datetime Target Date format, multi-segment Iteration Path, Assigned To name format variants, Blocked state, missing Story Points. Query installer flow added. JSON import removed from scope. |
| 1.1 | Added Story as third hierarchy level. Added PM date insertion. |
| 1.0 | Initial addendum — prescribed query, parser spec, common problems. |

---

## Purpose

This addendum defines the prescribed Azure DevOps query PMs must run before using Maestro, and specifies exactly how the Go parser handles real-world CSV export variance. It supersedes the import schema mapping in PRD §4.2.

The approach is built on one principle validated across three real exports: **`Parent` ID is the only reliable hierarchy signal.** Everything else — title columns, row order, column count — varies unpredictably across teams. The parser ignores title structure entirely and reconstructs hierarchy purely from the `Parent` foreign key.

---

## Section 1: What the Parser Relies On

Before anything else, this is the full picture of what Maestro can and cannot count on from a real DevOps export.

### Reliable — Parser Can Depend On These

| Field | Why Reliable |
|---|---|
| `ID` | Always present; primary key |
| `Work Item Type` | Always present; determines Epic / Feature / Story |
| `Parent` | Empty for Epics, populated for Features and Stories; clean foreign key regardless of team structure |
| `State` | Always present; values vary but field is always there |

### Soft-Reliable — Present in Most Exports, Handle Absence Gracefully

| Field | Variance | Fallback |
|---|---|---|
| `Iteration Path` | Always present but segment depth varies (2–5 segments); `Backlog` and `Archive` appear as non-sprint values | `NULL` sprint → item flagged as unscheduled |
| `Target Date` | Format varies (`YYYY-MM-DD`, `MM/DD/YYYY`, `MM/DD/YYYY HH:MM:SS AM/PM`); often empty on Features and Epics | `NULL` → post-import date assignment screen |
| `Assigned To` | Format varies (`First Last <email>` vs `Last, First <email>`); sometimes email only | Parse best effort; `NULL` is valid |
| Title columns | Count varies (1–5+); titles may appear in any column; unreliable for hierarchy | Scan all title columns; take last non-empty value as display title |

### Unreliable — Do Not Depend On

| Field | Reality |
|---|---|
| `Story Points` | Missing in majority of real exports; many teams use hours, t-shirt sizes, or nothing |
| `Area Path` | Segment depth varies; used for filtering only, never for hierarchy or metrics |
| Title column position | No consistent mapping between column index and hierarchy level |
| `Description`, `Tags`, `Priority` | Not in prescribed query; too variable to parse reliably |

---

## Section 2: Required vs Optional Columns

### Absolute Required (Import Rejected Without These)

```
Parent
ID
Work Item Type
```

These three are the skeleton. Without them the parser cannot create entities or link them.

### Strongly Recommended (Import Proceeds With Warnings If Missing)

```
[Any title column]   → display names; fallback to "Untitled {Type} {ID}"
Iteration Path       → sprint placement; missing = all items unscheduled
Target Date          → bar end dates; missing = date assignment screen
State                → health metrics; missing = all items treated as "New"
Assigned To          → ownership display; missing = "Unassigned"
```

### Optional (Enrichment Only)

```
Area Path            → List view filter dimension only
Story Points         → Sprint Load Index; permanently shows "—" if absent
```

---

## Section 3: One-Time Query Setup Guide

### Option A — Automatic Install (Recommended)

On the Maestro onboarding screen, enter:
- Azure DevOps org URL (e.g. `https://dev.azure.com/myorg`)
- Project name
- Personal Access Token (PAT) with **Work Items (Read)** scope

Maestro calls the Azure DevOps REST API to create the "Maestro Export" query directly in your Shared Queries folder. The PM never touches the query editor.

If the API call fails (org restrictions, insufficient PAT permissions), Maestro falls back to Option B and shows the manual setup guide.

---

### Option B — Manual Setup

**Step 1 — Navigate to Queries**

**Boards → Queries → New Query**

**Step 2 — Set Query Type**

Change the type dropdown to:

> **Tree of work items**

**Step 3 — Set Filters**

| And/Or | Field | Operator | Value |
|---|---|---|---|
| — | Work Item Type | = | Epic |
| Or | Work Item Type | = | Feature |
| Or | Work Item Type | = | User Story |

> **Scrum teams:** Replace `User Story` with `Product Backlog Item`
> **CMMI teams:** Replace `User Story` with `Requirement`
> **Custom templates:** Use the type name that sits below Feature in your hierarchy

Set link type to **Child Of** between levels.

**Step 4 — Add Required Columns**

Click **Column Options** and add in this order:

| # | Column | Required? |
|---|---|---|
| 1 | Parent | ✓ Required |
| 2 | ID | ✓ Required |
| 3 | Work Item Type | ✓ Required |
| 4 | Title 1 | Recommended |
| 5 | Title 2 | Recommended |
| 6 | Title 3 | Recommended |
| 7 | State | Recommended |
| 8 | Assigned To | Recommended |
| 9 | Iteration Path | Recommended |
| 10 | Story Points | Optional |
| 11 | Target Date | Recommended |
| 12 | Area Path | Optional |

> **On title columns:** Add as many title columns as your hierarchy depth requires. Maestro scans all of them and takes the last non-empty value as the display title regardless of which column it appears in. If your team has 5 levels, add Title 1 through Title 5.

**Step 5 — Save the Query**

Save as **"Maestro Export"** to **Shared Queries**.

**Step 6 — Export to CSV**

Every subsequent import:
1. **Boards → Queries → Shared Queries → Maestro Export**
2. Run the query
3. **⋯ → Export to CSV**
4. Drag the file into Maestro

Steps 1–5 are one-time. Only step 6 repeats.

---

## Section 4: WIQL Definition

This is the exact query Maestro installs via the API in Option A, and what the manual setup produces in Option B.

```sql
SELECT
  [System.Id],
  [System.Parent],
  [System.WorkItemType],
  [System.Title],
  [System.State],
  [System.AssignedTo],
  [System.AreaPath],
  [System.IterationPath],
  [Microsoft.VSTS.Scheduling.StoryPoints],
  [Microsoft.VSTS.Scheduling.TargetDate]
FROM workItemLinks
WHERE
  (
    [Source].[System.WorkItemType] IN ('Epic', 'Feature', 'User Story',
      'Product Backlog Item', 'Requirement')
  )
  AND [System.Links.LinkType] = 'System.LinkTypes.Hierarchy-Forward'
MODE (Recursive)
ORDER BY [System.Id]
```

---

## Section 5: Go Parser Specification

This section is the authoritative spec for the M1 import parser. It is written against real export data from three distinct Azure DevOps projects.

### 5.1 File Validation

Before parsing any rows:

1. File must be UTF-8 encoded
2. First row must be a header row
3. Header must contain `Parent`, `ID`, and `Work Item Type` — if any of these three are missing, reject import with a specific error message
4. At least one data row must exist

All other columns are detected opportunistically — their absence triggers warnings in the import report, not rejection.

### 5.2 Header Detection

Match column names case-insensitively with whitespace trimmed. Accept these aliases:

| Canonical | Also Accepted |
|---|---|
| `Parent` | `Parent ID`, `ParentId`, `parent` |
| `ID` | `Id`, `Work Item ID`, `WorkItemId` |
| `Work Item Type` | `WorkItemType`, `Type` |
| `Iteration Path` | `IterationPath`, `Sprint`, `Iteration` |
| `Story Points` | `StoryPoints`, `Effort`, `Size`, `Points` |
| `Target Date` | `TargetDate`, `Due Date`, `Finish Date`, `DueDate` |
| `Assigned To` | `AssignedTo`, `Owner`, `Assigned` |
| `Area Path` | `AreaPath`, `Area` |

Any column matching `Title \d+` or `Title\d+` is collected as a title column regardless of number.

### 5.3 ID and Parent Normalisation

Azure DevOps exports numeric IDs as floats (`500588.0`). Before any lookup or storage:

```go
// Strip decimal from ID and Parent fields
func normaliseID(raw string) string {
    raw = strings.TrimSpace(raw)
    if strings.Contains(raw, ".") {
        raw = strings.Split(raw, ".")[0]
    }
    return raw
}
```

Empty `Parent` field → item is a root-level entity (Epic or orphan).

### 5.4 Work Item Type Normalisation

Map all known variants to internal types:

| Raw CSV Value | Internal Type |
|---|---|
| `Epic` | `epic` |
| `Feature` | `feature` |
| `User Story` | `story` |
| `Product Backlog Item` | `story` |
| `Requirement` | `story` |

Any unrecognised Work Item Type → row logged as unknown, skipped, counted in import report. Parser does not fail — it continues.

### 5.5 Hierarchy Reconstruction via Parent ID

This is the complete hierarchy algorithm. It does not use title columns, row order, or any other signal.

```go
// Pass 1: Build entity map
entityMap := map[string]Entity{}
for each row:
    id = normaliseID(row["ID"])
    parentID = normaliseID(row["Parent"])
    itemType = normaliseType(row["Work Item Type"])
    title = extractTitle(row)  // see §5.6
    entity = Entity{id, parentID, itemType, title, ...}
    entityMap[id] = entity

// Pass 2: Link entities
for each entity in entityMap:
    if entity.parentID == "":
        if entity.type == "epic":
            → root epic, no parent needed
        else:
            → orphan, assign to synthetic Unassigned bucket
            → flag in import report
    else:
        parent = entityMap[entity.parentID]
        if parent == nil:
            → parent not in this export
            → orphan, assign to Unassigned bucket
            → flag in import report
        else:
            entity.parentRef = parent
```

Two-pass approach means row order in the CSV is irrelevant. A Story can appear before its parent Feature and still link correctly.

### 5.6 Title Extraction

Collect all columns matching `Title \d+` pattern. Scan them in ascending column order and return the last non-empty value found:

```go
func extractTitle(row map[string]string, titleCols []string) string {
    last := ""
    for _, col := range titleCols {  // titleCols sorted ascending: Title 1, Title 2...
        if val := strings.TrimSpace(row[col]); val != "" {
            last = val
        }
    }
    if last == "" {
        return fmt.Sprintf("Untitled %s %s", row["Work Item Type"], row["ID"])
    }
    return last
}
```

### 5.7 Iteration Path Parsing

Iteration Path format varies in depth and uses `\\` as separator. The parser must identify the sprint name intelligently rather than by fixed segment position.

Real examples from production exports:
```
FinDash\\Archive\\FY26 Q3              → unscheduled (Archive)
FinDash\\Backlog\\Queue                → unscheduled (Backlog)
FinDash\\Delivery\\FY26 Q3\\FY26 Q3.1 → sprint = "FY26 Q3", sub-sprint = "FY26 Q3.1"
FinDash\\Team-Finance-Insights\\Sprint 2\\DataModel → sprint = "Sprint 2"
Finance Analytics Dashboard\\Sprint 1  → sprint = "Sprint 1"
```

Algorithm:

```go
func parseIterationPath(path string) (sprintName string, isScheduled bool) {
    segments := strings.Split(path, `\\`)

    // Unscheduled signals — treat as backlog
    unscheduledKeywords := []string{"backlog", "archive", "queue"}
    for _, seg := range segments {
        for _, kw := range unscheduledKeywords {
            if strings.EqualFold(strings.TrimSpace(seg), kw) {
                return "", false
            }
        }
    }

    // Find the first segment that looks like a sprint
    // Sprint patterns: "Sprint N", "FY2N QN", "Q\d", "\d{4} Q\d", "Iteration N"
    sprintPatterns := []*regexp.Regexp{
        regexp.MustCompile(`(?i)^sprint\s+\d+`),
        regexp.MustCompile(`(?i)^FY\d{2}\s+Q\d`),
        regexp.MustCompile(`(?i)^Q\d{1,2}$`),
        regexp.MustCompile(`(?i)^iteration\s+\d+`),
        regexp.MustCompile(`(?i)^\d{4}\s+Q\d`),
    }

    for _, seg := range segments {
        seg = strings.TrimSpace(seg)
        for _, pattern := range sprintPatterns {
            if pattern.MatchString(seg) {
                return seg, true
            }
        }
    }

    // No recognisable sprint segment found → unscheduled
    return "", false
}
```

Sprint records are created from unique sprint names encountered during parsing. Start and end dates are populated from the sprint configuration step.

### 5.8 Target Date Parsing

Target Date field must handle multiple formats including full datetime strings. Parse in this order:

```go
var dateFormats = []string{
    "2006-01-02",                    // ISO 8601
    "01/02/2006",                    // US MM/DD/YYYY
    "02/01/2006",                    // UK DD/MM/YYYY
    "2006/01/02",                    // ISO with slashes
    "01/02/2006 3:04:05 PM",         // US datetime with 12hr (real DevOps format)
    "01/02/2006 15:04:05",           // US datetime with 24hr
    "2006-01-02T15:04:05Z07:00",     // ISO 8601 with timezone
    "January 2, 2006",               // Long form
    "2-Jan-2006",                    // Abbreviated month
}

func parseTargetDate(raw string) (*time.Time, error) {
    raw = strings.TrimSpace(raw)
    if raw == "" {
        return nil, nil  // valid — not an error
    }
    for _, format := range dateFormats {
        if t, err := time.Parse(format, raw); err == nil {
            // Strip time component — store date only
            d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
            return &d, nil
        }
    }
    return nil, fmt.Errorf("unparseable date: %q", raw)
}
```

Unparseable non-empty value → log warning, treat as empty, add to date assignment list.

> **Ambiguity note:** `01/02/2026` is ambiguous between US (Jan 2) and UK (Feb 1) format. Default to US. During M0, confirm the regional format used by the internal pilot teams.

### 5.9 Assigned To Parsing

Two formats observed in real exports:

```
"First Last <email@domain.com>"      → common format
"Last, First <email@domain.com>"     → also common
"email@domain.com"                   → email only, no display name
""                                   → unassigned
```

```go
func parseAssignedTo(raw string) (displayName, email string) {
    raw = strings.TrimSpace(raw)
    if raw == "" {
        return "", ""
    }

    // Extract email from angle brackets if present
    emailMatch := regexp.MustCompile(`<([^>]+)>`).FindStringSubmatch(raw)
    if len(emailMatch) > 1 {
        email = strings.TrimSpace(emailMatch[1])
        // Display name is everything before the angle bracket
        displayName = strings.TrimSpace(strings.Split(raw, "<")[0])
        // Handle "Last, First" → store as-is; UI displays what DevOps provides
        displayName = strings.Trim(displayName, `"' `)
        return displayName, email
    }

    // No angle bracket — check if it looks like an email
    if strings.Contains(raw, "@") {
        return "", raw  // email only
    }

    // Plain name, no email
    return raw, ""
}
```

### 5.10 State Normalisation

Map known state variants to Maestro's four internal states:

| Raw DevOps Value | Maestro State | Display Color |
|---|---|---|
| `New` | `not_started` | Grey |
| `To Do` | `not_started` | Grey |
| `Proposed` | `not_started` | Grey |
| `Active` | `in_progress` | Blue |
| `In Progress` | `in_progress` | Blue |
| `Committed` | `in_progress` | Blue |
| `Open` | `in_progress` | Blue |
| `Blocked` | `blocked` | Amber |
| `On Hold` | `blocked` | Amber |
| `Resolved` | `done` | Green |
| `Closed` | `done` | Green |
| `Completed` | `done` | Green |
| `Done` | `done` | Green |

Unrecognised state → stored as raw string, mapped to `not_started` for color coding, flagged in import report.

### 5.11 Story Points Parsing

- Empty string → `story_points = nil` (valid; Sprint Load Index shows "—")
- Float string (`"5.0"`) → parse as float, round to nearest integer
- Integer string (`"5"`) → parse directly
- Non-numeric → log warning, set nil

### 5.12 Import Report

After parsing, before the PM confirms, Maestro shows:

```
PARSED SUCCESSFULLY
✓  Epics:                     12
✓  Features:                  48
✓  Stories:                   180
✓  Sprints detected:          8

WARNINGS — import will proceed, but some items need attention
⚠  Items missing Target Date:    94  → Date Assignment screen after import
⚠  Items missing sprint:         23  → Unscheduled; assign sprint in List view
⚠  Items missing Story Points:  180  → Sprint Load Index unavailable
⚠  Items with blank titles:       3  → Assigned "Untitled Story {ID}"
⚠  Orphaned items:                2  → Assigned to Unassigned bucket
⚠  Unknown states found:          1  → "Expedite" stored as-is, mapped to In Progress

OPTIONAL DATA
ℹ  Area Path detected            →  4 unique areas available as List view filter
ℹ  Assigned To detected          →  12 unique team members
ℹ  Date format detected          →  US (MM/DD/YYYY HH:MM:SS AM/PM)

ROWS SKIPPED
✗  Unknown Work Item Types:       0
✗  Missing required fields:       0
```

The PM reviews this before confirming. If any rows were skipped, the report shows exactly which IDs were affected and why.

---

## Section 6: Post-Import Date Assignment

If any items have `NULL` Target Date after import, Maestro shows a Date Assignment screen before the PM reaches the Gantt.

```
┌──────────────────────────────────────────────────────────────┐
│  94 items need dates                                         │
│  Set them now or come back later from the List view.         │
│                                                              │
│  Filter: [All types ▾]  [All sprints ▾]                     │
│                                                              │
│  ○ Govern Drilldown adapter      Sprint 1   [pick date ▾]   │
│  ○ Calibrate Sandbox revenue     Sprint 4   [pick date ▾]   │
│  ○ Deploy Aster scheduler        Backlog    [assign sprint]  │
│  ○ Localize Closebook revenue    Sprint 2   [pick date ▾]   │
│  ... 90 more                                                 │
│                                                              │
│  [Set all to sprint end dates]   [Skip — I'll do it later]  │
└──────────────────────────────────────────────────────────────┘
```

**"Set all to sprint end dates"** is a bulk action — assigns `committed_end_date` = sprint `end_date` for every undated item that has a sprint assignment. `date_source` = `inherited`. Items with no sprint remain undated.

**First PM-assigned date becomes `original_end_date`** and locks — same rule as imported dates. Subsequent changes move only `committed_end_date`.

Items with `date_source = inherited` are excluded from Deadline Hit Rate and Slip Analysis — inherited dates are not commitments.

---

## Section 7: Tiered Functionality Model

Rather than rejecting imports that are missing recommended columns, Maestro degrades gracefully and tells the PM exactly what each missing column costs them.

| Tier | Columns Present | What Works |
|---|---|---|
| **Tier 1 — Minimum viable** | `Parent`, `ID`, `Work Item Type` | Hierarchy renders; items visible in List view; no Gantt bars; no metrics |
| **Tier 2 — Gantt renders** | + `Iteration Path`, `Target Date` | Gantt bars appear; sprint placement works; no health metrics |
| **Tier 3 — Health metrics** | + `State`, `Story Points` | Deadline Hit Rate, Sprint Load Index, Slip Analysis all available |
| **Tier 4 — Full enrichment** | + `Assigned To`, `Area Path` | Owner display, area filtering, full team visibility |

At import, Maestro detects which tier the export reaches and shows a banner:

> **"Your export supports Tier 2 — Gantt view is available. To enable health metrics, add State and Story Points columns to your DevOps query and re-export."**

This removes the pressure to get the query perfect on first run and makes Maestro immediately useful even with partial data.

---

## Section 8: Common Export Problems and Fixes

| Problem | Symptom | Fix |
|---|---|---|
| `Parent` column missing | Import rejected — missing required column | Add Parent via Column Options and re-export |
| Wrong query type (flat list) | All items show as orphans — no hierarchy | Re-run as "Tree of work items" |
| Story Points all empty | Sprint Load Index shows "—" for all sprints | Expected for many teams — use Tier 3 guidance or accept metric unavailability |
| Target Date empty on most items | 50+ items in date assignment screen | Expected — use bulk "set to sprint end" or assign individually |
| `Blocked` state not recognised | Items shown as grey instead of amber | Now handled in state normalisation table |
| Iteration Path segments unrecognised | Items show as unscheduled | Check sprint naming pattern against supported formats; contact Maestro support to add pattern |
| Parent ID shows as `500588.0` | No effect — parser strips decimal automatically | No action needed |
| `Assigned To` in `Last, First` format | No effect — parser handles both formats | No action needed |
| Items in Archive or Backlog path | Items show as unscheduled — correct | These are genuinely unscheduled; assign sprint in List view if needed |

---

## Section 9: PRD Sections Superseded by This Addendum

| PRD Section | Change |
|---|---|
| §4.1 Core Entities | `Parent` field added to Feature and Story entities. `date_source` field added to Story. |
| §4.2 Import Schema Mapping | Entirely superseded. Hierarchy via `Parent` ID, not title columns. Full parser spec in Addendum §5. |
| §5.1 Import — Accepted Formats | CSV only. JSON removed from POC scope. |
| §5.2 Gantt View | Tiered functionality model applies — Gantt renders at Tier 2, not Tier 1. |
| §5.3 List View | All items visible at Tier 1 regardless of date/sprint availability. |
| §5.4 Health Dashboard | Metrics require Tier 3. Missing Story Points permanently disables Sprint Load Index for that import. |
| §9 Milestones — M0 | M0 now includes: validate parser against all three example CSVs; confirm sprint name patterns for internal teams; confirm regional date format; measure SQLite baseline binary size. |
