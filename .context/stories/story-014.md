# Story 014: SVAR detail panel

**Status:** not-started
**Type:** ui
**Created:** 2026-05-09
**Last accessed:** —
**Completed:** —

---

## Goal
Build a slide-out detail panel that opens when a PM clicks any Gantt bar (Epic or Feature). The panel shows read-only original date, editable committed date, status, owner, sprint, net slip days, and a chronological audit log of all date changes for that item.

## Verification
In the Gantt view, click a Feature bar. Verify a detail panel slides in from the right showing: title, type, owner, sprint, original end date (strikethrough), committed end date, status, net slip days, and an audit log table. Click the original date to confirm it is read-only. Change the committed date via the inline picker and confirm the bar updates and an audit entry is created. Close the panel by clicking the X or outside the panel.

## Scope — files this story may touch
- `frontend/src/components/DetailPanel.svelte` (new)
- `frontend/src/screens/GanttView.svelte`
- `frontend/src/lib/api.ts` (audit log fetch)
- `frontend/src/components/DatePicker.svelte` (reuse)

## Out of scope — do not touch
- Drag interactions (already in story-013)
- List view detail panel (List uses inline editing)
- Health dashboard
- Settings screens

## Dependencies
- story-011
- story-012

---

## Checklist
- [ ] Build `DetailPanel.svelte` as a fixed-position slide-out overlay
- [ ] Wire `onTaskClick` on the SVAR Gantt to open the panel with the clicked task's ID
- [ ] Fetch full entity detail (Epic or Feature) from existing `GET /api/epics/:id` or `GET /api/features/:id`
- [ ] Fetch audit log via `GET /api/audit?entity_id={id}`
- [ ] Display fields: title, type pill, owner, sprint name, original end date (read-only, strikethrough), committed end date (editable via `DatePicker`), status, net slip days, slip event count
- [ ] Display audit log table: changed_at, old_date → new_date, delta_days, reason
- [ ] Close panel via X button, Escape key, or click outside
- [ ] Animate slide-in/slide-out
- [ ] Handle loading and empty states for audit log
- [ ] Verify committed date edit calls the same PATCH endpoint as drag and updates the bar optimistically

## Issues

## Completion Summary
