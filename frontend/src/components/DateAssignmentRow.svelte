<script lang="ts">
  import DatePicker from './DatePicker.svelte';
  import type { DateAssignmentCandidate } from '../lib/api';

  type AssignmentStatus = 'idle' | 'saving' | 'success' | 'error';

  export let candidate: DateAssignmentCandidate;
  export let sprintName = '';
  export let selected = false;
  export let dateValue = '';
  export let snapToSprint = true;
  export let disabled = false;
  export let status: AssignmentStatus = 'idle';
  export let error: string | null = null;
</script>

<div class="assignment-row" class:selected class:error={status === 'error'}>
  <label class="select-cell" aria-label={`Select ${candidate.title}`}>
    <input type="checkbox" bind:checked={selected} {disabled} />
  </label>

  <div class="item-cell">
    <div class="item-title">{candidate.title}</div>
    <div class="item-meta">
      <span>{candidate.work_item_type}</span>
      <span>#{candidate.id}</span>
      {#if candidate.assigned_owner}
        <span>{candidate.assigned_owner}</span>
      {/if}
    </div>
  </div>

  <div class="sprint-cell">
    <span class="label">Sprint</span>
    <span class="sprint-name">{sprintName || 'Unassigned'}</span>
  </div>

  <DatePicker bind:value={dateValue} bind:snapToSprint {sprintName} {disabled} {error} />

  <div class="status-cell" aria-live="polite">
    {#if status === 'saving'}
      Saving…
    {:else if status === 'success'}
      Saved
    {:else if status === 'error'}
      Error
    {/if}
  </div>
</div>

<style>
  .assignment-row {
    display: grid;
    grid-template-columns: 34px minmax(220px, 1fr) minmax(150px, 190px) minmax(220px, 260px) 70px;
    gap: 16px;
    align-items: center;
    padding: 16px 18px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--surface-panel);
  }

  .assignment-row.selected {
    border-color: color-mix(in srgb, var(--accent) 32%, transparent);
    background: color-mix(in srgb, var(--accent) 5%, transparent);
  }

  .assignment-row.error {
    border-color: color-mix(in srgb, var(--coral) 45%, transparent);
  }

  .select-cell {
    display: grid;
    place-items: center;
  }

  .select-cell input {
    width: 18px;
    height: 18px;
    accent-color: var(--accent);
  }

  .item-cell {
    min-width: 0;
  }

  .item-title {
    overflow: hidden;
    color: var(--text);
    font-size: 15px;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .item-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 5px;
    color: var(--text3);
    font-size: 12px;
    text-transform: capitalize;
  }

  .sprint-cell {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .label {
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .sprint-name {
    color: var(--text2);
    font-size: 14px;
  }

  .status-cell {
    color: var(--accent);
    font-size: 12px;
    min-height: 18px;
    text-align: right;
  }

  @media (max-width: 980px) {
    .assignment-row {
      grid-template-columns: 34px minmax(0, 1fr);
    }

    .sprint-cell,
    .status-cell {
      grid-column: 2;
    }
  }
</style>
