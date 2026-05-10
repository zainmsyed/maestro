# Story 011: SVAR data bridge + basic Gantt screen

**Status:** complete  
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** 2026-05-10  
**Completed:** 2026-05-10

---

## Goal
Build the production Gantt view using SVAR Svelte Gantt. Create a data bridge that maps Maestro Epic and Feature entities to SVAR's flat task array with correct hierarchy (Epic summary rows, Feature task rows), sprint-derived start dates, and committed end dates. Wire the Gantt view into the app's routing.

## Verification
Navigate to the Timeline (Gantt) view with real imported data. Verify all Epics appear as expandable summary rows, Features appear as bars under their parent Epic, the synthetic "Unassigned" Epic sits at the bottom, bars are positioned by sprint start dates, and end dates match `committed_end_date`. Expand and collapse an Epic row.

## Scope — files this story may touch
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/lib/svarBridge.ts` (new)
- `frontend/src/App.svelte` (routing)
- `frontend/src/lib/api.ts` (types if needed)
- `frontend/src/stores/view.ts` (view state)

## Out of scope — do not touch
- Status color coding (story-012)
- Today line or sprint boundary overlays (story-012)
- Drag interactions (story-013)
- Detail panel (story-014)
- List view, Health dashboard, Settings

## Dependencies
- story-004
- story-005
- story-010

---

## Checklist
- [x] Create `frontend/src/lib/svarBridge.ts`: two-pass builder that turns `GET /api/epics` response into SVAR task array
  - Pass 1: Epics → `type: "summary"`, `id = epic.id`, `parent = 0`
  - Pass 2: Features → `type: "task"`, `id = feature.id`, `parent = feature.epic_id`
  - Synthetic "Unassigned" Epic gets `is_synthetic` detail flag
- [x] Derive `start` from sprint start date fallback chain (sprint start → sprint end → import date)
  - Features with no valid start date are omitted from the task array
- [x] Derive `end` from `committed_end_date`
- [x] Derive `progress` from status (not started = 0, in progress = 50, done = 100)
- [x] Build `GanttView.svelte` shell: fetches epics, builds bridge, renders `<Gantt>` component
- [x] Wire Gantt route in `App.svelte` under Roadmap → Gantt tab
- [x] Handle loading and empty states
- [x] Verify expand/collapse works on Epic summary rows
- [x] Verify synthetic Unassigned Epic renders at bottom of row list

## Issues

## Completion Summary
Implemented the production SVAR Gantt route. Added `frontend/src/lib/svarBridge.ts` to map `GET /api/epics` into SVAR summary/task rows with epics as expandable summaries, features as child tasks, synthetic epics sorted last, sprint/import-derived starts, committed-date ends, and status-derived progress. Added `frontend/src/screens/GanttView.svelte` with loading, error, empty, summary, and Gantt states. Wired Roadmap → Gantt in `App.svelte`. Updated design context gaps from story-011 answers before component work. Verification performed: `npm run check` passes with 0 warnings, `npm test` passes 25 tests, and `npm run build` succeeds. Follow-up runtime fix: summary rows without children now render closed so SVAR does not recurse into null child arrays.
