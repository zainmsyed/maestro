<script lang="ts">
  import { formatSprintRange, type GanttSprint } from '../lib/ganttLayout';

  export let headerSprints: GanttSprint[] = [];
  export let headerColumnWidth = 172;
</script>

<div class="sprint-header" style={`width: ${Math.max(headerSprints.length * headerColumnWidth, headerColumnWidth)}px`}>
  {#each headerSprints as sprint}
    <div class="sprint-cell" style={`width: ${headerColumnWidth}px`}>
      <span class="sprint-name">{sprint.name}</span>
      <span class="sprint-range">{formatSprintRange(sprint.start, sprint.end)}</span>
    </div>
  {/each}
</div>

<style>
  .sprint-header {
    position: sticky;
    top: 0;
    z-index: 4;
    display: flex;
    height: 58px;
    border-bottom: 1px solid var(--border);
    background: color-mix(in srgb, var(--bg2) 96%, transparent);
  }

  .sprint-cell {
    flex: 0 0 auto;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 3px;
    padding: 0 14px;
    border-right: 1px solid var(--border);
  }

  .sprint-name {
    overflow: hidden;
    color: var(--text);
    font-size: 13px;
    letter-spacing: 0.04em;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .sprint-range {
    color: var(--text3);
    font-size: 12px;
  }
</style>
