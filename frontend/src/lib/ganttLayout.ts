import type { EpicRecord, FeatureRecord, SprintDraft } from './api';

export type GanttZoomLevel = 'sprint' | 'month' | 'quarter';
export type GanttStatusTone = 'not-started' | 'in-progress' | 'completed' | 'missed';
export type GanttRowType = 'epic' | 'feature';

export type GanttSprint = {
  id: string;
  name: string;
  start: Date;
  end: Date;
  left: number;
  width: number;
};

export type GanttBarLayout = {
  left: number;
  width: number;
  tone: GanttStatusTone;
  label: string;
};

export type GanttRow = {
  id: string;
  parentId: string | null;
  type: GanttRowType;
  title: string;
  status: string;
  owner: string;
  depth: number;
  isSynthetic: boolean;
  featureCount: number;
  bar: GanttBarLayout;
};

export type GanttLayout = {
  rows: GanttRow[];
  sprints: GanttSprint[];
  width: number;
  columnWidth: number;
  rowHeight: number;
  todayLeft: number | null;
};

type DateCarrier = {
  actual_end_date?: string | null;
  committed_end_date?: string | null;
  original_end_date?: string | null;
  created_at?: string | null;
};

type BarSource = DateCarrier & {
  title: string;
  status: string;
  sprint?: string;
  sprint_start?: string;
  sprint_end?: string;
};

type SprintDates = { start: Date; end: Date };

const ROW_HEIGHT = 48;
const DEFAULT_SPRINT_LENGTH_DAYS = 14;
const DAY_MS = 86_400_000;
const DONE_STATUS_PATTERN = /(done|complete|completed|closed|resolved)/;
const ACTIVE_STATUS_PATTERN = /(active|progress|doing|started|committed)/;
const ZOOM_COLUMN_WIDTH: Record<GanttZoomLevel, number> = {
  sprint: 172,
  month: 128,
  quarter: 96,
};

function parseDate(value: string | null | undefined): Date | null {
  if (!value) return null;
  const dateOnly = value.match(/^(\d{4})-(\d{2})-(\d{2})/);
  if (dateOnly) {
    return new Date(Number(dateOnly[1]), Number(dateOnly[2]) - 1, Number(dateOnly[3]));
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return null;
  return new Date(parsed.getFullYear(), parsed.getMonth(), parsed.getDate());
}

function addDays(date: Date, days: number): Date {
  const next = new Date(date);
  next.setDate(next.getDate() + days);
  return next;
}

function naturalCompare(left: string, right: string): number {
  return left.localeCompare(right, undefined, { numeric: true, sensitivity: 'base' });
}

function sprintLabel(value: string | null | undefined): string | null {
  const trimmed = value?.trim();
  return trimmed || null;
}

function uniqueSorted(values: Array<string | null | undefined>): string[] {
  return [...new Set(values.map(sprintLabel).filter((value): value is string => Boolean(value)))].sort(naturalCompare);
}

function parsedDates(values: Array<string | null | undefined>): Date[] {
  return values.map(parseDate).filter((date): date is Date => Boolean(date));
}

function targetDates(item: DateCarrier): Date[] {
  return parsedDates([item.actual_end_date, item.committed_end_date, item.original_end_date]);
}

function candidateDates(item: DateCarrier): Date[] {
  return [...targetDates(item), ...parsedDates([item.created_at])];
}

function earliestDate(dates: Date[]): Date | null {
  return dates.length === 0 ? null : dates.reduce((earliest, date) => (date < earliest ? date : earliest));
}

function latestDate(dates: Date[]): Date | null {
  return dates.length === 0 ? null : dates.reduce((latest, date) => (date > latest ? date : latest));
}

function featureSprintNames(feature: FeatureRecord): Array<string | null | undefined> {
  return [feature.sprint, ...(feature.stories ?? []).map((story) => story.sprint)];
}

function epicSprintNames(epic: EpicRecord): Array<string | null | undefined> {
  return [epic.sprint_start, epic.sprint_end, ...(epic.features ?? []).flatMap(featureSprintNames)];
}

function collectSprintNames(epics: EpicRecord[], drafts: SprintDraft[]): string[] {
  return uniqueSorted([...drafts.map((draft) => draft.name), ...epics.flatMap(epicSprintNames)]);
}

function appendDatesBySprint(dates: Map<string, Date[]>, sprint: string | null | undefined, item: DateCarrier) {
  const name = sprintLabel(sprint);
  if (!name) return;
  dates.set(name, [...(dates.get(name) ?? []), ...targetDates(item)]);
}

function appendFeatureDates(dates: Map<string, Date[]>, feature: FeatureRecord) {
  appendDatesBySprint(dates, feature.sprint, feature);
  for (const story of feature.stories ?? []) appendDatesBySprint(dates, story.sprint, story);
}

function collectDatesBySprint(epics: EpicRecord[]): Map<string, Date[]> {
  const dates = new Map<string, Date[]>();
  for (const epic of epics) {
    appendDatesBySprint(dates, epic.sprint_start, epic);
    appendDatesBySprint(dates, epic.sprint_end, epic);
    for (const feature of epic.features ?? []) appendFeatureDates(dates, feature);
  }
  return dates;
}

function draftSprintDates(draft: SprintDraft | undefined): SprintDates | null {
  const start = parseDate(draft?.start_date);
  if (!start) return null;
  const lengthDays = Math.max(1, draft?.length_days ?? DEFAULT_SPRINT_LENGTH_DAYS);
  return { start, end: addDays(start, lengthDays - 1) };
}

function inferredSprintDates(associatedDates: Date[] | undefined): SprintDates | null {
  const end = earliestDate(associatedDates ?? []);
  return end ? { start: addDays(end, -(DEFAULT_SPRINT_LENGTH_DAYS - 1)), end } : null;
}

function fallbackSprintDates(nextFallbackStart: Date): SprintDates {
  return { start: nextFallbackStart, end: addDays(nextFallbackStart, DEFAULT_SPRINT_LENGTH_DAYS - 1) };
}

function sprintDatesForName(
  name: string,
  drafts: Map<string, SprintDraft>,
  datesBySprint: Map<string, Date[]>,
  nextFallbackStart: Date,
): SprintDates {
  return draftSprintDates(drafts.get(name)) ?? inferredSprintDates(datesBySprint.get(name)) ?? fallbackSprintDates(nextFallbackStart);
}

function buildSprints(epics: EpicRecord[], drafts: SprintDraft[], importDate: Date, columnWidth: number): GanttSprint[] {
  const draftByName = new Map(drafts.map((draft) => [draft.name, draft]));
  const datesBySprint = collectDatesBySprint(epics);
  let nextFallbackStart = new Date(importDate.getFullYear(), importDate.getMonth(), importDate.getDate());

  return collectSprintNames(epics, drafts)
    .map((name) => {
      const dates = sprintDatesForName(name, draftByName, datesBySprint, nextFallbackStart);
      nextFallbackStart = addDays(dates.end, 1);
      return { id: name, name, ...dates, left: 0, width: columnWidth };
    })
    .sort((left, right) => left.start.getTime() - right.start.getTime() || naturalCompare(left.name, right.name))
    .map((sprint, index) => ({ ...sprint, left: index * columnWidth }));
}

function sprintIndexByName(sprints: GanttSprint[]): Map<string, number> {
  return new Map(sprints.map((sprint, index) => [sprint.name, index]));
}

function boundedDateIndex(sprints: GanttSprint[], date: Date): number | null {
  if (date < sprints[0].start) return 0;
  if (date > sprints[sprints.length - 1].end) return sprints.length - 1;
  return null;
}

function sprintIndexForDate(sprints: GanttSprint[], date: Date | null): number | null {
  if (!date || sprints.length === 0) return null;
  const exactIndex = sprints.findIndex((sprint) => date >= sprint.start && date <= sprint.end);
  return exactIndex >= 0 ? exactIndex : boundedDateIndex(sprints, date);
}

function sprintIndexForName(nameToIndex: Map<string, number>, value: string | null | undefined): number | null {
  const name = sprintLabel(value);
  return name ? nameToIndex.get(name) ?? null : null;
}

function firstNumber(...values: Array<number | null>): number | null {
  return values.find((value): value is number => typeof value === 'number') ?? null;
}

function isCompletedStatus(status: string): boolean {
  return DONE_STATUS_PATTERN.test(status.toLowerCase());
}

function isActiveStatus(status: string): boolean {
  return ACTIVE_STATUS_PATTERN.test(status.toLowerCase());
}

function completedTone(committed: Date | null, actual: Date | null): GanttStatusTone {
  return actual && committed && actual > committed ? 'missed' : 'completed';
}

function statusTone(status: string, committedEndDate: string | null, actualEndDate: string | null, today: Date): GanttStatusTone {
  const committed = parseDate(committedEndDate);
  const actual = parseDate(actualEndDate);
  if (isCompletedStatus(status)) return completedTone(committed, actual);
  if (committed && committed < today) return 'missed';
  return isActiveStatus(status) ? 'in-progress' : 'not-started';
}

function barIndexes(item: BarSource, sprints: GanttSprint[], nameToIndex: Map<string, number>) {
  const dates = candidateDates(item);
  const namedStart = sprintIndexForName(nameToIndex, item.sprint_start ?? item.sprint);
  const namedEnd = sprintIndexForName(nameToIndex, item.sprint_end ?? item.sprint);
  const dateStart = sprintIndexForDate(sprints, earliestDate(dates));
  const dateEnd = sprintIndexForDate(sprints, latestDate(dates));
  const startIndex = firstNumber(namedStart, namedEnd, dateStart) ?? 0;
  const endIndex = Math.max(startIndex, firstNumber(namedEnd, dateEnd, namedStart) ?? Math.max(0, sprints.length - 1));
  return { startIndex, endIndex };
}

function barBounds(startIndex: number, endIndex: number, columnWidth: number) {
  return {
    left: startIndex * columnWidth + 10,
    width: Math.max(44, (endIndex - startIndex + 1) * columnWidth - 20),
  };
}

function buildBar(item: BarSource, sprints: GanttSprint[], nameToIndex: Map<string, number>, today: Date): GanttBarLayout {
  const columnWidth = sprints[0]?.width ?? ZOOM_COLUMN_WIDTH.sprint;
  const { startIndex, endIndex } = barIndexes(item, sprints, nameToIndex);
  return {
    ...barBounds(startIndex, endIndex, columnWidth),
    tone: statusTone(item.status, item.committed_end_date ?? null, item.actual_end_date ?? null, today),
    label: item.title,
  };
}

function compareEpics(left: EpicRecord, right: EpicRecord): number {
  if (left.is_synthetic !== right.is_synthetic) return left.is_synthetic ? 1 : -1;
  return naturalCompare(left.title, right.title);
}

function compareFeatures(left: FeatureRecord, right: FeatureRecord): number {
  return naturalCompare(left.title, right.title);
}

function epicRow(epic: EpicRecord, sprints: GanttSprint[], nameToIndex: Map<string, number>, today: Date): GanttRow {
  return {
    id: epic.id,
    parentId: null,
    type: 'epic',
    title: epic.title,
    status: epic.status,
    owner: epic.owner,
    depth: 0,
    isSynthetic: epic.is_synthetic,
    featureCount: epic.features?.length ?? 0,
    bar: buildBar(epic, sprints, nameToIndex, today),
  };
}

function featureRow(feature: FeatureRecord, epic: EpicRecord, sprints: GanttSprint[], nameToIndex: Map<string, number>, today: Date): GanttRow {
  return {
    id: feature.id,
    parentId: epic.id,
    type: 'feature',
    title: feature.title,
    status: feature.status,
    owner: feature.owner,
    depth: 1,
    isSynthetic: epic.is_synthetic,
    featureCount: 0,
    bar: buildBar(feature, sprints, nameToIndex, today),
  };
}

function buildRows(epics: EpicRecord[], sprints: GanttSprint[], nameToIndex: Map<string, number>, today: Date): GanttRow[] {
  return [...epics].sort(compareEpics).flatMap((epic) => [
    epicRow(epic, sprints, nameToIndex, today),
    ...[...(epic.features ?? [])].sort(compareFeatures).map((feature) => featureRow(feature, epic, sprints, nameToIndex, today)),
  ]);
}

function todayLineLeft(sprints: GanttSprint[], today: Date): number | null {
  const index = sprintIndexForDate(sprints, today);
  if (index === null) return null;
  const sprint = sprints[index];
  const spanDays = Math.max(1, Math.round((sprint.end.getTime() - sprint.start.getTime()) / DAY_MS) + 1);
  const elapsedDays = Math.max(0, Math.round((today.getTime() - sprint.start.getTime()) / DAY_MS));
  return sprint.left + Math.min(1, elapsedDays / spanDays) * sprint.width;
}

export function buildGanttLayout(epics: EpicRecord[], drafts: SprintDraft[], zoom: GanttZoomLevel, todayInput: Date): GanttLayout {
  const today = new Date(todayInput.getFullYear(), todayInput.getMonth(), todayInput.getDate());
  const importDate = earliestDate(epics.flatMap(candidateDates)) ?? today;
  const columnWidth = ZOOM_COLUMN_WIDTH[zoom];
  const sprints = buildSprints(epics, drafts, importDate, columnWidth);
  const nameToIndex = sprintIndexByName(sprints);

  return {
    rows: buildRows(epics, sprints, nameToIndex, today),
    sprints,
    width: Math.max(sprints.length * columnWidth, columnWidth),
    columnWidth,
    rowHeight: ROW_HEIGHT,
    todayLeft: todayLineLeft(sprints, today),
  };
}

export function visibleGanttRows(rows: GanttRow[], expandedEpicIds: Set<string>): GanttRow[] {
  return rows.filter((row) => row.type === 'epic' || (row.parentId !== null && expandedEpicIds.has(row.parentId)));
}

export function formatSprintRange(start: Date, end: Date): string {
  const fmt = new Intl.DateTimeFormat(undefined, { month: 'short', day: 'numeric' });
  return `${fmt.format(start)} – ${fmt.format(end)}`;
}
