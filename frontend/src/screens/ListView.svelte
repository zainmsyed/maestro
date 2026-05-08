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
  let filters: FilterState = { epic: '', owner: '', sprint: '', status: '', type: '' };
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
      await api.patchDate(row.type, row.id, committedDate);
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

  async function reassignStory(row: ListRow, featureID: string) {
    try {
      await api.patchStoryFeature(row.id, featureID);
      await load();
    } catch (err) {
      error = formatApiError(err, 'Could not reassign story.');
    }
  }

  function exportCurrentView() {
    downloadCsv('maestro-list-view.csv', buildExportRows(groupedRows, groupBy));
  }

  $: flatRows = flattenRows(epics, audits);
  $: filteredRows = filterRows(flatRows, filters);
  $: groupedRows = buildGroups(filteredRows, groupBy, sortKey, sortDirection);

  const optionConfig: Record<
    keyof FilterState,
    { accessor: (r: ListRow) => string | null | undefined; fallback?: string; filterEmpty: boolean }
  > = {
    epic: { accessor: (r) => r.epicTitle, filterEmpty: true },
    owner: { accessor: (r) => r.owner, filterEmpty: true },
    sprint: { accessor: (r) => r.sprint, fallback: 'Unassigned', filterEmpty: false },
    status: { accessor: (r) => r.status, filterEmpty: true },
    type: { accessor: (r) => r.type, filterEmpty: true },
  };

  function optionsForKey(
    key: keyof FilterState,
    rows: ListRow[],
    currentFilters: FilterState,
  ): string[] {
    const otherFilters = { ...currentFilters, [key]: '' };
    const filtered = filterRows(rows, otherFilters);
    const source = filtered.length > 0 ? filtered : rows;
    const active = currentFilters[key];
    const config = optionConfig[key];

    let values = source.map(config.accessor);
    if (config.filterEmpty) {
      values = values.filter(Boolean);
    }
    if (active) values.push(active);

    return uniqueLabels(values, config.fallback);
  }

  $: epicOptions = optionsForKey('epic', flatRows, filters);
  $: ownerOptions = optionsForKey('owner', flatRows, filters);
  $: sprintOptions = optionsForKey('sprint', flatRows, filters);
  $: statusOptions = optionsForKey('status', flatRows, filters);
  $: typeOptions = optionsForKey('type', flatRows, filters);

  $: reassignEpicOptions = epics.filter((epic) => !epic.is_synthetic).map((epic) => ({ id: epic.id, title: epic.title }));
  $: reassignFeatureOptions = flatRows
    .filter((row) => row.type === 'feature' && !row.isSynthetic)
    .map((row) => ({ id: row.id, title: row.title }));

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
      <div><span>{filteredRows.filter((row) => row.type === 'story').length}</span><small>Stories</small></div>
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
    {typeOptions}
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
    epicOptions={reassignEpicOptions}
    featureOptions={reassignFeatureOptions}
    onSort={sortBy}
    onSaveDate={saveDate}
    onReassignFeature={reassignFeature}
    onReassignStory={reassignStory}
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
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 12px;
  }

  .header-summary div {
    min-width: 110px;
    padding: 14px 16px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--surface-glass);
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
    border: 1px solid color-mix(in srgb, var(--coral) 45%, transparent);
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
