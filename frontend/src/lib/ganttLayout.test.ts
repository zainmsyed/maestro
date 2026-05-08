import { describe, expect, it } from 'vitest';
import type { EpicRecord, FeatureRecord, SprintDraft } from './api';
import { buildGanttLayout, visibleGanttRows } from './ganttLayout';

function feature(overrides: Partial<FeatureRecord>): FeatureRecord {
  return {
    id: 'F-1',
    epic_id: 'E-1',
    title: 'Feature one',
    description: '',
    status: 'New',
    owner: 'Iris',
    sprint: 'Sprint 1',
    story_points: null,
    original_end_date: null,
    committed_end_date: null,
    actual_end_date: null,
    date_source: 'imported',
    stories: [],
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

function epic(overrides: Partial<EpicRecord>): EpicRecord {
  return {
    id: 'E-1',
    title: 'Alpha Epic',
    description: '',
    status: 'In Progress',
    owner: 'Tanaka',
    sprint_start: 'Sprint 1',
    sprint_end: 'Sprint 2',
    original_end_date: null,
    committed_end_date: null,
    actual_end_date: null,
    is_synthetic: false,
    features: [],
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

function dateLabel(date: Date): string {
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${date.getFullYear()}-${month}-${day}`;
}

const sprintDrafts: SprintDraft[] = [
  { name: 'Sprint 1', start_date: '2026-01-01', length_days: 14, source: 'generated' },
  { name: 'Sprint 2', start_date: '2026-01-15', length_days: 14, source: 'generated' },
];

describe('buildGanttLayout', () => {
  it('orders synthetic unassigned epics last and expands feature rows under parents', () => {
    const layout = buildGanttLayout([
      epic({ id: 'epic-unassigned', title: 'Unassigned', is_synthetic: true, features: [feature({ id: 'F-U', title: 'Loose feature' })] }),
      epic({ id: 'E-1', title: 'Alpha Epic', features: [feature({ id: 'F-1', title: 'Planned feature' })] }),
    ], sprintDrafts, 'sprint', new Date('2026-01-08'));

    expect(layout.rows.map((row) => row.id)).toEqual(['E-1', 'F-1', 'epic-unassigned', 'F-U']);
    expect(layout.rows.find((row) => row.id === 'epic-unassigned')?.isSynthetic).toBe(true);
  });

  it('uses sprint names and draft dates to position bars and today line', () => {
    const layout = buildGanttLayout([
      epic({ features: [feature({ id: 'F-1', sprint: 'Sprint 2' })] }),
    ], sprintDrafts, 'sprint', new Date('2026-01-08'));

    const epicRow = layout.rows.find((row) => row.id === 'E-1');
    const featureRow = layout.rows.find((row) => row.id === 'F-1');

    expect(layout.sprints.map((sprint) => sprint.name)).toEqual(['Sprint 1', 'Sprint 2']);
    expect(epicRow?.bar.left).toBe(10);
    expect(epicRow?.bar.width).toBe((172 * 2) - 20);
    expect(featureRow?.bar.left).toBe(172 + 10);
    expect(layout.todayLeft).toBeGreaterThan(0);
    expect(layout.todayLeft).toBeLessThan(172);
  });

  it('infers sprint dates from item dates when sprint drafts are unavailable', () => {
    const layout = buildGanttLayout([
      epic({ sprint_start: 'FY26 Sprint 4', sprint_end: 'FY26 Sprint 4', committed_end_date: '2026-03-20' }),
    ], [], 'sprint', new Date('2026-03-01'));

    expect(layout.sprints[0].name).toBe('FY26 Sprint 4');
    expect(dateLabel(layout.sprints[0].end)).toBe('2026-03-20');
    expect(dateLabel(layout.sprints[0].start)).toBe('2026-03-07');
  });

  it('maps status tones for active, on-time completed, and missed work', () => {
    const layout = buildGanttLayout([
      epic({
        features: [
          feature({ id: 'active', status: 'In Progress', committed_end_date: '2026-01-20' }),
          feature({ id: 'done', status: 'Completed', committed_end_date: '2026-01-20', actual_end_date: '2026-01-19' }),
          feature({ id: 'missed', status: 'Completed', committed_end_date: '2026-01-20', actual_end_date: '2026-01-21' }),
          feature({ id: 'overdue', status: 'New', committed_end_date: '2026-01-02' }),
        ],
      }),
    ], sprintDrafts, 'sprint', new Date('2026-01-10'));

    const tones = new Map(layout.rows.map((row) => [row.id, row.bar.tone]));
    expect(tones.get('active')).toBe('in-progress');
    expect(tones.get('done')).toBe('completed');
    expect(tones.get('missed')).toBe('missed');
    expect(tones.get('overdue')).toBe('missed');
  });
});

describe('visibleGanttRows', () => {
  it('hides feature rows for collapsed epics', () => {
    const layout = buildGanttLayout([
      epic({ id: 'E-1', features: [feature({ id: 'F-1' })] }),
      epic({ id: 'E-2', features: [feature({ id: 'F-2', epic_id: 'E-2' })] }),
    ], sprintDrafts, 'sprint', new Date('2026-01-08'));

    expect(visibleGanttRows(layout.rows, new Set(['E-1'])).map((row) => row.id)).toEqual(['E-1', 'F-1', 'E-2']);
  });
});
