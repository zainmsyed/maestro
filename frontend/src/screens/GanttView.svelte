<script lang="ts">
  import { onMount } from 'svelte';
  import GanttGrid from '../components/GanttGrid.svelte';
  import RowLabels from '../components/RowLabels.svelte';
  import { api, formatApiError, type EpicRecord } from '../lib/api';
  import { buildGanttLayout, visibleGanttRows, type GanttZoomLevel } from '../lib/ganttLayout';
  import { project } from '../stores/project';

  let loading = true;
  let error: string | null = null;
  let epics: EpicRecord[] = [];
  let zoom: GanttZoomLevel = 'sprint';
  let expandedEpicIds = new Set<string>();

  async function load() {
    loading = true;
    error = null;
    try {
      epics = await api.get<EpicRecord[]>('/epics');
      expandedEpicIds = new Set(epics.filter((epic) => !epic.is_synthetic).map((epic) => epic.id));
    } catch (err) {
      error = formatApiError(err, 'Could not load timeline data.');
    } finally {
      loading = false;
    }
  }

  function toggleEpic(id: string) {
    const next = new Set(expandedEpicIds);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    expandedEpicIds = next;
  }

  onMount(load);

  $: layout = buildGanttLayout(epics, $project.sprintDrafts, zoom, new Date());
  $: visibleRows = visibleGanttRows(layout.rows, expandedEpicIds);
  $: epicCount = epics.length;
  $: featureCount = epics.reduce((total, epic) => total + (epic.features?.length ?? 0), 0);
</script>

<section class="gantt-view" aria-label="Timeline view">
  <div class="screen-header">
    <div>
      <p class="eyebrow">Timeline</p>
      <h1>Roadmap timeline</h1>
      <p class="intro">Sprint-aligned epics and features with status colour, today marker, and expandable rows.</p>
    </div>

    <div class="header-actions">
      <div class="summary-card"><span>{epicCount}</span><small>Epics</small></div>
      <div class="summary-card"><span>{featureCount}</span><small>Features</small></div>
      <label class="zoom-control">
        <span>Zoom</span>
        <select bind:value={zoom} aria-label="Timeline zoom level">
          <option value="sprint">Sprint</option>
          <option value="month">Month</option>
          <option value="quarter">Quarter</option>
        </select>
      </label>
    </div>
  </div>

  {#if error}
    <div class="error-banner">{error}</div>
  {/if}

  <div class="timeline-shell" class:loading>
    {#if loading}
      <div class="loading-state">Loading timeline…</div>
    {:else if layout.rows.length === 0}
      <div class="empty-state">No roadmap data found. Import a CSV to populate the timeline.</div>
    {:else}
      <div class="labels-pane">
        <RowLabels
          labelRows={visibleRows}
          labelExpandedEpicIds={expandedEpicIds}
          labelRowHeight={layout.rowHeight}
          handleToggleEpic={toggleEpic}
        />
      </div>
      <div class="grid-pane">
        <GanttGrid
          gridRows={visibleRows}
          gridSprints={layout.sprints}
          gridWidth={layout.width}
          gridRowHeight={layout.rowHeight}
          gridColumnWidth={layout.columnWidth}
          gridTodayLeft={layout.todayLeft}
        />
      </div>
    {/if}
  </div>
</section>

<style>
  .gantt-view {
    display: flex;
    flex-direction: column;
    gap: 20px;
    min-height: 100%;
    width: 100%;
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

  h1 {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-size: clamp(34px, 5vw, 56px);
    font-weight: 300;
    letter-spacing: -0.05em;
    line-height: 1.02;
  }

  .intro {
    max-width: 760px;
    margin: 12px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .header-actions {
    display: flex;
    align-items: stretch;
    gap: 12px;
  }

  .summary-card,
  .zoom-control {
    min-width: 108px;
    padding: 12px 14px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--surface-glass);
  }

  .summary-card span {
    display: block;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 26px;
    font-weight: 300;
  }

  .summary-card small,
  .zoom-control span {
    display: block;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .zoom-control select {
    width: 100%;
    margin-top: 8px;
    border: 1px solid var(--border2);
    border-radius: 6px;
    background: var(--bg);
    color: var(--text2);
    font: inherit;
    font-size: 13px;
  }

  .error-banner,
  .loading-state,
  .empty-state {
    padding: 16px 18px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--bg2);
    color: var(--text2);
  }

  .error-banner {
    border-color: color-mix(in srgb, var(--coral) 45%, transparent);
    color: var(--coral);
  }

  .timeline-shell {
    min-height: 520px;
    flex: 1;
    display: grid;
    grid-template-columns: minmax(420px, 34vw) minmax(0, 1fr);
    overflow: hidden;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
  }

  .timeline-shell.loading {
    display: block;
    padding: 18px;
  }

  .labels-pane,
  .grid-pane {
    min-height: 0;
  }

  .labels-pane {
    overflow: visible;
  }

  .grid-pane {
    position: relative;
    overflow-x: auto;
    overflow-y: hidden;
  }

  @media (max-width: 980px) {
    .screen-header,
    .header-actions {
      flex-direction: column;
    }

    .timeline-shell {
      grid-template-columns: 340px minmax(0, 1fr);
    }
  }
</style>
