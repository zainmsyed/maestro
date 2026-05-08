<script lang="ts">
  import DatePicker from './DatePicker.svelte';
  import DateSourceBadge from './DateSourceBadge.svelte';
  import type { ListGroup, ListRow, SortKey } from '../lib/listView';

  export let groups: ListGroup[] = [];
  export let sortKey: SortKey = 'title';
  export let sortDirection: 'asc' | 'desc' = 'asc';
  export let loading = false;
  export let savingDateId: string | null = null;
  export let epicOptions: Array<{ id: string; title: string }> = [];
  export let featureOptions: Array<{ id: string; title: string }> = [];
  export let onSort: (key: SortKey) => void;
  export let onSaveDate: (row: ListRow, date: string) => void;
  export let onReassignFeature: (row: ListRow, epicID: string) => void;
  export let onReassignStory: (row: ListRow, featureID: string) => void;

  const columns: Array<{ key: SortKey; label: string }> = [
    { key: 'title', label: 'Title' },
    { key: 'type', label: 'Type' },
    { key: 'owner', label: 'Owner' },
    { key: 'sprint', label: 'Sprint' },
    { key: 'original', label: 'Original Date' },
    { key: 'committed', label: 'Committed Date' },
    { key: 'actual', label: 'Actual Date' },
    { key: 'status', label: 'Status' },
    { key: 'health', label: 'Health' },
    { key: 'slip', label: 'Slip Events' },
  ];

  let editingId: string | null = null;
  let editingDate = '';
  let collapsedEpics: Record<string, boolean> = {};
  let collapsedFeatures: Record<string, boolean> = {};

  function collapseAll() {
    const nextEpics: Record<string, boolean> = {};
    const nextFeatures: Record<string, boolean> = {};
    for (const group of groups) {
      for (const row of group.rows) {
        if (row.type === 'epic') nextEpics[row.id] = true;
        if (row.type === 'feature') nextFeatures[row.id] = true;
      }
    }
    collapsedEpics = nextEpics;
    collapsedFeatures = nextFeatures;
  }

  function expandAll() {
    collapsedEpics = {};
    collapsedFeatures = {};
  }

  function toggleEpic(epicId: string) {
    collapsedEpics = { ...collapsedEpics, [epicId]: !collapsedEpics[epicId] };
  }

  function toggleFeature(featureId: string) {
    collapsedFeatures = { ...collapsedFeatures, [featureId]: !collapsedFeatures[featureId] };
  }

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

  function reassignFeature(row: ListRow, event: Event) {
    const select = event.currentTarget as HTMLSelectElement;
    if (!select.value) return;
    onReassignFeature(row, select.value);
  }

  function reassignStory(row: ListRow, event: Event) {
    const select = event.currentTarget as HTMLSelectElement;
    if (!select.value) return;
    onReassignStory(row, select.value);
  }

  function sortLabel(key: SortKey): string {
    if (sortKey !== key) return '';
    return sortDirection === 'asc' ? ' ↑' : ' ↓';
  }

  function indentClass(type: ListRow['type']): string {
    if (type === 'story') return 'indent-story';
    if (type === 'feature') return 'indent-feature';
    return '';
  }


</script>

<section class="data-table" aria-label="Roadmap list table">
  {#if loading}
    <div class="state-card">Loading list view…</div>
  {:else if groups.length === 0}
    <div class="state-card">No roadmap items match the current filters.</div>
  {:else}
    <div class="table-actions">
      <button type="button" class="action-link" on:click={collapseAll}>Collapse all</button>
      <button type="button" class="action-link" on:click={expandAll}>Expand all</button>
    </div>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            {#each columns as column}
              <th scope="col" class:col-numeric={column.key === 'slip'}>
                <button type="button" class="sort-button" on:click={() => onSort(column.key)}>
                  {column.label}{sortLabel(column.key)}
                </button>
              </th>
            {/each}
            <th scope="col" class="col-numeric">Story Points</th>
            <th scope="col">Date Source</th>
          </tr>
        </thead>
        <tbody>
          {#each groups as group}
            <tr class="group-row">
              <td colspan={12}>{group.label}</td>
            </tr>
            {#each group.rows as row}
              {#if row.type === 'epic' || (row.type === 'feature' && !collapsedEpics[row.epicId ?? 'none']) || (row.type === 'story' && !collapsedEpics[row.epicId ?? 'none'] && !collapsedFeatures[row.featureId ?? 'none'])}
                <tr
                  class:feature-row={row.type === 'feature'}
                  class:story-row={row.type === 'story'}
                  class:synthetic={row.isSynthetic}
                >
                  <td>
                    <div class="title-cell {indentClass(row.type)}">
                      {#if row.type === 'epic'}
                        <button
                          type="button"
                          class="expand-toggle"
                          on:click={() => toggleEpic(row.id)}
                          aria-label={!collapsedEpics[row.id] ? 'Collapse epic' : 'Expand epic'}
                        >
                          {!collapsedEpics[row.id] ? '▼' : '▶'}
                        </button>
                      {:else if row.type === 'feature'}
                        <button
                          type="button"
                          class="expand-toggle"
                          on:click={() => toggleFeature(row.id)}
                          aria-label={!collapsedFeatures[row.id] ? 'Collapse feature' : 'Expand feature'}
                        >
                          {!collapsedFeatures[row.id] ? '▼' : '▶'}
                        </button>
                      {:else}
                        <span class="expand-spacer"></span>
                      {/if}
                      <span class="title-text">{row.title}</span>
                      {#if row.type === 'feature' && row.isSynthetic}
                        <label class="reassign-label">
                          <span>Reassign</span>
                          <select on:change={(event) => reassignFeature(row, event)}>
                            <option value="">Choose epic</option>
                            {#each epicOptions as option}
                              <option value={option.id}>{option.title}</option>
                            {/each}
                          </select>
                        </label>
                      {/if}
                      {#if row.type === 'story' && row.isSynthetic}
                        <label class="reassign-label">
                          <span>Reassign</span>
                          <select on:change={(event) => reassignStory(row, event)}>
                            <option value="">Choose feature</option>
                            {#each featureOptions as option}
                              <option value={option.id}>{option.title}</option>
                            {/each}
                          </select>
                        </label>
                      {/if}
                    </div>
                  </td>
                  <td><span class="type-pill type-{row.type}">{row.type}</span></td>
                  <td>{row.owner || '—'}</td>
                  <td>{row.sprint || 'Unassigned'}</td>
                  <td>{row.originalDate || '—'}</td>
                  <td>
                    {#if row.type !== 'epic'}
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
                            <button
                              class="cell-action primary"
                              type="button"
                              disabled={!editingDate || savingDateId === row.id}
                              on:click={() => saveEdit(row)}
                            >
                              Save
                            </button>
                            <button
                              class="cell-action"
                              type="button"
                              disabled={savingDateId === row.id}
                              on:click={cancelEdit}
                            >
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
                  <td>{row.status || '—'}</td>
                  <td>
                    <span class={`health health-${row.health}`} title={row.healthLabel}></span>
                  </td>
                  <td class="col-numeric">{row.slipEvents}</td>
                  <td class="col-numeric">{row.storyPoints ?? '—'}</td>
                  <td>
                    {#if row.dateSource}
                      <DateSourceBadge source={row.dateSource} />
                    {:else}
                      <span class="muted">—</span>
                    {/if}
                  </td>
                </tr>
              {/if}
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
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    min-width: 1080px;
  }

  thead th {
    position: sticky;
    top: 0;
    z-index: 2;
    padding: 10px 12px;
    border-bottom: 1px solid var(--border);
    background: var(--bg2);
    text-align: left;
    white-space: nowrap;
    color: var(--text3);
    font-size: 11px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .col-numeric {
    text-align: right;
  }

  .sort-button {
    border: 0;
    background: transparent;
    color: inherit;
    cursor: pointer;
    font: inherit;
  }

  .sort-button:hover {
    color: var(--text);
  }

  td {
    padding: 10px 12px;
    border-bottom: 1px solid var(--border);
    color: var(--text2);
    font-size: 12px;
    vertical-align: top;
    white-space: nowrap;
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

  .story-row td {
    background: transparent;
  }

  tr.synthetic td {
    color: var(--text3);
  }

  .title-cell {
    display: flex;
    align-items: flex-start;
    gap: 8px;
  }

  .indent-feature {
    padding-left: 28px;
  }

  .indent-story {
    padding-left: 52px;
  }

  .expand-toggle {
    flex-shrink: 0;
    width: 16px;
    height: 16px;
    padding: 0;
    border: 0;
    background: transparent;
    color: var(--text3);
    cursor: pointer;
    font-size: 9px;
    line-height: 16px;
    text-align: center;
  }

  .expand-toggle:hover {
    color: var(--text);
  }

  .expand-spacer {
    flex-shrink: 0;
    width: 18px;
    display: inline-block;
  }

  .title-text {
    color: var(--text);
    flex: 1;
    max-width: 280px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .type-pill {
    display: inline-flex;
    padding: 2px 6px;
    border: 1px solid var(--border);
    border-radius: 999px;
    color: var(--text2);
    font-size: 11px;
    text-transform: capitalize;
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

  .type-story {
    border-color: color-mix(in srgb, var(--green) 35%, transparent);
    background: color-mix(in srgb, var(--green) 8%, transparent);
    color: var(--green-light);
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

  .muted {
    color: var(--text3);
  }

  .table-actions {
    display: flex;
    gap: 16px;
    padding: 12px 14px;
    border-bottom: 1px solid var(--border);
  }

  .action-link {
    padding: 0;
    border: 0;
    background: transparent;
    color: var(--text3);
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 12px;
    letter-spacing: 0.04em;
    text-decoration: underline;
    text-underline-offset: 3px;
  }

  .action-link:hover {
    color: var(--accent);
  }

  .state-card {
    padding: 24px;
    color: var(--text2);
  }
</style>
