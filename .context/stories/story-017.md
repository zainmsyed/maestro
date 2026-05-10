# Story 017: Settings screens

**Status:** not-started
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Build the Settings screen with sidebar navigation, project configuration, sprint management table with inline capacity editing, metrics thresholds, and import history with re-import and clear-data actions.

## Verification
Open Settings → Sprints and edit a sprint's capacity inline. Refresh the page and confirm the value persisted. Navigate to Settings → Metrics and change the scope creep threshold to 20%. Return to the Health dashboard and confirm the warning boundary shifted.

## Scope — files this story may touch
- `frontend/src/screens/SettingsView.svelte`
- `frontend/src/components/SettingsNav.svelte`
- `frontend/src/components/SprintTable.svelte`
- `frontend/src/components/MetricsThresholds.svelte`
- `frontend/src/components/ImportHistory.svelte`
- `frontend/src/lib/api.ts`
- `frontend/src/stores/settings.ts`

## Out of scope — do not touch
- Auth or user management
- Multi-project switching
- External integrations beyond re-import

## Dependencies
- story-004
- story-005
- story-015

---

## Checklist
- [ ] Create `SettingsView.svelte` with sidebar nav and page content area
- [ ] Project settings page: project name, team name inputs with save
- [ ] Sprints page: full sprint table with inline capacity editing
- [ ] Add "+ add sprint" button with manual creation form
- [ ] Snap-to-sprint toggle with persistence
- [ ] "Reset to import defaults" button with confirmation dialog
- [ ] Metrics thresholds page: sprint load warning, scope creep warning, deadline hit rate target
- [ ] Import history page: last import file, timestamp, item counts
- [ ] Re-import button that routes to onboarding
- [ ] "Clear all data" danger button with confirmation
- [ ] Ensure all changes call the appropriate PATCH/POST APIs

---

## Issues

---

## Completion Summary
