<script lang="ts">
  import type { GanttRow } from '../lib/ganttLayout';

  export let barRow: GanttRow;
  export let barTop = 0;
  export let barRowHeight = 42;
</script>

<div
  class="gantt-bar"
  class:epic={barRow.type === 'epic'}
  class:feature={barRow.type === 'feature'}
  class:synthetic={barRow.isSynthetic}
  class:not-started={barRow.bar.tone === 'not-started'}
  class:in-progress={barRow.bar.tone === 'in-progress'}
  class:completed={barRow.bar.tone === 'completed'}
  class:missed={barRow.bar.tone === 'missed'}
  style={`left: ${barRow.bar.left}px; top: ${barTop + 8}px; width: ${barRow.bar.width}px; height: ${Math.max(18, barRowHeight - 16)}px`}
  title={`${barRow.title} · ${barRow.status || 'No status'}`}
>
  <span>{barRow.title}</span>
</div>

<style>
  .gantt-bar {
    position: absolute;
    z-index: 2;
    display: flex;
    align-items: center;
    min-width: 44px;
    overflow: hidden;
    border: 1px solid color-mix(in srgb, var(--text2) 18%, transparent);
    border-radius: 999px;
    background: var(--bg4);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.22);
    color: var(--text);
    font-size: 12px;
    letter-spacing: -0.01em;
    white-space: nowrap;
  }

  .gantt-bar span {
    overflow: hidden;
    padding: 0 12px;
    text-overflow: ellipsis;
  }

  .gantt-bar.epic {
    height: 28px;
    border-radius: 7px;
    font-weight: 500;
  }

  .gantt-bar.feature {
    opacity: 0.9;
    height: 22px;
    border-radius: 999px;
  }

  .gantt-bar.not-started {
    border-color: color-mix(in srgb, var(--text3) 45%, transparent);
    background: color-mix(in srgb, var(--bg4) 76%, var(--text3));
    color: var(--text2);
  }

  .gantt-bar.in-progress {
    border-color: color-mix(in srgb, var(--blue-light) 50%, transparent);
    background: color-mix(in srgb, var(--blue) 44%, var(--bg3));
  }

  .gantt-bar.completed {
    border-color: color-mix(in srgb, var(--green-light) 50%, transparent);
    background: color-mix(in srgb, var(--green) 46%, var(--bg3));
  }

  .gantt-bar.missed {
    border-color: color-mix(in srgb, var(--coral-light) 55%, transparent);
    background: color-mix(in srgb, var(--coral) 50%, var(--bg3));
  }

  .gantt-bar.synthetic {
    border-color: var(--border);
    background: color-mix(in srgb, var(--bg4) 62%, transparent);
    box-shadow: none;
    color: var(--text3);
    font-style: italic;
  }
</style>
