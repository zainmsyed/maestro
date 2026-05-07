<script lang="ts">
  import { onMount } from 'svelte';
  import DataTable from '../components/DataTable.svelte';
  import FilterBar from '../components/FilterBar.svelte';
  import { api, formatApiError, type AuditRecord, type EpicRecord } from '../lib/api';
  import { downloadCsv } from '../lib/csvExport';
  import {
    buildExportRows,
    buildGroups,
    filterRows,
    flattenRows,
    uniqueLabels,
    type FilterState,
    type GroupBy,
    type ListGroup,
    type ListRow,
    type SortKey,
  } from '../lib/listView';

  let loading = true;
  let savingDateId: string | null = null;
  let error: string | null = null;
  let epics: EpicRecord[] = [];
  let audits: AuditRecord[] = [];
  let filters: FilterState = { epic: '', owner: '', sprint: '', status: '' };
  let groupBy: GroupBy = 'epic';
  let sortKey: SortKey = 'title';
  let sortDirection: 'asc' | 'desc' = 'asc';


  function sortBy(key: SortKey) {
    if (sortKey === key) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
      return;
    }
    sortKey = key;
    sortDirection = 'asc';
  }

  async function load() {
    loading = true;
    error = null;
    try {
      const [nextEpics, nextAudits] = await Promise.all([
        api.get<EpicRecord[]>('/epics'),
        api.get<AuditRecord[]>('/audit').catch(() => []),
      ]);
      epics = nextEpics;
      audits = nextAudits;
    } catch (err) {
      error = formatApiError(err, 'Could not load list view data.');
    } finally {
      loading = false;
    }
  }

  async function saveDate(row: ListRow, committedDate: string) {
    savingDateId = row.id;
    try {
      await api.patchDate('feature', row.id, committedDate);
      await load();
    } catch (err) {
      error = formatApiError(err, 'Could not save committed date.');
    } finally {
      savingDateId = null;
    }
  }

  async function reassignFeature(row: ListRow, epicID: string) {
    try {
      await api.patchFeatureEpic(row.id, epicID);
      await load();
    } catch (err) {
      error = formatApiError(err, 'Could not reassign feature.');
    }
  }

  function exportCurrentView() {
    downloadCsv('maestro-list-view.csv', buildExportRows(groupedRows, groupBy));
  }

  $: flatRows = flattenRows(epics, audits);
  $: filteredRows = filterRows(flatRows, filters);
  $: groupedRows = buildGroups(filteredRows, groupBy, sortKey, sortDirection);

  // Build options for one filter key from rows filtered by ALL OTHER keys.
  // This gives true cascading: picking an epic narrows Owner/Sprint/Status,
  // but the Epic dropdown still shows every epic that matches the other picks.
  function optionsForKey(
    key: keyof FilterState,
    rows: ListRow[],
    currentFilters: FilterState,
    fallback?: string,
  ): string[] {
    const otherFilters = { ...currentFilters, [key]: '' };
    const filtered = filterRows(rows, otherFilters);
    const source = filtered.length > 0 ? filtered : rows;
    const active = currentFilters[key];

    if (key === 'epic') {
      return uniqueLabels([...source.map((r) => r.epicTitle).filter(Boolean), active].filter(Boolean));
    }
    if (key === 'owner') {
      return uniqueLabels([...source.map((r) => r.owner).filter(Boolean), active].filter(Boolean));
    }
    if (key === 'sprint') {
      return uniqueLabels([...source.map((r) => r.sprint), ...(active ? [active] : [])], fallback);
    }
    return uniqueLabels([...source.map((r) => r.status).filter(Boolean), active].filter(Boolean));
  }

  $: epicOptions = optionsForKey('epic', flatRows, filters);
  $: ownerOptions = optionsForKey('owner', flatRows, filters);
  $: sprintOptions = optionsForKey('sprint', flatRows, filters, 'Unassigned');
  $: statusOptions = optionsForKey('status', flatRows, filters);

  $: reassignOptions = epics.filter((epic) => !epic.is_synthetic).map((epic) => ({ id: epic.id, title: epic.title }));

  function handleFilterChange(next: FilterState) {
    filters = next;
  }

  onMount(load);
</script>

<section class="list-view" aria-label="List view">
  <div class="screen-header">
    <div>
      <p class="eyebrow">List</p>
      <h1>Roadmap list view</h1>
      <p class="intro">
        Sort, filter, group, edit committed dates inline, and export the current working view.
      </p>
    </div>
    <div class="header-summary">
      <div><span>{filteredRows.filter((row) => row.type === 'epic').length}</span><small>Epics</small></div>
      <div><span>{filteredRows.filter((row) => row.type === 'feature').length}</span><small>Features</small></div>
      <div><span>{groupedRows.length}</span><small>Groups</small></div>
    </div>
  </div>

  <FilterBar
    {filters}
    onFilterChange={handleFilterChange}
    bind:groupBy
    {epicOptions}
    {ownerOptions}
    {sprintOptions}
    {statusOptions}
    onExport={exportCurrentView}
  />

  {#if error}
    <div class="error-banner">{error}</div>
  {/if}

  <DataTable
    groups={groupedRows}
    {sortKey}
    {sortDirection}
    {loading}
    {savingDateId}
    epicOptions={reassignOptions}
    onSort={sortBy}
    onSaveDate={saveDate}
    onReassignFeature={reassignFeature}
  />
</section>

<style>
  .list-view {
    display: flex;
    flex-direction: column;
    gap: 20px;
    width: min(1320px, 100%);
    margin: 0 auto;
  }

  .screen-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 24px;
  }

  .eyebrow {
    margin: 0 0 10px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  h1 {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-size: clamp(34px, 5vw, 56px);
    font-weight: 300;
    letter-spacing: -0.05em;
    line-height: 1.02;
  }

  .intro {
    max-width: 680px;
    margin: 12px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .header-summary {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
  }

  .header-summary div {
    min-width: 110px;
    padding: 14px 16px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: linear-gradient(145deg, rgba(29, 29, 32, 0.94), rgba(22, 22, 24, 0.94));
  }

  .header-summary span {
    display: block;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 28px;
    font-weight: 300;
  }

  .header-summary small {
    display: block;
    margin-top: 4px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .error-banner {
    padding: 16px 18px;
    border: 1px solid rgba(196, 97, 74, 0.45);
    border-radius: 8px;
    background: var(--bg2);
    color: var(--coral);
  }

  @media (max-width: 980px) {
    .screen-header {
      flex-direction: column;
    }
  }
</style>
