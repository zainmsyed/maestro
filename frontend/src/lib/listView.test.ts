import { describe, expect, it } from 'vitest';
import type { AuditRecord, EpicRecord } from './api';
import {
  buildExportRows,
  buildGroups,
  filterRows,
  flattenRows,
  uniqueLabels,
  type FilterState,
} from './listView';

const epics: EpicRecord[] = [
  {
    id: 'E-1',
    title: 'Alpha Epic',
    description: '',
    status: 'Active',
    owner: 'Alice',
    sprint_start: '',
    sprint_end: 'Sprint 1',
    original_end_date: '2026-05-01T00:00:00Z',
    committed_end_date: '2026-05-03T00:00:00Z',
    actual_end_date: null,
    is_synthetic: false,
    created_at: '',
    updated_at: '',
    features: [
      {
        id: 'F-1',
        epic_id: 'E-1',
        title: 'Feature One',
        description: '',
        status: 'Active',
        owner: 'Bob',
        sprint: 'Sprint 2',
        story_points: null,
        original_end_date: '2026-05-02T00:00:00Z',
        committed_end_date: '2026-05-04T00:00:00Z',
        actual_end_date: null,
        date_source: 'imported',
        stories: [],
        created_at: '',
        updated_at: '',
      },
    ],
  },
  {
    id: 'epic-unassigned',
    title: 'Unassigned Epic',
    description: '',
    status: 'Imported',
    owner: '  alice  ',
    sprint_start: '',
    sprint_end: '',
    original_end_date: null,
    committed_end_date: null,
    actual_end_date: null,
    is_synthetic: true,
    created_at: '',
    updated_at: '',
    features: [
      {
        id: 'feature-unassigned',
        epic_id: 'epic-unassigned',
        title: 'Unassigned Feature',
        description: '',
        status: 'Active',
        owner: 'alice',
        sprint: '',
        story_points: null,
        original_end_date: null,
        committed_end_date: null,
        actual_end_date: null,
        date_source: 'imported',
        stories: [],
        created_at: '',
        updated_at: '',
      },
    ],
  },
];

const audits: AuditRecord[] = [
  {
    id: 1,
    entity_type: 'feature',
    entity_id: 'F-1',
    changed_by: 'pm',
    old_date: '2026-05-03T00:00:00Z',
    new_date: '2026-05-04T00:00:00Z',
    delta_days: 1,
    reason: 'Slip',
    changed_at: '2026-05-01T00:00:00Z',
  },
];

describe('listView helpers', () => {
  it('flattens epics and features with slip counts', () => {
    const rows = flattenRows(epics, audits);
    expect(rows).toHaveLength(4);
    const feature = rows.find((row) => row.id === 'F-1');
    expect(feature?.slipEvents).toBe(1);
    expect(feature?.health).toBe('amber');
  });

  it('filters and groups rows by sprint', () => {
    const rows = flattenRows(epics, audits);
    const filters: FilterState = { epic: '', owner: 'Bob', sprint: 'Sprint 2', status: 'Active' };
    const filtered = filterRows(rows, filters);
    expect(filtered.map((row) => row.id)).toEqual(['F-1']);

    const groups = buildGroups(filtered, 'sprint', 'title', 'asc');
    expect(groups).toHaveLength(1);
    expect(groups[0].label).toBe('Sprint 2');
    expect(groups[0].rows[0].id).toBe('F-1');
  });

  it('deduplicates filter labels case-insensitively and trims whitespace', () => {
    expect(uniqueLabels(['Alice', ' alice ', 'ALICE', 'Bob'])).toEqual(['Alice', 'Bob']);
  });


  it('shapes export rows from grouped view', () => {
    const rows = flattenRows(epics, audits);
    const groups = buildGroups(rows, 'epic', 'title', 'asc');
    const exportRows = buildExportRows(groups, 'epic');
    expect(exportRows[0]).toMatchObject({
      group_by: 'epic',
      group: 'Alpha Epic',
      title: 'Alpha Epic',
      type: 'epic',
    });
  });
});
