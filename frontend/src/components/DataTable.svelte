<script lang="ts">
  import DatePicker from './DatePicker.svelte';
  import type { ListGroup, ListRow, SortKey } from '../lib/listView';

  export let groups: ListGroup[] = [];
  export let sortKey: SortKey = 'title';
  export let sortDirection: 'asc' | 'desc' = 'asc';
  export let loading = false;
  export let savingDateId: string | null = null;
  export let epicOptions: Array<{ id: string; title: string }> = [];
  export let onSort: (key: SortKey) => void;
  export let onSaveDate: (row: ListRow, date: string) => void;
  export let onReassignFeature: (row: ListRow, epicID: string) => void;

  const columns: Array<{ key: SortKey; label: string }> = [
    { key: 'title', label: 'Title' },
    { key: 'type', label: 'Type' },
    { key: 'owner', label: 'Owner' },
    { key: 'sprint', label: 'Sprint' },
    { key: 'original', label: 'Original Date' },
    { key: 'committed', label: 'Committed Date' },
    { key: 'actual', label: 'Actual Date' },
    { key: 'slip', label: 'Slip Events' },
    { key: 'status', label: 'Status' },
    { key: 'health', label: 'Health' },
  ];

  let editingId: string | null = null;
  let editingDate = '';

  function beginEdit(row: ListRow) {
    editingId = row.id;
    editingDate = row.committedDate;
  }

  function cancelEdit() {
    editingId = null;
    editingDate = '';
  }

  function saveEdit(row: ListRow) {
    if (!editingDate) return;
    onSaveDate(row, editingDate);
    editingId = null;
  }

  function reassign(row: ListRow, event: Event) {
    const select = event.currentTarget as HTMLSelectElement;
    if (!select.value) return;
    onReassignFeature(row, select.value);
  }

  function sortLabel(key: SortKey): string {
    if (sortKey !== key) return '';
    return sortDirection === 'asc' ? ' ↑' : ' ↓';
  }
</script>

<section class="data-table" aria-label="Roadmap list table">
  {#if loading}
    <div class="state-card">Loading list view…</div>
  {:else if groups.length === 0}
    <div class="state-card">No roadmap items match the current filters.</div>
  {:else}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            {#each columns as column}
              <th scope="col">
                <button type="button" class="sort-button" on:click={() => onSort(column.key)}>
                  {column.label}{sortLabel(column.key)}
                </button>
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each groups as group}
            <tr class="group-row">
              <td colspan={columns.length}>{group.label}</td>
            </tr>
            {#each group.rows as row}
              <tr class:feature-row={row.type === 'feature'} class:synthetic={row.isSynthetic}>
                <td>
                  <div class="title-cell" class:indented={row.type === 'feature'}>
                    <span class="title-text">{row.title}</span>
                    {#if row.type === 'feature' && row.isSynthetic}
                      <label class="reassign-label">
                        <span>Reassign</span>
                        <select on:change={(event) => reassign(row, event)}>
                          <option value="">Choose epic</option>
                          {#each epicOptions as option}
                            <option value={option.id}>{option.title}</option>
                          {/each}
                        </select>
                      </label>
                    {/if}
                  </div>
                </td>
                <td><span class="type-pill">{row.type}</span></td>
                <td>{row.owner || '—'}</td>
                <td>{row.sprint || 'Unassigned'}</td>
                <td>{row.originalDate || '—'}</td>
                <td>
                  {#if row.type === 'feature'}
                    {#if editingId === row.id}
                      <div class="edit-cell">
                        <DatePicker
                          bind:value={editingDate}
                          showLabel={false}
                          showSnapToggle={false}
                          compact={true}
                          label="Committed date"
                          inputAriaLabel={`Committed date for ${row.title}`}
                          disabled={savingDateId === row.id}
                        />
                        <div class="edit-actions">
                          <button class="cell-action primary" type="button" disabled={!editingDate || savingDateId === row.id} on:click={() => saveEdit(row)}>
                            Save
                          </button>
                          <button class="cell-action" type="button" disabled={savingDateId === row.id} on:click={cancelEdit}>
                            Cancel
                          </button>
                        </div>
                      </div>
                    {:else}
                      <button class="date-display" type="button" on:click={() => beginEdit(row)}>
                        {row.committedDate || '— set date'}
                      </button>
                    {/if}
                  {:else}
                    {row.committedDate || '—'}
                  {/if}
                </td>
                <td>{row.actualDate || '—'}</td>
                <td>{row.slipEvents}</td>
                <td>{row.status || '—'}</td>
                <td>
                  <span class={`health health-${row.health}`} title={row.healthLabel}></span>
                </td>
              </tr>
            {/each}
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</section>

<style>
  .data-table {
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
    overflow: hidden;
  }

  .table-wrap {
    max-height: calc(100vh - 280px);
    overflow: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    min-width: 1120px;
  }

  thead th {
    position: sticky;
    top: 0;
    z-index: 2;
    padding: 12px 14px;
    border-bottom: 1px solid var(--border);
    background: var(--bg2);
    text-align: left;
  }

  .sort-button {
    border: 0;
    background: transparent;
    color: var(--text3);
    cursor: pointer;
    font-size: 12px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  .sort-button:hover {
    color: var(--text);
  }

  td {
    padding: 14px;
    border-bottom: 1px solid var(--border);
    color: var(--text2);
    font-size: 14px;
    vertical-align: top;
  }

  tr:last-child td {
    border-bottom: 0;
  }

  .group-row td {
    padding: 10px 14px;
    background: color-mix(in srgb, var(--accent) 5%, transparent);
    color: var(--accent);
    font-size: 12px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  .feature-row td {
    background: color-mix(in srgb, var(--text) 1%, transparent);
  }

  tr.synthetic td {
    color: var(--text3);
  }

  .title-cell {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .title-cell.indented {
    padding-left: 22px;
  }

  .title-text {
    color: var(--text);
  }

  .type-pill {
    display: inline-flex;
    padding: 3px 8px;
    border: 1px solid var(--border);
    border-radius: 999px;
    color: var(--text2);
    font-size: 12px;
    text-transform: capitalize;
  }

  .reassign-label {
    display: flex;
    flex-direction: column;
    gap: 6px;
    color: var(--text3);
    font-size: 11px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .reassign-label select {
    height: 34px;
    max-width: 220px;
    padding: 0 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 13px;
  }

  .date-display,
  .cell-action {
    border: 1px solid var(--border2);
    border-radius: 6px;
    background: transparent;
    color: var(--text2);
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    letter-spacing: 0.03em;
  }

  .date-display {
    padding: 7px 10px;
  }

  .date-display:hover,
  .cell-action:hover {
    background: var(--accent-a8);
    color: var(--text);
  }

  .edit-cell {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .edit-actions {
    display: flex;
    gap: 8px;
  }

  .cell-action {
    padding: 6px 10px;
  }

  .cell-action.primary {
    border-color: var(--accent2);
    color: var(--accent);
  }

  .health {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 999px;
  }

  .health-green { background: var(--green-light); }
  .health-amber { background: var(--amber); }
  .health-red { background: var(--coral); }
  .health-blue { background: var(--blue-light); }

  .state-card {
    padding: 24px;
    color: var(--text2);
  }
</style>
