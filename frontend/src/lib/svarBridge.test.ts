import { describe, expect, it } from 'vitest';
import type { EpicRecord, FeatureRecord } from './api';
import { buildSvarTasks } from './svarBridge';

function feature(overrides: Partial<FeatureRecord> = {}): FeatureRecord {
  return {
    id: 'F-1',
    epic_id: 'E-1',
    title: 'Feature One',
    description: '',
    status: 'Active',
    owner: 'Owner',
    sprint: 'FY26 Q1',
    story_points: null,
    original_end_date: null,
    committed_end_date: '2026-05-20T00:00:00Z',
    actual_end_date: null,
    date_source: 'imported',
    stories: [],
    created_at: '2026-05-01T00:00:00Z',
    updated_at: '2026-05-01T00:00:00Z',
    ...overrides,
  };
}

function epic(overrides: Partial<EpicRecord> = {}): EpicRecord {
  return {
    id: 'E-1',
    title: 'Alpha Epic',
    description: '',
    status: 'New',
    owner: 'Epic Owner',
    sprint_start: '2026-05-06T00:00:00Z',
    sprint_end: '2026-05-13T00:00:00Z',
    original_end_date: null,
    committed_end_date: null,
    actual_end_date: null,
    is_synthetic: false,
    features: [feature()],
    created_at: '2026-05-01T00:00:00Z',
    updated_at: '2026-05-01T00:00:00Z',
    ...overrides,
  };
}

describe('buildSvarTasks', () => {
  it('builds epic summary rows before feature task rows with hierarchy intact', () => {
    const tasks = buildSvarTasks([epic()]);

    expect(tasks).toHaveLength(2);
    expect(tasks[0]).toMatchObject({ id: 'E-1', parent: 0, type: 'summary', source_type: 'epic' });
    expect(tasks[1]).toMatchObject({ id: 'F-1', parent: 'E-1', type: 'task', source_type: 'feature' });
  });

  it('sorts synthetic epics last and preserves the synthetic detail flag', () => {
    const tasks = buildSvarTasks([
      epic({ id: 'E-synthetic', title: 'Unassigned Epic', is_synthetic: true, features: [] }),
      epic({ id: 'E-real', title: 'Alpha Epic', features: [] }),
    ]);

    expect(tasks.map((task) => task.id)).toEqual(['E-real', 'E-synthetic']);
    expect(tasks[1]).toMatchObject({ is_synthetic: true, open: false });
  });

  it('uses sprint start fallback chain and committed end date for feature bars', () => {
    const tasks = buildSvarTasks([
      epic({
        sprint_start: '',
        sprint_end: '2026-05-10T00:00:00Z',
        features: [feature({ committed_end_date: '2026-05-22T00:00:00Z' })],
      }),
    ]);

    const task = tasks.find((item) => item.source_type === 'feature');
    expect(task?.start?.toISOString()).toBe('2026-05-10T00:00:00.000Z');
    expect(task?.end?.toISOString()).toBe('2026-05-22T00:00:00.000Z');
  });

  it('maps status to progress buckets and CSS status classes', () => {
    const tasks = buildSvarTasks([
      epic({
        features: [
          feature({ id: 'todo', status: 'New' }),
          feature({ id: 'active', status: 'Active' }),
          feature({ id: 'done', status: 'Completed' }),
        ],
      }),
    ]);

    expect(tasks.find((task) => task.id === 'todo')).toMatchObject({
      progress: 0,
      status_variant: 'not-started',
      css: expect.stringContaining('maestro-status-not-started'),
    });
    expect(tasks.find((task) => task.id === 'active')).toMatchObject({
      progress: 50,
      status_variant: 'in-progress',
      css: expect.stringContaining('maestro-status-in-progress'),
    });
    expect(tasks.find((task) => task.id === 'done')).toMatchObject({
      progress: 100,
      status_variant: 'done',
      css: expect.stringContaining('maestro-status-done'),
    });
  });

  it('marks completed work as missed when actual end is after committed end', () => {
    const tasks = buildSvarTasks([
      epic({
        features: [
          feature({
            id: 'late',
            status: 'Completed',
            committed_end_date: '2026-05-20T00:00:00Z',
            actual_end_date: '2026-05-22T00:00:00Z',
          }),
        ],
      }),
    ]);

    expect(tasks.find((task) => task.id === 'late')).toMatchObject({
      status_variant: 'missed',
      is_missed: true,
      css: expect.stringContaining('maestro-status-missed'),
    });
  });

  it('marks synthetic feature rows with muted CSS class', () => {
    const tasks = buildSvarTasks([
      epic({ id: 'synthetic', is_synthetic: true, features: [feature({ id: 'orphan', epic_id: 'synthetic' })] }),
    ]);

    expect(tasks.find((task) => task.id === 'orphan')).toMatchObject({
      is_synthetic: true,
      css: expect.stringContaining('maestro-synthetic'),
    });
  });
});
