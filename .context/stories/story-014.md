# Story 014: Health dashboard with 6 metrics

**Status:** not-started
**Type:** ui
**Created:** 2026-05-06
**Last accessed:** 2026-05-06
**Completed:** —

---

## Goal
Build the Health dashboard screen that surfaces all 6 execution metrics with clear visualizations, per-sprint charts, and a dedicated "Recoveries" callout for schedule improvements.

## Verification
Open the Health dashboard with test data. Verify each metric matches hand-calculated expectations for the fixture dataset. Confirm the Recoveries callout shows features with negative net slip days. Change a filter (epic or date range) and verify metrics recalculate.

## Scope — files this story may touch
- `frontend/src/screens/HealthView.svelte`
- `frontend/src/components/MetricCard.svelte`
- `frontend/src/components/SprintLoadChart.svelte`
- `frontend/src/components/VelocityChart.svelte`
- `frontend/src/components/SlipDistribution.svelte`
- `frontend/src/components/RecoveriesCallout.svelte`
- `frontend/src/components/FilterBar.svelte`
- `frontend/src/lib/api.ts`

## Out of scope — do not touch
- Metrics threshold editing (Settings → Metrics comes in story-016)
- Gantt view
- List view

## Dependencies
- story-004
- story-005

---

## Checklist
- [ ] Create `HealthView.svelte` with grid layout for metric cards
- [ ] Deadline Hit Rate: percentage card with 4-sprint trend line
- [ ] Sprint Load Index: per-sprint bar chart with underloaded/healthy/overloaded bands
- [ ] Scope Creep Rate: percentage card with trend indicator
- [ ] Slip Analysis: Slip Event Rate percentage, Avg Net Slip Days, distribution breakdown (1× / 2× / 3+ slip events)
- [ ] Velocity Trend: rolling 4-sprint average line chart
- [ ] Orphaned Feature Rate: percentage card with reassignment shortcut
- [ ] Recoveries callout: count, total days recovered, top 3 features by days recovered
- [ ] Add Epic and date-range filters that apply to all metrics
- [ ] Handle missing capacity gracefully (show "—" with inline prompt)
- [ ] Verify all formulas match PRD §5.4 exactly

---

## Issues

---

## Completion Summary
