<script lang="ts">
  export let value = '';
  export let sprintName = '';
  export let snapToSprint = true;
  export let disabled = false;
  export let error: string | null = null;
  export let label = 'Target date';
  export let showLabel = true;
  export let showSnapToggle = true;
  export let compact = false;
  export let inputAriaLabel: string | undefined = undefined;

  const isoDatePattern = /\d{4}-\d{2}-\d{2}/g;

  function sprintBoundaryDate(name: string): string | null {
    const matches = name.match(isoDatePattern);
    if (!matches || matches.length === 0) return null;
    return matches[matches.length - 1];
  }

  $: boundaryDate = sprintBoundaryDate(sprintName);
  $: if (showSnapToggle && snapToSprint && boundaryDate && !value) {
    value = boundaryDate;
  }

  function handleSnapChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    snapToSprint = target.checked;
    if (snapToSprint && boundaryDate) {
      value = boundaryDate;
    }
  }
</script>

<div class="date-picker" class:compact>
  {#if showLabel}
    <label class="date-label">
      {label}
      <input
        class:error={!!error}
        type="date"
        bind:value
        {disabled}
        aria-invalid={!!error}
        aria-label={inputAriaLabel}
      />
    </label>
  {:else}
    <input
      class:error={!!error}
      class:compact-input={compact}
      type="date"
      bind:value
      {disabled}
      aria-invalid={!!error}
      aria-label={inputAriaLabel ?? label}
    />
  {/if}

  {#if showSnapToggle}
    <label class="snap-toggle" title={boundaryDate ? `Sprint boundary: ${boundaryDate}` : 'No sprint boundary date available'}>
      <input type="checkbox" bind:checked={snapToSprint} disabled={disabled} on:change={handleSnapChange} />
      <span>Snap to sprint boundary</span>
    </label>
  {/if}

  {#if showSnapToggle && snapToSprint && boundaryDate}
    <div class="date-hint">Using sprint boundary {boundaryDate}</div>
  {:else if showSnapToggle && snapToSprint}
    <div class="date-hint muted">No sprint boundary date detected</div>
  {/if}

  {#if error}
    <div class="date-error">{error}</div>
  {/if}
</div>

<style>
  .date-picker {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .date-picker.compact {
    gap: 6px;
  }

  .date-label {
    display: flex;
    flex-direction: column;
    gap: 6px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  input[type='date'] {
    min-width: 150px;
    height: 36px;
    padding: 0 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 14px;
  }

  input[type='date'].compact-input {
    min-width: 132px;
    height: 34px;
  }

  input[type='date']:focus {
    outline: none;
    border-color: var(--accent2);
  }

  input[type='date'].error {
    border-color: var(--coral);
  }

  .snap-toggle {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    color: var(--text2);
    font-size: 13px;
  }

  .snap-toggle input {
    accent-color: var(--accent);
  }

  .date-hint,
  .date-error {
    font-size: 12px;
  }

  .date-hint {
    color: var(--accent);
  }

  .date-hint.muted {
    color: var(--text3);
  }

  .date-error {
    color: var(--coral);
  }
</style>
