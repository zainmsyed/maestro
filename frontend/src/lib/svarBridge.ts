import type { ITask } from '@svar-ui/svelte-gantt';
import type { EpicRecord, FeatureRecord } from './api';

export type SvarStatusVariant = 'not-started' | 'in-progress' | 'done' | 'missed';

export type SvarTask = ITask & {
  id: string;
  parent: string | number;
  text: string;
  status?: string;
  owner?: string;
  sprint?: string;
  is_synthetic?: boolean;
  source_type: 'epic' | 'feature';
  maestro_id: string;
  missing_committed_end?: boolean;
  status_variant: SvarStatusVariant;
  is_missed?: boolean;
  css: string;
};

const ONE_DAY_MS = 24 * 60 * 60 * 1000;

function parseDate(value: string | null | undefined): Date | null {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;
  return date;
}

function addDays(date: Date, days: number): Date {
  return new Date(date.getTime() + days * ONE_DAY_MS);
}

function maxDate(dates: Array<Date | null | undefined>): Date | null {
  const valid = dates.filter((date): date is Date => date instanceof Date && !Number.isNaN(date.getTime()));
  if (valid.length === 0) return null;
  return new Date(Math.max(...valid.map((date) => date.getTime())));
}

function minDate(dates: Array<Date | null | undefined>): Date | null {
  const valid = dates.filter((date): date is Date => date instanceof Date && !Number.isNaN(date.getTime()));
  if (valid.length === 0) return null;
  return new Date(Math.min(...valid.map((date) => date.getTime())));
}

function featureStart(feature: FeatureRecord, epic: EpicRecord): Date | null {
  return parseDate(epic.sprint_start) ?? parseDate(epic.sprint_end) ?? parseDate(feature.created_at);
}

function epicStart(epic: EpicRecord, children: SvarTask[] = []): Date | null {
  return parseDate(epic.sprint_start) ?? parseDate(epic.sprint_end) ?? minDate(children.map((task) => task.start)) ?? parseDate(epic.created_at);
}

function normalizedStatus(status: string): string {
  return status.trim().toLowerCase().replace(/[_-]+/g, ' ');
}

function isDoneStatus(status: string): boolean {
  return ['done', 'complete', 'completed', 'closed', 'resolved'].includes(normalizedStatus(status));
}

function isInProgressStatus(status: string): boolean {
  return ['in progress', 'active', 'committed', 'doing'].includes(normalizedStatus(status));
}

function isMissed(committedEndDate: string | null | undefined, actualEndDate: string | null | undefined): boolean {
  const committed = parseDate(committedEndDate);
  const actual = parseDate(actualEndDate);
  return Boolean(committed && actual && actual.getTime() > committed.getTime());
}

function statusVariant(status: string, missed: boolean): SvarStatusVariant {
  if (isDoneStatus(status) && missed) return 'missed';
  if (isDoneStatus(status)) return 'done';
  if (isInProgressStatus(status)) return 'in-progress';
  return 'not-started';
}

function progressFromStatus(status: string): number {
  if (isDoneStatus(status)) return 100;
  if (isInProgressStatus(status)) return 50;
  return 0;
}

function taskCss(task: Pick<SvarTask, 'status_variant' | 'is_synthetic' | 'source_type'>): string {
  return [
    'maestro-gantt-bar',
    `maestro-status-${task.status_variant}`,
    `maestro-${task.source_type}`,
    task.is_synthetic ? 'maestro-synthetic' : '',
  ]
    .filter(Boolean)
    .join(' ');
}

function buildFeatureTask(feature: FeatureRecord, epic: EpicRecord): SvarTask | null {
  const start = featureStart(feature, epic);
  if (!start) return null;

  const committedEnd = parseDate(feature.committed_end_date);
  const end = committedEnd ?? addDays(start, 1);
  const missed = isMissed(feature.committed_end_date, feature.actual_end_date);
  const variant = statusVariant(feature.status, missed);
  const task: SvarTask = {
    id: feature.id,
    parent: feature.epic_id ?? epic.id,
    text: feature.title,
    start,
    end,
    progress: progressFromStatus(feature.status),
    type: 'task',
    status: feature.status,
    owner: feature.owner,
    sprint: feature.sprint,
    is_synthetic: epic.is_synthetic,
    source_type: 'feature',
    maestro_id: feature.id,
    missing_committed_end: !committedEnd,
    status_variant: variant,
    is_missed: missed,
    css: '',
  };
  task.css = taskCss(task);
  return task;
}

function buildEpicTask(epic: EpicRecord, children: SvarTask[]): SvarTask {
  const start = epicStart(epic, children) ?? new Date();
  const committedEnd = parseDate(epic.committed_end_date);
  const childEnd = maxDate(children.map((task) => task.end));
  const end = committedEnd ?? childEnd ?? addDays(start, 1);
  const childMissed = children.some((task) => task.is_missed);
  const missed = isMissed(epic.committed_end_date, epic.actual_end_date) || childMissed;
  const variant = statusVariant(epic.status, missed);
  const task: SvarTask = {
    id: epic.id,
    parent: 0,
    text: epic.title,
    start,
    end,
    progress: progressFromStatus(epic.status),
    type: 'summary',
    open: children.length > 0,
    status: epic.status,
    owner: epic.owner,
    is_synthetic: epic.is_synthetic,
    source_type: 'epic',
    maestro_id: epic.id,
    missing_committed_end: !committedEnd,
    status_variant: variant,
    is_missed: missed,
    css: '',
  };
  task.css = taskCss(task);
  return task;
}

function sortEpics(epics: EpicRecord[]): EpicRecord[] {
  return [...epics].sort((a, b) => {
    if (a.is_synthetic !== b.is_synthetic) return a.is_synthetic ? 1 : -1;
    return a.title.localeCompare(b.title, undefined, { sensitivity: 'base' });
  });
}

function buildFeatureTaskMap(epics: EpicRecord[]): Map<string, SvarTask[]> {
  return new Map(
    epics.map((epic) => [
      epic.id,
      epic.features
        .map((feature) => buildFeatureTask(feature, epic))
        .filter((task): task is SvarTask => task !== null),
    ]),
  );
}

function buildEpicTasks(epics: EpicRecord[], childrenByEpic: Map<string, SvarTask[]>): SvarTask[] {
  return epics.map((epic) => buildEpicTask(epic, childrenByEpic.get(epic.id) ?? []));
}

function flattenFeatureTasks(epics: EpicRecord[], childrenByEpic: Map<string, SvarTask[]>): SvarTask[] {
  return epics.flatMap((epic) => childrenByEpic.get(epic.id) ?? []);
}

export function buildSvarTasks(epics: EpicRecord[]): SvarTask[] {
  const orderedEpics = sortEpics(epics);
  const featureTasksByEpic = buildFeatureTaskMap(orderedEpics);

  return [
    // Pass 1: epics become summary rows.
    ...buildEpicTasks(orderedEpics, featureTasksByEpic),
    // Pass 2: features become task rows under their parent epic.
    ...flattenFeatureTasks(orderedEpics, featureTasksByEpic),
  ];
}
