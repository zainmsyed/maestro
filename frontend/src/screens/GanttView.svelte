<script lang="ts">
  import { onMount } from 'svelte';
  import { Gantt } from '@svar-ui/svelte-gantt';
  import type { IConfig, IColumnConfig, IMarker, IApi } from '@svar-ui/svelte-gantt';
  import { api, formatApiError, type EpicRecord } from '../lib/api';
  import { buildSvarTasks, type SvarTask } from '../lib/svarBridge';

  let loading = true;
  let error: string | null = null;
  let epics: EpicRecord[] = [];
  let collapsedEpics = new Set<string>();

  const links: IConfig['links'] = [];
  const cellWidth = 28;
  const dayMs = 24 * 60 * 60 * 1000;
  const scales: IConfig['scales'] = [
    { unit: 'month', step: 1, format: '%F %Y' },
  ];
  const columns: IColumnConfig[] = [
    { id: 'text', header: 'Work item', flexgrow: 1 },
    { id: 'owner', header: 'Owner', width: 160 },
  ];

  async function load() {
    loading = true;
    error = null;
    try {
      epics = await api.get<EpicRecord[]>('/epics');
    } catch (err) {
      error = formatApiError(err, 'Could not load the Gantt view.');
    } finally {
      loading = false;
    }
  }

  function addDays(date: Date, days: number): Date {
    return new Date(date.getTime() + days * dayMs);
  }

  function minTaskDate(items: SvarTask[]): Date | null {
    const dates = items.map((task) => task.start).filter((date): date is Date => date instanceof Date);
    if (dates.length === 0) return null;
    return new Date(Math.min(...dates.map((date) => date.getTime())));
  }

  function maxTaskDate(items: SvarTask[]): Date | null {
    const dates = items.map((task) => task.end).filter((date): date is Date => date instanceof Date);
    if (dates.length === 0) return null;
    return new Date(Math.max(...dates.map((date) => date.getTime())));
  }

  function parseDate(value: string | null | undefined): Date | null {
    if (!value) return null;
    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? null : date;
  }

  function sprintBoundariesFor(records: EpicRecord[]): Date[] {
    const byDay = new Map<string, Date>();
    for (const epic of records) {
      const start = parseDate(epic.sprint_start);
      if (start) byDay.set(start.toISOString().slice(0, 10), start);
    }
    return [...byDay.values()].sort((a, b) => a.getTime() - b.getTime());
  }

  function buildMarkers(boundaries: Date[]): IMarker[] {
    return [
      { start: new Date(), text: 'Today', css: 'maestro-today-marker' },
      ...boundaries.map((start) => ({ start, text: '', css: 'maestro-sprint-boundary' })),
    ];
  }

  function dataIdSelector(id: string): string {
    return Array.from(`:${id}`)
      .map((char) => (/^[a-zA-Z0-9_:-]$/.test(char) ? char : `\\${char.charCodeAt(0).toString(16)} `))
      .join('');
  }

  const variantStyles: Record<SvarTask['status_variant'] | 'synthetic', string> = {
    synthetic: 'background: var(--gantt-muted) !important; border-color: var(--gantt-muted-border) !important; color: var(--text2) !important; opacity: 0.74;',
    missed: 'background: var(--gantt-red) !important; border-color: var(--gantt-red-border) !important; color: var(--text) !important;',
    done: 'background: var(--gantt-green) !important; border-color: var(--gantt-green-border) !important; color: var(--text) !important;',
    'in-progress': 'background: var(--gantt-blue) !important; border-color: var(--gantt-blue-border) !important; color: var(--text) !important;',
    'not-started': 'background: var(--gantt-grey) !important; border-color: var(--gantt-grey-border) !important; color: var(--text2) !important;',
  };

  function styleForVariant(task: SvarTask): string {
    return variantStyles[task.is_synthetic ? 'synthetic' : task.status_variant];
  }

  function taskStyleRules(items: SvarTask[]): string {
    return items
      .map((task) => {
        const selector = `.gantt-card .wx-bar[data-id="${dataIdSelector(task.id)}"]`;
        const contentSelector = `.gantt-card .wx-bar[data-id="${dataIdSelector(task.id)}"] .wx-content`;
        const syntheticLabel = task.is_synthetic ? `${contentSelector}{font-style:italic;color:var(--text2)!important;}` : '';
        return `${selector}{${styleForVariant(task)}}${syntheticLabel}`;
      })
      .join('\n');
  }

  function onGanttInit(ganttApi: IApi) {
    ganttApi.intercept('open-task', ({ id }: { id: string | number; mode: boolean }) => {
      const sid = String(id);
      collapsedEpics = new Set(collapsedEpics);
      if (collapsedEpics.has(sid)) {
        collapsedEpics.delete(sid);
      } else {
        collapsedEpics.add(sid);
      }
      return false;
    });
  }

  function collapseAll() {
    const all = tasks.filter((t: SvarTask) => t.source_type === 'epic').map((t) => t.id);
    collapsedEpics = new Set(all);
  }

  function expandAll() {
    collapsedEpics = new Set();
  }

  $: tasks = buildSvarTasks(epics);
  $: visibleTasks = tasks.map((task: SvarTask) => {
    if (task.source_type === 'epic') {
      return { ...task, open: !collapsedEpics.has(task.id) };
    }
    return task;
  });
  $: featureCount = tasks.filter((task: SvarTask) => task.source_type === 'feature').length;
  $: missingDateCount = tasks.filter((task: SvarTask) => task.source_type === 'feature' && task.missing_committed_end).length;
  $: timelineStart = minTaskDate(visibleTasks) ? addDays(minTaskDate(visibleTasks) as Date, -7) : null;
  $: timelineEnd = maxTaskDate(visibleTasks) ? addDays(maxTaskDate(visibleTasks) as Date, 14) : null;
  $: sprintBoundaries = sprintBoundariesFor(epics);
  $: markers = buildMarkers(sprintBoundaries);
  $: ganttStart = timelineStart ?? undefined;
  $: ganttEnd = timelineEnd ?? undefined;
  $: ganttStyleRules = taskStyleRules(visibleTasks);

  onMount(load);
</script>

<svelte:head>
  <link rel="stylesheet" href="https://cdn.svar.dev/fonts/wxi/wx-icons.css" />
</svelte:head>

<section class="gantt-view" aria-label="Gantt roadmap view">
  <div class="screen-header">
    <div>
      <p class="eyebrow">Gantt</p>
      <h1>Timeline roadmap</h1>
      <p class="intro">
        A cozy timeline view of epics and features, grouped into expandable epic rows.
      </p>
    </div>
    <div class="header-summary" aria-label="Gantt summary">
      <div><span>{epics.length}</span><small>Epics</small></div>
      <div><span>{featureCount}</span><small>Features</small></div>
      <div><span>{missingDateCount}</span><small>Missing end</small></div>
    </div>
  </div>

  {#if error}
    <div class="error-banner" role="alert">
      <strong>Timeline unavailable.</strong>
      <span>{error}</span>
      <button type="button" on:click={load}>Try again</button>
    </div>
  {/if}

  {#if loading}
    <section class="state-card" aria-live="polite">
      <p class="eyebrow">Loading</p>
      <h2>Building the timeline</h2>
      <p>Fetching epics and features, then arranging them into SVAR tasks…</p>
    </section>
  {:else if epics.length === 0}
    <section class="state-card">
      <p class="eyebrow">Empty</p>
      <h2>No roadmap data yet</h2>
      <p>Import a CSV to populate the Gantt timeline.</p>
    </section>
  {:else if visibleTasks.length === 0}
    <section class="state-card">
      <p class="eyebrow">Empty</p>
      <h2>No schedulable tasks</h2>
      <p>No epics or features had enough date information to render on the timeline.</p>
    </section>
  {:else}
    <div class="gantt-toolbar">
      <button class="action-link" type="button" on:click={collapseAll}>Collapse all</button>
      <button class="action-link" type="button" on:click={expandAll}>Expand all</button>
    </div>
    <div class="gantt-card">
      <Gantt
        tasks={visibleTasks}
        {links}
        {scales}
        {columns}
        {markers}
        readonly={true}
        cellBorders="full"
        start={ganttStart}
        end={ganttEnd}
        lengthUnit="day"
        {cellWidth}
        cellHeight={48}
        scaleHeight={44}
        init={onGanttInit}
      />
      {@html `<style>${ganttStyleRules}</style>`}
    </div>
  {/if}
</section>

<style>
  .gantt-view {
    display: flex;
    flex-direction: column;
    gap: 24px;
    width: min(1440px, 100%);
    min-height: 100%;
    margin: 0 auto;
  }

  .screen-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 24px;
  }

  .eyebrow {
    margin: 0 0 10px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  h1,
  h2 {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-weight: 300;
    letter-spacing: -0.05em;
  }

  h1 {
    font-size: clamp(34px, 5vw, 56px);
    line-height: 1.02;
  }

  h2 {
    font-size: 30px;
    line-height: 1.12;
  }

  .intro,
  .state-card p,
  .error-banner span {
    color: var(--text2);
  }

  .intro {
    max-width: 680px;
    margin: 12px 0 0;
    font-size: 15px;
  }

  .header-summary {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
  }

  .header-summary div,
  .state-card,
  .gantt-card,
  .error-banner {
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
  }

  .header-summary div {
    min-width: 120px;
    padding: 14px 16px;
  }

  .header-summary span {
    display: block;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 28px;
    font-weight: 300;
  }

  .header-summary small {
    display: block;
    margin-top: 4px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .state-card {
    padding: 32px;
  }

  .state-card p:last-child {
    margin: 12px 0 0;
    font-size: 15px;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 18px;
    border-color: color-mix(in srgb, var(--coral) 45%, transparent);
  }

  .error-banner strong {
    color: var(--coral-light);
    font-weight: 500;
  }

  .error-banner button {
    margin-left: auto;
    padding: 8px 14px;
    border: 1px solid var(--accent2);
    border-radius: 6px;
    background: transparent;
    color: var(--accent);
    cursor: pointer;
  }

  .error-banner button:hover {
    background: var(--accent-a8);
  }

  .gantt-toolbar {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 16px;
    padding: 12px 14px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
  }

  .action-link {
    padding: 0;
    border: 0;
    background: transparent;
    color: var(--text3);
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 12px;
    letter-spacing: 0.04em;
    text-decoration: underline;
    text-underline-offset: 3px;
  }

  .action-link:hover {
    color: var(--accent);
  }

  .gantt-card {
    --gantt-grey: color-mix(in srgb, var(--text2) 42%, var(--bg4));
    --gantt-grey-border: color-mix(in srgb, var(--text2) 58%, transparent);
    --gantt-blue: var(--blue);
    --gantt-blue-border: var(--blue-light);
    --gantt-green: var(--green);
    --gantt-green-border: var(--green-light);
    --gantt-red: var(--red);
    --gantt-red-border: var(--coral-light);
    --gantt-muted: color-mix(in srgb, var(--text2) 24%, var(--bg3));
    --gantt-muted-border: color-mix(in srgb, var(--text2) 36%, transparent);
    position: relative;
    min-height: 620px;
    height: calc(100vh - 220px);
    overflow: hidden;
  }

  :global(.gantt-card .wx-gantt) {
    height: 100%;
    color: var(--text);
  }

  :global(.gantt-card .wx-marker) {
    pointer-events: none;
  }

  :global(.gantt-card .wx-marker.maestro-today-marker) {
    width: 2px;
    background: var(--coral-light);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--coral-light) 24%, transparent);
  }

  :global(.gantt-card .wx-marker.maestro-today-marker .wx-content) {
    min-width: auto;
    padding: 3px 6px;
    border: 1px solid color-mix(in srgb, var(--coral-light) 45%, transparent);
    border-radius: 999px;
    background: var(--bg2);
    color: var(--coral-light);
    font-family: var(--font-mono);
    font-size: 11px;
    letter-spacing: 0.04em;
  }

  :global(.gantt-card .wx-marker.maestro-sprint-boundary) {
    width: 1px;
    background: color-mix(in srgb, var(--accent) 34%, transparent);
  }

  :global(.gantt-card .wx-marker.maestro-sprint-boundary .wx-content) {
    display: none;
  }

  :global(.gantt-card .wx-toggle-icon) {
    color: var(--text3) !important;
    cursor: pointer;
  }

  :global(.gantt-card .wx-toggle-icon:hover) {
    color: var(--text) !important;
  }

  @media (max-width: 980px) {
    .screen-header {
      flex-direction: column;
    }

    .header-summary {
      width: 100%;
    }
  }

  @media (max-width: 720px) {
    .header-summary {
      grid-template-columns: 1fr;
    }
  }
</style>