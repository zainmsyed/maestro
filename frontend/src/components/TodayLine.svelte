<script lang="ts">
  export let left: number | null = null;
  export let height = 0;
</script>

{#if left !== null}
  <div class="today-line" style={`left: ${left}px; height: ${height}px`} aria-label="Today">
    <span>today</span>
  </div>
{/if}

<style>
  .today-line {
    position: absolute;
    top: 0;
    z-index: 3;
    width: 1px;
    background: color-mix(in srgb, var(--accent) 72%, transparent);
    pointer-events: none;
  }

  .today-line span {
    position: sticky;
    top: 8px;
    display: inline-flex;
    align-items: center;
    gap: 5px;
    transform: translateX(-50%);
    padding: 2px 7px;
    border: 1px solid var(--accent2);
    border-radius: 999px;
    background: var(--bg);
    color: var(--accent);
    font-size: 11px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .today-line span::before {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 999px;
    background: var(--accent);
  }
</style>
