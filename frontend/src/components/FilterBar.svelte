<script lang="ts">
  import GroupBySelect from './GroupBySelect.svelte';
  import type { FilterState, GroupBy } from '../lib/listView';

  export let filters: FilterState = { epic: '', owner: '', sprint: '', status: '' };
  export let onFilterChange: (filters: FilterState) => void = () => {};
  export let groupBy: GroupBy = 'epic';
  export let epicOptions: string[] = [];
  export let ownerOptions: string[] = [];
  export let sprintOptions: string[] = [];
  export let statusOptions: string[] = [];
  export let onExport: () => void;

  const filterLabels: Record<keyof FilterState, string> = {
    epic: 'Epic',
    owner: 'Owner',
    sprint: 'Sprint',
    status: 'Status',
  };

  function resetFilters() {
    onFilterChange({ epic: '', owner: '', sprint: '', status: '' });
  }

  function updateFilter(key: keyof FilterState, event: Event) {
    const select = event.currentTarget as HTMLSelectElement;
    onFilterChange({ ...filters, [key]: select.value });
  }

  function clearFilter(key: keyof FilterState) {
    onFilterChange({ ...filters, [key]: '' });
  }

  $: activeFilters = (Object.entries(filters) as Array<[keyof FilterState, string]>).filter(([, value]) => Boolean(value));
</script>

<section class="filter-bar" aria-label="List view filters">
  <div class="filters-grid">
    <label>
      <span>Epic</span>
      <select value={filters.epic} on:change={(event) => updateFilter('epic', event)}>
        <option value="">All epics</option>
        {#each epicOptions as option}
          <option value={option}>{option}</option>
        {/each}
      </select>
    </label>

    <label>
      <span>Owner</span>
      <select value={filters.owner} on:change={(event) => updateFilter('owner', event)}>
        <option value="">All owners</option>
        {#each ownerOptions as option}
          <option value={option}>{option}</option>
        {/each}
      </select>
    </label>

    <label>
      <span>Sprint</span>
      <select value={filters.sprint} on:change={(event) => updateFilter('sprint', event)}>
        <option value="">All sprints</option>
        {#each sprintOptions as option}
          <option value={option}>{option}</option>
        {/each}
      </select>
    </label>

    <label>
      <span>Status</span>
      <select value={filters.status} on:change={(event) => updateFilter('status', event)}>
        <option value="">All statuses</option>
        {#each statusOptions as option}
          <option value={option}>{option}</option>
        {/each}
      </select>
    </label>

    <GroupBySelect bind:value={groupBy} />
  </div>

  <div class="actions">
    <button class="secondary" type="button" on:click={resetFilters}>Clear filters</button>
    <button class="primary" type="button" on:click={onExport}>Export CSV</button>
  </div>

  {#if activeFilters.length > 0}
    <div class="active-filters" aria-label="Active filters">
      {#each activeFilters as [key, value]}
        <button class="filter-chip" type="button" on:click={() => clearFilter(key)}>
          <span>{filterLabels[key]}:</span>
          <strong>{value}</strong>
          <small>×</small>
        </button>
      {/each}
    </div>
  {/if}
</section>

<style>
  .filter-bar {
    display: flex;
    flex-wrap: wrap;
    align-items: end;
    justify-content: space-between;
    gap: 20px;
    padding: 18px 20px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: linear-gradient(145deg, rgba(29, 29, 32, 0.94), rgba(22, 22, 24, 0.94));
  }

  .filters-grid {
    display: grid;
    grid-template-columns: repeat(5, minmax(0, 1fr));
    gap: 14px;
    flex: 1;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 6px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  select {
    height: 36px;
    min-width: 120px;
    padding: 0 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 14px;
  }

  select:focus {
    outline: none;
    border-color: var(--accent2);
  }

  .actions {
    display: flex;
    gap: 12px;
  }

  .active-filters {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    width: 100%;
    padding-top: 10px;
    border-top: 1px solid var(--border);
    margin-top: 4px;
  }

  .filter-chip {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    height: 32px;
    padding: 0 12px;
    border: 1px solid rgba(200, 169, 110, 0.35);
    border-radius: 999px;
    background: rgba(200, 169, 110, 0.1);
    color: var(--text2);
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 12px;
    letter-spacing: 0.04em;
    transition: background 0.15s ease;
  }

  .filter-chip:hover {
    background: rgba(200, 169, 110, 0.18);
  }

  .filter-chip span,
  .filter-chip small {
    color: var(--text3);
  }

  .filter-chip strong {
    color: var(--accent);
    font-weight: 500;
  }

  button {
    height: 36px;
    padding: 0 14px;
    border-radius: 6px;
    background: transparent;
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    letter-spacing: 0.04em;
  }

  .primary {
    border: 1px solid var(--accent2);
    color: var(--accent);
  }

  .secondary {
    border: 1px solid var(--border2);
    color: var(--text2);
  }

  .primary:hover,
  .secondary:hover {
    background: rgba(200, 169, 110, 0.08);
  }

  @media (max-width: 1100px) {
    .filter-bar {
      flex-direction: column;
      align-items: stretch;
    }

    .filters-grid {
      grid-template-columns: 1fr 1fr;
    }

    .actions {
      justify-content: flex-end;
    }
  }

  @media (max-width: 720px) {
    .filters-grid {
      grid-template-columns: 1fr;
    }

    .actions {
      flex-direction: column;
    }
  }
</style>
