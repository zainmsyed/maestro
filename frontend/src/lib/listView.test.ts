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
        stories: [
          {
            id: 'S-1',
            feature_id: 'F-1',
            title: 'Story One',
            description: '',
            status: 'Active',
            owner: 'Bob',
            sprint: 'Sprint 2',
            story_points: 5,
            original_end_date: '2026-05-02T00:00:00Z',
            committed_end_date: null,
            actual_end_date: null,
            date_source: 'imported',
            created_at: '',
            updated_at: '',
          },
        ],
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
  it('flattens epics, features, and stories with slip counts', () => {
    const rows = flattenRows(epics, audits);
    expect(rows).toHaveLength(5); // 2 epics + 2 features + 1 story
    const feature = rows.find((row) => row.id === 'F-1');
    expect(feature?.slipEvents).toBe(1);
    expect(feature?.health).toBe('amber');
    const story = rows.find((row) => row.id === 'S-1');
    expect(story?.type).toBe('story');
    expect(story?.epicTitle).toBe('Alpha Epic');
    expect(story?.featureTitle).toBe('Feature One');
    expect(story?.storyPoints).toBe(5);
    expect(story?.dateSource).toBe('imported');
  });

  it('filters by type', () => {
    const rows = flattenRows(epics, audits);
    const filters: FilterState = { epic: '', owner: '', sprint: '', status: '', type: 'story' };
    const filtered = filterRows(rows, filters);
    expect(filtered.every((r) => r.type === 'story')).toBe(true);
    expect(filtered).toHaveLength(1);
  });

  it('filters and groups rows by sprint', () => {
    const rows = flattenRows(epics, audits);
    const filters: FilterState = { epic: '', owner: 'Bob', sprint: 'Sprint 2', status: 'Active', type: '' };
    const filtered = filterRows(rows, filters);
    expect(filtered.map((row) => row.id)).toEqual(['F-1', 'S-1']);

    const groups = buildGroups(filtered, 'sprint', 'title', 'asc');
    expect(groups).toHaveLength(1);
    expect(groups[0].label).toBe('Sprint 2');
  });

  it('deduplicates filter labels case-insensitively and trims whitespace', () => {
    expect(uniqueLabels(['Alice', ' alice ', 'ALICE', 'Bob'])).toEqual(['Alice', 'Bob']);
  });

  it('groups by epic with three-level hierarchy', () => {
    const rows = flattenRows(epics, audits);
    const groups = buildGroups(rows, 'epic', 'title', 'asc');
    const alphaGroup = groups.find((g) => g.label === 'Alpha Epic');
    expect(alphaGroup?.rows[0].type).toBe('epic');
    expect(alphaGroup?.rows[1].type).toBe('feature');
    expect(alphaGroup?.rows[2].type).toBe('story');
    expect(alphaGroup?.rows[2].id).toBe('S-1');
  });

  it('shapes export rows from grouped view with date_source', () => {
    const rows = flattenRows(epics, audits);
    const groups = buildGroups(rows, 'epic', 'title', 'asc');
    const exportRows = buildExportRows(groups, 'epic');
    expect(exportRows[0]).toMatchObject({
      group_by: 'epic',
      group: 'Alpha Epic',
      title: 'Alpha Epic',
      type: 'epic',
    });
    const storyExport = exportRows.find((r) => r.title === 'Story One');
    expect(storyExport).toMatchObject({
      type: 'story',
      date_source: 'imported',
      story_points: 5,
    });
  });
});
