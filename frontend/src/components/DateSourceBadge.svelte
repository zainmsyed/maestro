<script lang="ts">
  export let source: string;

  const labels: Record<string, string> = {
    imported: 'imported',
    pm_assigned: 'pm assigned',
    inherited: 'inherited',
    manual: 'manual',
  };

  const tones: Record<string, string> = {
    imported: 'badge-blue',
    pm_assigned: 'badge-green',
    inherited: 'badge-purple',
    manual: 'badge-amber',
  };

  $: display = labels[source] || source || '—';
  $: tone = tones[source] || 'badge-muted';
</script>

<span class="date-source-badge {tone}">{display}</span>

<style>
  .date-source-badge {
    display: inline-flex;
    padding: 2px 8px;
    border-radius: 999px;
    font-size: 11px;
    letter-spacing: 0.04em;
    text-transform: lowercase;
    line-height: 1.4;
  }

  .badge-blue {
    border: 1px solid color-mix(in srgb, var(--blue) 40%, transparent);
    background: color-mix(in srgb, var(--blue) 12%, transparent);
    color: var(--blue-light);
  }

  .badge-green {
    border: 1px solid color-mix(in srgb, var(--green) 40%, transparent);
    background: color-mix(in srgb, var(--green) 12%, transparent);
    color: var(--green-light);
  }

  .badge-purple {
    border: 1px solid color-mix(in srgb, var(--purple) 40%, transparent);
    background: color-mix(in srgb, var(--purple) 12%, transparent);
    color: var(--purple-light);
  }

  .badge-amber {
    border: 1px solid color-mix(in srgb, var(--amber) 40%, transparent);
    background: color-mix(in srgb, var(--amber) 12%, transparent);
    color: var(--amber);
  }

  .badge-muted {
    border: 1px solid var(--border);
    background: var(--bg3);
    color: var(--text3);
  }
</style>
