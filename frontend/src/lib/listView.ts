import type { AuditRecord, EpicRecord } from './api';

export type GroupBy = 'epic' | 'sprint' | 'owner';
export type SortKey = 'title' | 'type' | 'owner' | 'sprint' | 'original' | 'committed' | 'actual' | 'slip' | 'status' | 'health';
export type HealthTone = 'green' | 'amber' | 'red' | 'blue';

export type ListRow = {
  id: string;
  type: 'epic' | 'feature' | 'story';
  title: string;
  owner: string;
  sprint: string;
  originalDate: string;
  committedDate: string;
  actualDate: string;
  slipEvents: number;
  status: string;
  health: HealthTone;
  healthLabel: string;
  epicId: string | null;
  epicTitle: string;
  featureId: string | null;
  featureTitle: string;
  isSynthetic: boolean;
  dateSource: string;
  storyPoints: number | null;
};

export type ListGroup = { key: string; label: string; rows: ListRow[] };

export type FilterState = {
  epic: string;
  owner: string;
  sprint: string;
  status: string;
  type: string;
};

function normalizeDate(value: string | null | undefined): string {
  return value ? value.slice(0, 10) : '';
}

function deriveHealth(row: {
  status: string;
  originalDate: string;
  committedDate: string;
  actualDate: string;
  slipEvents: number;
}): { tone: HealthTone; label: string } {
  const status = row.status.toLowerCase();
  if (row.actualDate && row.committedDate && row.actualDate > row.committedDate) {
    return { tone: 'red', label: 'Late actual completion' };
  }
  if (row.committedDate && row.originalDate && row.committedDate > row.originalDate) {
    return { tone: 'amber', label: 'Committed date slipped' };
  }
  if (row.slipEvents > 0) {
    return { tone: 'amber', label: 'Has slip events' };
  }
  if (['done', 'closed', 'completed', 'resolved'].includes(status)) {
    return { tone: 'green', label: 'Completed without active risk' };
  }
  return { tone: 'blue', label: 'In progress / neutral' };
}

export function flattenRows(epicList: EpicRecord[], auditList: AuditRecord[]): ListRow[] {
  const slipCounts = new Map<string, number>();
  for (const audit of auditList) {
    const key = `${audit.entity_type}:${audit.entity_id}`;
    slipCounts.set(key, (slipCounts.get(key) ?? 0) + 1);
  }

  return epicList.flatMap((epic) => {
    const featureRows = epic.features.flatMap((feature) => {
      const storyRows = (feature.stories ?? []).map((story) => {
        const originalDate = normalizeDate(story.original_end_date);
        const committedDate = normalizeDate(story.committed_end_date);
        const actualDate = normalizeDate(story.actual_end_date);
        const slipEvents = slipCounts.get(`story:${story.id}`) ?? 0;
        const health = deriveHealth({
          status: story.status,
          originalDate,
          committedDate,
          actualDate,
          slipEvents,
        });
        return {
          id: story.id,
          type: 'story' as const,
          title: story.title,
          owner: story.owner,
          sprint: story.sprint,
          originalDate,
          committedDate,
          actualDate,
          slipEvents,
          status: story.status,
          health: health.tone,
          healthLabel: health.label,
          epicId: epic.id,
          epicTitle: epic.title,
          featureId: feature.id,
          featureTitle: feature.title,
          isSynthetic: epic.is_synthetic,
          dateSource: story.date_source,
          storyPoints: story.story_points,
        };
      });

      const originalDate = normalizeDate(feature.original_end_date);
      const committedDate = normalizeDate(feature.committed_end_date);
      const actualDate = normalizeDate(feature.actual_end_date);
      const slipEvents = slipCounts.get(`feature:${feature.id}`) ?? 0;
      const health = deriveHealth({
        status: feature.status,
        originalDate,
        committedDate,
        actualDate,
        slipEvents,
      });
      const featureRow: ListRow = {
        id: feature.id,
        type: 'feature' as const,
        title: feature.title,
        owner: feature.owner,
        sprint: feature.sprint,
        originalDate,
        committedDate,
        actualDate,
        slipEvents,
        status: feature.status,
        health: health.tone,
        healthLabel: health.label,
        epicId: epic.id,
        epicTitle: epic.title,
        featureId: null,
        featureTitle: '',
        isSynthetic: epic.is_synthetic,
        dateSource: feature.date_source,
        storyPoints: feature.story_points,
      };

      return [featureRow, ...storyRows];
    });

    const originalDate = normalizeDate(epic.original_end_date);
    const committedDate = normalizeDate(epic.committed_end_date);
    const actualDate = normalizeDate(epic.actual_end_date);
    const slipEvents = featureRows.reduce((sum, row) => sum + row.slipEvents, 0);
    const health = deriveHealth({
      status: epic.status,
      originalDate,
      committedDate,
      actualDate,
      slipEvents,
    });

    const epicRow: ListRow = {
      id: epic.id,
      type: 'epic' as const,
      title: epic.title,
      owner: epic.owner,
      sprint: epic.sprint_end,
      originalDate,
      committedDate,
      actualDate,
      slipEvents,
      status: epic.status,
      health: health.tone,
      healthLabel: health.label,
      epicId: epic.id,
      epicTitle: epic.title,
      featureId: null,
      featureTitle: '',
      isSynthetic: epic.is_synthetic,
      dateSource: '',
      storyPoints: null,
    };

    return [epicRow, ...featureRows];
  });
}

function compareValues(a: string | number, b: string | number): number {
  if (typeof a === 'number' && typeof b === 'number') return a - b;
  return String(a).localeCompare(String(b), undefined, { sensitivity: 'base' });
}

function normalizeOptionLabel(value: string | null | undefined, fallback = ''): string {
  const normalized = (value ?? fallback)
    .trim()
    .replace(/\s+/g, ' ');
  return normalized || fallback;
}

export function uniqueLabels(values: Array<string | null | undefined>, fallback = ''): string[] {
  const seen = new Set<string>();
  const labels: string[] = [];
  for (const value of values) {
    const label = normalizeOptionLabel(value, fallback);
    const key = label.toLocaleLowerCase();
    if (!seen.has(key)) {
      seen.add(key);
      labels.push(label);
    }
  }
  return labels.sort((a, b) => a.localeCompare(b, undefined, { sensitivity: 'base' }));
}

const sortValueByKey: Record<SortKey, (row: ListRow) => string | number> = {
  title: (row) => row.title,
  type: (row) => row.type,
  owner: (row) => row.owner || 'zzzz',
  sprint: (row) => row.sprint || 'zzzz',
  original: (row) => row.originalDate || '9999-99-99',
  committed: (row) => row.committedDate || '9999-99-99',
  actual: (row) => row.actualDate || '9999-99-99',
  slip: (row) => row.slipEvents,
  status: (row) => row.status || 'zzzz',
  health: (row) => row.health,
};

function rowSortValue(row: ListRow, key: SortKey): string | number {
  return sortValueByKey[key](row);
}

const filterPredicates: Record<keyof FilterState, (row: ListRow, value: string) => boolean> = {
  epic: (row, value) => row.epicTitle === value,
  owner: (row, value) => row.owner === value,
  sprint: (row, value) => (row.sprint || 'Unassigned') === value,
  status: (row, value) => row.status === value,
  type: (row, value) => row.type === value,
};

export function filterRows(rows: ListRow[], activeFilters: FilterState): ListRow[] {
  const activeEntries = Object.entries(activeFilters).filter(([, value]) => Boolean(value)) as Array<
    [keyof FilterState, string]
  >;
  if (activeEntries.length === 0) return rows;
  return rows.filter((row) => activeEntries.every(([key, value]) => filterPredicates[key](row, value)));
}

function sortRows(rows: ListRow[], activeSortKey: SortKey, activeSortDirection: 'asc' | 'desc'): ListRow[] {
  return [...rows].sort((left, right) => {
    const result = compareValues(rowSortValue(left, activeSortKey), rowSortValue(right, activeSortKey));
    return activeSortDirection === 'asc' ? result : -result;
  });
}

function groupByEpic(rows: ListRow[], activeSortKey: SortKey, activeSortDirection: 'asc' | 'desc'): ListGroup[] {
  const epicMap = new Map<string, ListGroup>();
  for (const row of rows) {
    const groupKey = row.epicId ?? 'none';
    const groupLabel = row.epicTitle || 'Unassigned';
    if (!epicMap.has(groupKey)) {
      epicMap.set(groupKey, { key: groupKey, label: groupLabel, rows: [] });
    }
    epicMap.get(groupKey)?.rows.push(row);
  }

  return [...epicMap.values()].map((group) => {
    const epicRow = group.rows.find((row) => row.type === 'epic');
    const featureRows = group.rows.filter((row) => row.type === 'feature');
    const storyRows = group.rows.filter((row) => row.type === 'story');

    const sortedFeatures = sortRows(featureRows, activeSortKey, activeSortDirection);
    const orderedRows: ListRow[] = [];

    if (epicRow) {
      orderedRows.push(epicRow);
    }

    for (const feature of sortedFeatures) {
      orderedRows.push(feature);
      const children = storyRows.filter((s) => s.featureId === feature.id);
      orderedRows.push(...sortRows(children, activeSortKey, activeSortDirection));
    }

    return { ...group, rows: orderedRows };
  });
}

function groupByField(rows: ListRow[], activeGroupBy: Exclude<GroupBy, 'epic'>): ListGroup[] {
  const groups = new Map<string, ListGroup>();
  for (const row of rows) {
    const groupLabel = activeGroupBy === 'sprint'
      ? row.sprint || 'Unassigned'
      : row.owner || 'Unassigned';
    const key = `${activeGroupBy}:${groupLabel}`;
    if (!groups.has(key)) {
      groups.set(key, { key, label: groupLabel, rows: [] });
    }
    groups.get(key)?.rows.push(row);
  }
  return [...groups.values()];
}

export function buildGroups(
  rows: ListRow[],
  activeGroupBy: GroupBy,
  activeSortKey: SortKey,
  activeSortDirection: 'asc' | 'desc',
): ListGroup[] {
  const sortedRows = sortRows(rows, activeSortKey, activeSortDirection);
  if (activeGroupBy === 'epic') {
    return groupByEpic(sortedRows, activeSortKey, activeSortDirection);
  }
  return groupByField(sortedRows, activeGroupBy);
}

export function buildExportRows(groups: ListGroup[], groupBy: GroupBy): Array<Record<string, string | number>> {
  return groups.flatMap((group) =>
    group.rows.map((row) => ({
      group_by: groupBy,
      group: group.label,
      title: row.title,
      type: row.type,
      epic: row.epicTitle,
      owner: row.owner || '',
      sprint: row.sprint || '',
      original_date: row.originalDate || '',
      committed_date: row.committedDate || '',
      actual_date: row.actualDate || '',
      story_points: row.storyPoints ?? '',
      date_source: row.dateSource,
      slip_events: row.slipEvents,
      status: row.status || '',
      health: row.healthLabel,
    })),
  );
}
