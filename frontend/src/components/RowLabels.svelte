<script lang="ts">
  import type { GanttRow } from '../lib/ganttLayout';

  export let labelRows: GanttRow[] = [];
  export let labelExpandedEpicIds: Set<string> = new Set();
  export let labelRowHeight = 42;
  export let handleToggleEpic: (id: string) => void = () => {};
</script>

<div class="row-labels" role="tree" aria-label="Roadmap rows">
  <div class="label-header">
    <span>Work item</span>
    <span>Type</span>
  </div>
  {#each labelRows as row}
    <div
      class="label-row"
      class:epic={row.type === 'epic'}
      class:feature={row.type === 'feature'}
      class:synthetic={row.isSynthetic}
      style={`height: ${labelRowHeight}px`}
      role="treeitem"
      aria-selected="false"
      aria-expanded={row.type === 'epic' ? labelExpandedEpicIds.has(row.id) : undefined}
    >
      {#if row.type === 'epic'}
        <button
          class="chevron"
          class:expanded={labelExpandedEpicIds.has(row.id)}
          type="button"
          aria-label={`${labelExpandedEpicIds.has(row.id) ? 'Collapse' : 'Expand'} ${row.title}`}
          on:click={() => handleToggleEpic(row.id)}
        >
          ›
        </button>
      {:else}
        <span class="indent" aria-hidden="true"></span>
      {/if}

      <div class="label-main">
        <div class="label-title">{row.title}</div>
        <div class="label-meta">
          {#if row.owner}{row.owner}{:else}No owner{/if}
          · {row.status || 'No status'}
          {#if row.type === 'epic'} · {row.featureCount} features{/if}
        </div>
      </div>

      <span class="type-pill type-{row.type}" class:hide-pill={row.isSynthetic}>{row.type}</span>
    </div>
  {/each}
</div>

<style>
  .row-labels {
    width: 100%;
    min-width: 0;
    border-right: 1px solid var(--border);
    background: color-mix(in srgb, var(--bg2) 92%, transparent);
  }

  .label-header {
    position: sticky;
    top: 0;
    z-index: 4;
    display: grid;
    grid-template-columns: minmax(0, 1fr) 96px;
    align-items: center;
    height: 58px;
    padding: 0 14px;
    border-bottom: 1px solid var(--border);
    background: color-mix(in srgb, var(--bg2) 96%, transparent);
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  .label-row {
    display: grid;
    grid-template-columns: 24px minmax(0, 1fr) 96px;
    align-items: center;
    gap: 8px;
    padding: 0 12px;
    border-bottom: 1px solid color-mix(in srgb, var(--border) 72%, transparent);
  }

  .label-row.epic {
    border-top: 1px solid color-mix(in srgb, var(--accent) 10%, transparent);
    background: color-mix(in srgb, var(--accent) 5%, transparent);
  }

  .label-row.feature {
    background: color-mix(in srgb, var(--text) 1.5%, transparent);
  }

  .label-row.synthetic {
    color: var(--text3);
    font-style: italic;
  }

  .chevron {
    width: 22px;
    height: 22px;
    display: grid;
    place-items: center;
    border: 1px solid transparent;
    border-radius: 5px;
    background: transparent;
    color: var(--text2);
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
    transition: transform 0.15s ease, background 0.15s ease, color 0.15s ease;
  }

  .chevron:hover {
    background: var(--bg4);
    color: var(--accent);
  }

  .chevron.expanded {
    transform: rotate(90deg);
    color: var(--accent);
  }

  .indent {
    width: 22px;
    height: 1px;
    justify-self: center;
    background: var(--border2);
  }

  .label-main {
    min-width: 0;
  }

  .label-title,
  .label-meta {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .label-title {
    color: var(--text);
    font-size: 13px;
    letter-spacing: -0.02em;
  }

  .epic .label-title {
    color: var(--accent);
    font-size: 12px;
    letter-spacing: 0.09em;
    text-transform: uppercase;
  }

  .feature .label-title {
    padding-left: 4px;
  }

  .synthetic .label-title {
    color: var(--text3);
  }

  .label-meta {
    margin-top: 2px;
    color: var(--text3);
    font-size: 11px;
  }

  .type-pill {
    width: fit-content;
    max-width: 96px;
    justify-self: end;
    overflow: hidden;
    padding: 3px 8px;
    border: 1px solid var(--border);
    border-radius: 999px;
    color: var(--text2);
    font-size: 11px;
    text-transform: capitalize;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .type-epic {
    border-color: color-mix(in srgb, var(--blue) 35%, transparent);
    background: color-mix(in srgb, var(--blue) 8%, transparent);
    color: var(--blue-light);
  }

  .type-feature {
    border-color: color-mix(in srgb, var(--purple) 35%, transparent);
    background: color-mix(in srgb, var(--purple) 8%, transparent);
    color: var(--purple-light);
  }

  .type-pill.hide-pill {
    visibility: hidden;
  }
</style>
