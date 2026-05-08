<script lang="ts">
  import GanttBar from './GanttBar.svelte';
  import SprintHeader from './SprintHeader.svelte';
  import TodayLine from './TodayLine.svelte';
  import type { GanttRow, GanttSprint } from '../lib/ganttLayout';

  export let gridRows: GanttRow[] = [];
  export let gridSprints: GanttSprint[] = [];
  export let gridWidth = 0;
  export let gridRowHeight = 42;
  export let gridColumnWidth = 172;
  export let gridTodayLeft: number | null = null;

  $: bodyHeight = gridRows.length * gridRowHeight;
  $: gridHeight = bodyHeight + 58;
</script>

<div class="gantt-grid" style={`width: ${gridWidth}px`}>
  <SprintHeader headerSprints={gridSprints} headerColumnWidth={gridColumnWidth} />

  <div class="grid-body" style={`height: ${bodyHeight}px; width: ${gridWidth}px`}>
    {#each gridSprints as sprint}
      <div class="vertical-line" style={`left: ${sprint.left}px; height: ${bodyHeight}px`}></div>
    {/each}
    <div class="vertical-line end" style={`left: ${gridWidth}px; height: ${bodyHeight}px`}></div>

    {#each gridRows as row, index}
      <div
        class="row-band"
        class:epic={row.type === 'epic'}
        class:feature={row.type === 'feature'}
        class:synthetic={row.isSynthetic}
        style={`top: ${index * gridRowHeight}px; height: ${gridRowHeight}px`}
      ></div>
      <GanttBar barRow={row} barTop={index * gridRowHeight} barRowHeight={gridRowHeight} />
    {/each}

    <TodayLine left={gridTodayLeft} height={gridHeight} />
  </div>
</div>

<style>
  .gantt-grid {
    position: relative;
    min-width: 100%;
  }

  .grid-body {
    position: relative;
    min-width: 100%;
    background:
      linear-gradient(to bottom, transparent calc(100% - 1px), color-mix(in srgb, var(--border) 72%, transparent) 0) 0 0 / 100% 48px,
      color-mix(in srgb, var(--bg) 94%, transparent);
  }

  .vertical-line {
    position: absolute;
    top: 0;
    width: 1px;
    background: color-mix(in srgb, var(--border) 82%, transparent);
  }

  .vertical-line.end {
    background: var(--border);
  }

  .row-band {
    position: absolute;
    left: 0;
    right: 0;
    border-bottom: 1px solid color-mix(in srgb, var(--border) 72%, transparent);
    pointer-events: none;
  }

  .row-band.epic {
    border-top: 1px solid color-mix(in srgb, var(--accent) 10%, transparent);
    background: color-mix(in srgb, var(--accent) 4%, transparent);
  }

  .row-band.feature {
    background: color-mix(in srgb, var(--text) 1.5%, transparent);
  }

  .row-band.synthetic {
    background: color-mix(in srgb, var(--bg4) 18%, transparent);
  }
</style>
